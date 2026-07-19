package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
)

func TestLogger_WriteAndHMACChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")

	cfg := config.AuditConfig{
		Path:    path,
		HMACKey: "test-hmac-key-12345",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// Write 3 entries
	entries := []AuditEntry{
		{ID: uuid.New().String(), Timestamp: time.Now(), Identity: "agent1", Tool: "read_db", Decision: "allow", Duration: 5},
		{ID: uuid.New().String(), Timestamp: time.Now(), Identity: "agent2", Tool: "write_db", Decision: "block", Duration: 2},
		{ID: uuid.New().String(), Timestamp: time.Now(), Identity: "agent1", Tool: "execute_payout", Decision: "pending", Params: map[string]any{"amount": 100}},
	}

	for _, e := range entries {
		if err := logger.Write(e); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
	}

	// Verify each entry has HMAC and PrevHMAC chain
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Read audit file failed: %v", err)
	}

	// Verify file has 3 lines
	lines := splitLines(data)
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}

	// Verify HMAC chain by reading through Verifier
	verifier := NewVerifier(path)
	valid, err := verifier.Verify()
	if err != nil {
		t.Fatalf("Verifier.Verify failed: %v (chain may be broken)", err)
	}
	if !valid {
		t.Fatal("HMAC chain is broken")
	}
}

func TestLogger_HMACChainBroken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")

	cfg := config.AuditConfig{
		Path:    path,
		HMACKey: "test-key",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Write one entry
	logger.Write(AuditEntry{
		ID: uuid.New().String(), Timestamp: time.Now(), Identity: "agent", Tool: "test", Decision: "allow",
	})
	logger.Close()

	// Tamper with the file
	data, _ := os.ReadFile(path)
	tampered := string(data)
	tampered = tampered[:len(tampered)-1] + " {\\\"tampered\\\":true}\n"
	os.WriteFile(path, []byte(tampered), 0600)

	// Verify should detect tampering
	verifier := NewVerifier(path)
	valid, err := verifier.Verify()
	if err == nil && valid {
		t.Fatal("Expected tampered audit log to fail verification")
	}
}

func TestLogger_EmptyLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")

	verifier := NewVerifier(path)
	valid, err := verifier.Verify()
	// Non-existent file is not an error — just returns false
	if err == nil && valid {
		t.Fatal("Empty log should not be 'valid'")
	}
}

func TestLogger_AppendPreservesChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "append.jsonl")

	cfg := config.AuditConfig{Path: path, HMACKey: "append-key"}

	// First session
	logger1, _ := NewLogger(cfg)
	logger1.Write(AuditEntry{ID: uuid.New().String(), Timestamp: time.Now(), Identity: "a", Tool: "t1", Decision: "allow"})
	logger1.Close()

	// Second session (append)
	logger2, _ := NewLogger(cfg)
	logger2.Write(AuditEntry{ID: uuid.New().String(), Timestamp: time.Now(), Identity: "a", Tool: "t2", Decision: "block"})
	logger2.Close()

	// Verify full chain intact
	verifier := NewVerifier(path)
	valid, err := verifier.Verify()
	if err != nil {
		t.Fatalf("Append verification error: %v", err)
	}
	if !valid {
		t.Fatal("Append should preserve HMAC chain")
	}
}

func TestLogger_ComputeHMAC(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hmac-test.jsonl")
	cfg := config.AuditConfig{Path: path, HMACKey: "test-key"}
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// Same entry should produce same HMAC
	entry1 := AuditEntry{ID: "test-1", Timestamp: time.Now(), Identity: "agent", Tool: "tool1", Decision: "allow"}
	entry2 := AuditEntry{ID: "test-1", Timestamp: entry1.Timestamp, Identity: "agent", Tool: "tool1", Decision: "allow"}

	h1 := logger.computeHMAC(entry1)
	h2 := logger.computeHMAC(entry2)

	if h1 != h2 {
		t.Error("Same entry should produce same HMAC")
	}

	// Different entry should produce different HMAC
	entry3 := AuditEntry{ID: "test-2", Timestamp: time.Now(), Identity: "agent", Tool: "tool2", Decision: "block"}
	h3 := logger.computeHMAC(entry3)
	if h1 == h3 {
		t.Error("Different entries should produce different HMACs")
	}
}

func TestLogger_KeyFromEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "env-key.jsonl")

	os.Setenv("MCP_GUARD_HMAC_KEY", "env-based-key")
	defer os.Unsetenv("MCP_GUARD_HMAC_KEY")

	cfg := config.AuditConfig{
		Path:   path,
		HMACKey: "${MCP_GUARD_HMAC_KEY}",
	}

	// Simulate expansion (normally done by config.Load)
	cfg.HMACKey = os.Getenv("MCP_GUARD_HMAC_KEY")

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger with env key failed: %v", err)
	}
	defer logger.Close()

	err = logger.Write(AuditEntry{
		ID: uuid.New().String(), Timestamp: time.Now(), Identity: "agent", Tool: "test", Decision: "allow",
	})
	if err != nil {
		t.Fatalf("Write with env key failed: %v", err)
	}
}

func TestLogger_ConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "concurrent.jsonl")

	logger, _ := NewLogger(config.AuditConfig{Path: path, HMACKey: "concurrent-key"})
	defer logger.Close()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			err := logger.Write(AuditEntry{
				ID: uuid.New().String(), Timestamp: time.Now(),
				Identity: "agent", Tool: "tool", Decision: "allow",
			})
			if err != nil {
				t.Errorf("Concurrent write %d failed: %v", n, err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all 10 entries are in the file
	verifier := NewVerifier(path)
	valid, err := verifier.Verify()
	if err != nil {
		t.Fatalf("Concurrent verification error: %v", err)
	}
	if !valid {
		t.Fatal("Concurrent writes broke HMAC chain")
	}
}
