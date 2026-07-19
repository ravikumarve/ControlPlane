package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
)

// Pin represents a stored schema pin for an MCP server.
type Pin struct {
	ServerURL   string            `json:"server_url"`
	ToolHashes  map[string]string `json:"tool_hashes"`
	PinnedAt    time.Time         `json:"pinned_at"`
	LastChecked time.Time         `json:"last_checked"`
}

// DriftReport describes a schema change detected during verification.
type DriftReport struct {
	Server string
	Tool   string
	Old    string
	New    string
}

// Pinner manages schema pinning for MCP server tool definitions.
// It hashes tool definitions and detects supply-chain poisoning
// when servers change their tool schemas after initial registration.
type Pinner struct {
	cfg    config.SchemaConfig
	mu     sync.RWMutex
	pins   map[string]*Pin // serverURL -> pin
	hashes map[string]string // toolName -> currentHash cache
}

// NewPinner creates a schema pinner.
func NewPinner(cfg config.SchemaConfig) *Pinner {
	return &Pinner{
		cfg:    cfg,
		pins:   make(map[string]*Pin),
		hashes: make(map[string]string),
	}
}

// HashTool computes a SHA-256 hash of a tool's definition JSON.
func (p *Pinner) HashTool(name string, definition any) string {
	data, _ := json.Marshal(map[string]any{
		"name":        name,
		"definition":  definition,
	})
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// CheckAndPin checks a tool call against pinned schemas and updates hashes.
// This should be called for every allowed tools/list response.
func (p *Pinner) CheckAndPin(toolName string, rawData []byte) {
	// For MVP, we hash tool names as seen in requests.
	// In production, we'd parse tools/list responses.
	hash := sha256.Sum256(rawData)
	hashStr := hex.EncodeToString(hash[:])

	p.mu.Lock()
	defer p.mu.Unlock()

	// Check for drift
	if existing, ok := p.hashes[toolName]; ok {
		if existing != hashStr && p.cfg.Mode == "block" {
			log.Warn().
				Str("tool", toolName).
				Str("old_hash", existing).
				Str("new_hash", hashStr).
				Msg("SCHEMA DRIFT DETECTED — tool definition changed")
		}
	}

	p.hashes[toolName] = hashStr
}

// LoadPins reads the stored pins from disk.
func (p *Pinner) LoadPins() ([]Pin, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	data, err := os.ReadFile(p.cfg.Store)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var pins []Pin
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, err
	}

	// Rebuild map
	for i := range pins {
		p.pins[pins[i].ServerURL] = &pins[i]
	}

	return pins, nil
}

// SavePins writes the current pins to disk.
func (p *Pinner) SavePins(pins []Pin) error {
	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return err
	}

	dir := dirFromPath(p.cfg.Store)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir pin store: %w", err)
		}
	}

	return os.WriteFile(p.cfg.Store, data, 0600)
}

// VerifyAll checks all pinned schemas for drift.
// Returns a list of drifts found.
func (p *Pinner) VerifyAll() ([]DriftReport, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.hashes) == 0 {
		return nil, fmt.Errorf("no hashes loaded — connect to servers first")
	}

	// For MVP, we compare the current in-memory hashes with stored pins
	pins, err := p.LoadPins()
	if err != nil {
		return nil, err
	}

	var drifts []DriftReport
	for _, pin := range pins {
		for tool, storedHash := range pin.ToolHashes {
			if currentHash, ok := p.hashes[tool]; ok {
				if storedHash != currentHash {
					drifts = append(drifts, DriftReport{
						Server: pin.ServerURL,
						Tool:   tool,
						Old:    storedHash,
						New:    currentHash,
					})
				}
			}
		}
	}

	return drifts, nil
}

func dirFromPath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return ""
}
