package testmcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/audit"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/policy"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/proxy"
)

// testPolicies returns policies tailored for integration testing.
func testPolicies() []config.PolicyCfg {
	return []config.PolicyCfg{
		{
			Name:   "allow-reads",
			Action: "allow",
			Match: config.PolicyMatch{
				Identity: "*",
				Tools:    []string{"read_*", "echo", "tools/list"},
			},
		},
		{
			Name:   "block-deletes",
			Action: "block",
			Match: config.PolicyMatch{
				Tools: []string{"delete_*"},
			},
			Alert: true,
		},
		{
			Name:   "payment-approval",
			Action: "hitl",
			Match: config.PolicyMatch{
				Tools: []string{"execute_payout"},
			},
			Constraints: &config.Constraints{MaxAmount: 500},
		},
	}
}

// sendJSONRPC sends a JSON-RPC message and returns the decoded response.
func sendJSONRPC(t *testing.T, conn net.Conn, method string, params any) map[string]any {
	t.Helper()
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      uuid.New().String(),
		"method":  method,
	}
	if params != nil {
		req["params"] = params
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if _, err = fmt.Fprintln(conn, string(data)); err != nil {
		t.Fatalf("send request: %v", err)
	}
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	var resp map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &resp); err != nil {
		t.Fatalf("parse response: %v (raw: %s)", err, line)
	}
	return resp
}

// TestIntegration_WithFixedPort runs with a predetermined proxy port.
func TestIntegration_WithFixedPort(t *testing.T) {
	// 1. Start test MCP server
	server, err := NewServer()
	if err != nil {
		t.Fatalf("start test server: %v", err)
	}
	server.Start()
	defer server.Stop()

	// 2. Find a free port for the proxy
	proxyListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve proxy port: %v", err)
	}
	proxyPort := proxyListener.Addr().(*net.TCPAddr).Port
	proxyListener.Close() // Release it — proxy will re-listen
	proxyAddr := fmt.Sprintf("127.0.0.1:%d", proxyPort)

	t.Logf("Proxy listening on %s, upstream to %s", proxyAddr, server.Addr())

	// 3. Create audit logger
	auditDir := t.TempDir()
	auditPath := filepath.Join(auditDir, "audit.jsonl")
	auditLogger, err := audit.NewLogger(config.AuditConfig{
		Path:    auditPath,
		HMACKey: "integration-test-key",
	})
	if err != nil {
		t.Fatalf("audit logger: %v", err)
	}
	defer auditLogger.Close()

	// 4. Create and start proxy
	p := proxy.New(proxy.Options{
		Mode:     "tcp",
		Listen:   proxyAddr,
		Upstream: server.Addr(),
		Policy:   policy.NewEngine(testPolicies()),
		Audit:    auditLogger,
	})

	proxyErr := make(chan error, 1)
	go func() {
		proxyErr <- p.Start()
	}()
	defer p.Stop()

	time.Sleep(200 * time.Millisecond) // Wait for proxy to start

	// 5. Connect test client to proxy
	client, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		t.Fatalf("connect to proxy: %v", err)
	}
	defer client.Close()

	// 6. Initialize MCP handshake
	t.Log("Sending initialize...")
	initResp := sendJSONRPC(t, client, "initialize", map[string]any{
		"protocolVersion": "2025-11-05",
		"clientInfo": map[string]any{
			"name":    "integration-test",
			"version": "1.0.0",
		},
	})
	if initResp["error"] != nil {
		t.Fatalf("initialize error: %v", initResp["error"])
	}
	t.Logf("Initialize OK")

	// 7. Tools/list
	t.Log("Sending tools/list...")
	listResp := sendJSONRPC(t, client, "tools/list", nil)
	if listResp["error"] != nil {
		t.Fatalf("tools/list error: %v", listResp["error"])
	}
	t.Logf("Tools/list OK")

	// 8. Test: read_db → should be ALLOWED
	t.Log("Test: read_db → expect ALLOW")
	resp := sendJSONRPC(t, client, "tools/call", map[string]any{
		"name": "read_db",
		"arguments": map[string]any{"query": "SELECT * FROM users"},
	})
	if resp["error"] != nil {
		t.Fatalf("read_db blocked but should be allowed: %v", resp["error"])
	}
	t.Logf("  ✓ read_db allowed")

	// 9. Test: delete_db → should be BLOCKED
	t.Log("Test: delete_db → expect BLOCK")
	resp = sendJSONRPC(t, client, "tools/call", map[string]any{
		"name": "delete_db",
		"arguments": map[string]any{"table": "users", "id": 1},
	})
	if resp["error"] == nil {
		t.Fatal("delete_db was allowed but should be blocked!")
	}
	errMsg := resp["error"].(map[string]any)
	if errMsg["code"] != float64(-32000) {
		t.Errorf("expected error code -32000, got %v", errMsg["code"])
	}
	t.Logf("  ✓ delete_db blocked")

	// 10. Test: echo → should be ALLOWED
	t.Log("Test: echo → expect ALLOW")
	resp = sendJSONRPC(t, client, "tools/call", map[string]any{
		"name": "echo",
		"arguments": map[string]any{"message": "hello"},
	})
	if resp["error"] != nil {
		t.Fatalf("echo blocked but should be allowed: %v", resp["error"])
	}
	t.Logf("  ✓ echo allowed")

	// 11. Test: execute_payout → should be HITL PENDING
	t.Log("Test: execute_payout → expect HITL")
	resp = sendJSONRPC(t, client, "tools/call", map[string]any{
		"name": "execute_payout",
		"arguments": map[string]any{"amount": 500, "currency": "USD"},
	})
	if resp["error"] == nil {
		t.Fatal("execute_payout was allowed but should require HITL!")
	}
	errMsg = resp["error"].(map[string]any)
	t.Logf("  ✓ execute_payout blocked (HITL): %s", errMsg["message"])

	// 12. Test: unknown_tool → should be BLOCKED (default deny)
	t.Log("Test: unknown_tool → expect BLOCK (default deny)")
	resp = sendJSONRPC(t, client, "tools/call", map[string]any{
		"name": "unknown_tool",
		"arguments": map[string]any{},
	})
	if resp["error"] == nil {
		t.Fatal("unknown_tool was allowed but should be blocked by default deny!")
	}
	t.Logf("  ✓ unknown_tool blocked (default deny)")

	// 13. Verify audit log has entries
	client.Close()
	p.Stop()

	auditData, err := os.ReadFile(auditPath)
	if err != nil {
		t.Fatalf("read audit log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(auditData)), "\n")
	if len(lines) < 5 {
		t.Fatalf("expected at least 5 audit entries, got %d", len(lines))
	}
	t.Logf("  ✓ Audit log has %d entries", len(lines))

	// 14. Verify HMAC chain intact
	verifier := audit.NewVerifier(auditPath)
	valid, err := verifier.Verify()
	if err != nil {
		t.Fatalf("audit verify error: %v", err)
	}
	if !valid {
		t.Fatal("audit HMAC chain broken!")
	}
	t.Logf("  ✓ Audit HMAC chain intact")

	// Check proxy exit
	select {
	case err := <-proxyErr:
		if err != nil {
			t.Fatalf("proxy error: %v", err)
		}
	default:
	}
}
