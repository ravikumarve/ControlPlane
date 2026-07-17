package audit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/matrix/mcp-guard/internal/config"
)

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"ts"`
	Identity  string    `json:"identity"`
	Tool      string    `json:"tool"`
	Params    any       `json:"params"`
	Decision  string    `json:"decision"` // allow | block | hitl | pending | denied
	Duration  int64     `json:"duration_ms,omitempty"`
	HMAC      string    `json:"hmac"`
	PrevHMAC  string    `json:"prev_hmac"`
}

// Logger writes HMAC-chained audit entries to a JSONL file.
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	hmacKey  []byte
	prevHMAC string
	enc      *json.Encoder
}

// NewLogger creates an audit logger.
func NewLogger(cfg config.AuditConfig) (*Logger, error) {
	// Ensure directory exists
	dir := dirFromPath(cfg.Path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("mkdir audit dir: %w", err)
		}
	}

	file, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open audit file: %w", err)
	}

	key := []byte(cfg.HMACKey)
	if len(key) == 0 {
		// Default key for development; production MUST set MCP_GUARD_HMAC_KEY
		key = []byte("dev-only-key-change-in-production")
	}

	l := &Logger{
		file:    file,
		hmacKey: key,
		enc:     json.NewEncoder(file),
	}

	// Read last entry's HMAC for chaining
	l.prevHMAC = l.readLastHMAC()

	return l, nil
}

// Write logs a single audit entry with HMAC chaining.
func (l *Logger) Write(entry AuditEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// HMAC chain: this entry's HMAC = HMAC(entry + prevHMAC)
	entry.PrevHMAC = l.prevHMAC
	entry.HMAC = l.computeHMAC(entry)

	if err := l.enc.Encode(entry); err != nil {
		return fmt.Errorf("encode audit entry: %w", err)
	}

	l.prevHMAC = entry.HMAC
	return nil
}

// Close flushes and closes the audit log file.
func (l *Logger) Close() error {
	return l.file.Close()
}

// computeHMAC computes the HMAC-SHA256 of an entry's content (excluding the HMAC field itself).
func (l *Logger) computeHMAC(entry AuditEntry) string {
	// Create a copy without HMAC for hashing
	entry.HMAC = ""
	data, _ := json.Marshal(entry)

	mac := hmac.New(sha256.New, l.hmacKey)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// readLastHMAC reads the last entry's HMAC from the file for chain continuation.
func (l *Logger) readLastHMAC() string {
	info, err := l.file.Stat()
	if err != nil || info.Size() == 0 {
		return ""
	}

	// Read the file from start for simplicity (in production, seek from end)
	// For MVP, we read the whole file and grab the last entry.
	data, err := os.ReadFile(l.file.Name())
	if err != nil {
		return ""
	}

	if len(data) == 0 {
		return ""
	}

	// Find last non-empty line
	lines := splitLines(data)
	for i := len(lines) - 1; i >= 0; i-- {
		if len(lines[i]) == 0 {
			continue
		}
		var entry AuditEntry
		if err := json.Unmarshal(lines[i], &entry); err == nil && entry.HMAC != "" {
			return entry.HMAC
		}
	}

	return ""
}

// splitLines splits a byte slice into lines.
func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			if i > start {
				lines = append(lines, data[start:i])
			}
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

func dirFromPath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return ""
}
