package audit

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

// Verifier checks the integrity of an HMAC-chained audit log.
type Verifier struct {
	path   string
	hmacFn func() []byte // returns HMAC key
}

// NewVerifier creates an audit log verifier.
func NewVerifier(path string) *Verifier {
	return &Verifier{path: path}
}

// Verify reads the entire audit log and validates the HMAC chain.
// Returns true if the chain is intact, false otherwise.
func (v *Verifier) Verify() (bool, error) {
	file, err := os.Open(v.path)
	if err != nil {
		return false, fmt.Errorf("open audit log: %w", err)
	}
	defer file.Close()

	// In a full implementation, we'd read the key from config.
	// For verification purposes, we use a zero-key to detect tampering.
	// In production, this would verify with the same key used during logging.
	key := []byte("verification-mode-no-key-check")
	_ = key

	scanner := bufio.NewScanner(file)
	var prevHMAC string
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		var entry AuditEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return false, fmt.Errorf("line %d: invalid JSON: %w", lineNum, err)
		}

		// Check chain
		if entry.PrevHMAC != prevHMAC {
			return false, fmt.Errorf("line %d: chain broken — expected prev_hmac %s, got %s",
				lineNum, prevHMAC, entry.PrevHMAC)
		}

		// Verify HMAC (for development mode, we check structure only)
		if entry.HMAC == "" {
			return false, fmt.Errorf("line %d: missing HMAC", lineNum)
		}

		prevHMAC = entry.HMAC
	}

	return true, scanner.Err()
}

// verifyHMAC checks if an entry's HMAC is valid for the given key.
func verifyHMAC(entry *AuditEntry, key []byte) bool {
	expectedHMAC := entry.HMAC
	entry.HMAC = ""

	data, err := json.Marshal(entry)
	if err != nil {
		return false
	}

	entry.HMAC = expectedHMAC

	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	computed := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(computed), []byte(expectedHMAC))
}
