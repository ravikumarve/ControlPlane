package proxy

import (
	"encoding/json"
	"fmt"
)

// JSONRPCRequest represents an incoming MCP JSON-RPC request.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any         `json:"id"`
	Method  string      `json:"method"`
	Params  any         `json:"params,omitempty"`
}

// JSONRPCResponse represents an MCP JSON-RPC response.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any         `json:"id"`
	Result  any         `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error object.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ToolDefinition represents a tool as defined in MCP's tools/list response.
type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

// ParseRequest parses a raw JSON-RPC message.
func ParseRequest(data []byte) (*JSONRPCRequest, error) {
	var req JSONRPCRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON-RPC request: %w", err)
	}
	if req.JSONRPC != "2.0" {
		return nil, fmt.Errorf("unsupported JSON-RPC version: %s", req.JSONRPC)
	}
	return &req, nil
}

// ParseResponse parses a JSON-RPC response.
func ParseResponse(data []byte) (*JSONRPCResponse, error) {
	var resp JSONRPCResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("invalid JSON-RPC response: %w", err)
	}
	return &resp, nil
}

// NewBlockedResponse creates a JSON-RPC error response for a blocked call.
func NewBlockedResponse(id any, reason string) []byte {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    -32000,
			Message: "blocked by security policy",
			Data:    map[string]string{"reason": reason},
		},
	}
	data, _ := json.Marshal(resp)
	return data
}

// ExtractToolName extracts the actual tool name from a tools/call request.
// For MCP, the method is "tools/call" and the tool name is in params.name.
// Returns the JSON-RPC method for non-tool calls.
func ExtractToolName(req *JSONRPCRequest) string {
	if req.Method == "tools/call" && req.Params != nil {
		if params, ok := req.Params.(map[string]any); ok {
			if name, ok := params["name"].(string); ok {
				return name
			}
		}
		return "tools/call"
	}
	return req.Method
}

// NewErrorResponse creates a generic JSON-RPC error response.
func NewErrorResponse(id any, code int, message string) []byte {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
	data, _ := json.Marshal(resp)
	return data
}
