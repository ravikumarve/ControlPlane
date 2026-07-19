package testmcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
)

// ToolDefinition for MCP tools/list response.
type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

// Server is a minimal MCP test server.
type Server struct {
	addr     string
	listener net.Listener
	tools    []ToolDefinition
	wg       sync.WaitGroup
	stopped  bool
	mu       sync.Mutex
}

// NewServer creates a test MCP server on a random port.
func NewServer() (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("testmcp listen: %w", err)
	}

	tools := []ToolDefinition{
		{
			Name:        "read_db",
			Description: "Read data from the database",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{"type": "string"},
				},
			},
		},
		{
			Name:        "delete_db",
			Description: "Delete data from the database",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"table": map[string]any{"type": "string"},
					"id":    map[string]any{"type": "integer"},
				},
			},
		},
		{
			Name:        "execute_payout",
			Description: "Execute a payment payout",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"amount":   map[string]any{"type": "number"},
					"currency": map[string]any{"type": "string"},
				},
			},
		},
		{
			Name:        "echo",
			Description: "Echo back the input",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{"type": "string"},
				},
			},
		},
	}

	return &Server{
		addr:     listener.Addr().String(),
		listener: listener,
		tools:    tools,
	}, nil
}

// Addr returns the server's listen address.
func (s *Server) Addr() string { return s.addr }

// Start begins accepting connections in the background.
func (s *Server) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.mu.Lock()
				stopped := s.stopped
				s.mu.Unlock()
				if stopped {
					return
				}
				continue
			}
			s.wg.Add(1)
			go s.handleConn(conn)
		}
	}()
}

// Stop shuts down the server.
func (s *Server) Stop() {
	s.mu.Lock()
	s.stopped = true
	s.mu.Unlock()
	s.listener.Close()
	// Don't wait for stuck reader goroutines — listener close is sufficient
}

// handleConn handles a single TCP connection.
func (s *Server) handleConn(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		response := s.processMessage(line)
		if response != nil {
			fmt.Fprintln(conn, string(response))
		}
	}
}

// processMessage handles a single JSON-RPC message and returns a response.
func (s *Server) processMessage(data string) []byte {
	var msg struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Method  string `json:"method"`
		Params  any    `json:"params,omitempty"`
	}
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return s.errorResponse(nil, -32700, "Parse error")
	}
	if msg.JSONRPC != "2.0" {
		return s.errorResponse(msg.ID, -32600, "Invalid Request")
	}

	switch msg.Method {
	case "initialize":
		return s.initializeResponse(msg.ID)
	case "notifications/initialized":
		return nil // no response needed for notifications
	case "tools/list":
		return s.toolsListResponse(msg.ID)
	case "tools/call":
		return s.toolsCallResponse(msg.ID, msg.Params)
	default:
		return s.errorResponse(msg.ID, -32601, "Method not found")
	}
}

func (s *Server) initializeResponse(id any) []byte {
	return s.response(id, map[string]any{
		"protocolVersion": "2025-11-05",
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
		"serverInfo": map[string]any{
			"name":    "testmcp-server",
			"version": "1.0.0",
		},
	})
}

func (s *Server) toolsListResponse(id any) []byte {
	return s.response(id, map[string]any{
		"tools": s.tools,
	})
}

func (s *Server) toolsCallResponse(id any, params any) []byte {
	paramsMap, ok := params.(map[string]any)
	if !ok {
		return s.errorResponse(id, -32602, "Invalid params")
	}
	name, _ := paramsMap["name"].(string)

	switch name {
	case "read_db":
		return s.response(id, map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": `{"data": "test_result"}`},
			},
		})
	case "delete_db":
		return s.response(id, map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": `{"deleted": true}`},
			},
		})
	case "execute_payout":
		return s.response(id, map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": `{"status": "executed", "amount": 100}`},
			},
		})
	case "echo":
		return s.response(id, map[string]any{
			"content": []map[string]any{
				{"type": "text", "text": `{"echo": "ok"}`},
			},
		})
	default:
		return s.errorResponse(id, -32602, "Unknown tool: "+name)
	}
}

func (s *Server) response(id any, result any) []byte {
	resp := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
	data, _ := json.Marshal(resp)
	return data
}

func (s *Server) errorResponse(id any, code int, message string) []byte {
	resp := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	}
	data, _ := json.Marshal(resp)
	return data
}
