package proxy

import (
	"testing"
)

func TestParseRequest_Valid(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":1,"method":"read_database","params":{"query":"SELECT * FROM users"}}`)

	req, err := ParseRequest(data)
	if err != nil {
		t.Fatalf("ParseRequest failed: %v", err)
	}

	if req.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q; want 2.0", req.JSONRPC)
	}
	if req.Method != "read_database" {
		t.Errorf("Method = %q; want read_database", req.Method)
	}
	if req.ID != float64(1) { // JSON numbers decode as float64
		t.Errorf("ID = %v (%T); want 1", req.ID, req.ID)
	}
}

func TestParseRequest_InvalidJSON(t *testing.T) {
	data := []byte(`not json`)

	_, err := ParseRequest(data)
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestParseRequest_WrongVersion(t *testing.T) {
	data := []byte(`{"jsonrpc":"1.0","id":1,"method":"test"}`)

	_, err := ParseRequest(data)
	if err == nil {
		t.Fatal("Expected error for wrong JSON-RPC version")
	}
}

func TestParseRequest_NoMethod(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":1}`)

	req, err := ParseRequest(data)
	if err != nil {
		t.Fatalf("ParseRequest failed: %v", err)
	}
	if req.Method != "" {
		t.Errorf("Method should be empty; got %q", req.Method)
	}
}

func TestParseRequest_WithParams(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":2,"method":"execute_payout","params":{"amount":500,"currency":"USD"}}`)

	req, err := ParseRequest(data)
	if err != nil {
		t.Fatalf("ParseRequest failed: %v", err)
	}

	params, ok := req.Params.(map[string]any)
	if !ok {
		t.Fatalf("Params type = %T; want map[string]any", req.Params)
	}
	if params["amount"] != float64(500) {
		t.Errorf("amount = %v; want 500", params["amount"])
	}
	if params["currency"] != "USD" {
		t.Errorf("currency = %v; want USD", params["currency"])
	}
}

func TestParseResponse_Valid(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":1,"result":{"data":"ok"}}`)

	resp, err := ParseResponse(data)
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q; want 2.0", resp.JSONRPC)
	}
	if resp.Error != nil {
		t.Errorf("Error should be nil; got %v", resp.Error)
	}
}

func TestParseResponse_WithError(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"Method not found"}}`)

	resp, err := ParseResponse(data)
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("Error should not be nil")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("Error code = %d; want -32601", resp.Error.Code)
	}
	if resp.Error.Message != "Method not found" {
		t.Errorf("Error message = %q; want 'Method not found'", resp.Error.Message)
	}
}

func TestNewBlockedResponse(t *testing.T) {
	data := NewBlockedResponse(1, "blocked by policy: block-dangerous")

	resp, err := ParseResponse(data)
	if err != nil {
		t.Fatalf("Parse blocked response failed: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("Blocked response should have error")
	}
	if resp.Error.Code != -32000 {
		t.Errorf("Error code = %d; want -32000", resp.Error.Code)
	}
	if resp.Error.Message != "blocked by security policy" {
		t.Errorf("Error message = %q; want 'blocked by security policy'", resp.Error.Message)
	}
}

func TestNewBlockedResponse_StringID(t *testing.T) {
	data := NewBlockedResponse("req-123", "requires approval")

	resp, err := ParseResponse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if resp.ID != "req-123" {
		t.Errorf("ID = %v; want req-123", resp.ID)
	}
}

func TestNewErrorResponse(t *testing.T) {
	data := NewErrorResponse(1, -32001, "custom error")

	resp, err := ParseResponse(data)
	if err != nil {
		t.Fatalf("Parse error response failed: %v", err)
	}

	if resp.Error.Code != -32001 {
		t.Errorf("Error code = %d; want -32001", resp.Error.Code)
	}
	if resp.Error.Message != "custom error" {
		t.Errorf("Error message = %q; want 'custom error'", resp.Error.Message)
	}
}

func TestDetectIdentity_Anonymous(t *testing.T) {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test",
	}

	id := detectIdentity(req)
	if id != "anonymous" {
		t.Errorf("detectIdentity = %q; want anonymous", id)
	}
}

func TestDetectIdentity_WithIdentity(t *testing.T) {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test",
		Params:  map[string]any{"_identity": "payment-bot"},
	}

	id := detectIdentity(req)
	if id != "payment-bot" {
		t.Errorf("detectIdentity = %q; want payment-bot", id)
	}
}

func TestDetectIdentity_NilParams(t *testing.T) {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test",
		Params:  nil,
	}

	id := detectIdentity(req)
	if id != "anonymous" {
		t.Errorf("detectIdentity = %q; want anonymous", id)
	}
}
