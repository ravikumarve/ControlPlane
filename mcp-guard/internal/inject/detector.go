package inject

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Result holds the outcome of an injection scan.
type Result struct {
	Detected     bool
	Severity     float64 // 0.0–1.0
	Reason       string
	PatternID    string
	Confusables  []string // list of confusable runs found
}

// String returns a human-readable summary.
func (r Result) String() string {
	if !r.Detected {
		return "clean"
	}
	return fmt.Sprintf("%s (severity: %.1f)", r.Reason, r.Severity)
}

// ScanResult is returned by the full pipeline.
type ScanResult struct {
	Injection Result
	// Future: schema drift, rate limit, etc.
}

// Config configures the injection detector.
type Config struct {
	// MaxParamDepth rejects JSON deeper than this (bomb protection).
	MaxParamDepth int
	// MaxParamLength rejects string values longer than this.
	MaxParamLength int
	// MinConfusableScore is the fraction of confusable chars to flag (0.0–1.0).
	MinConfusableScore float64
	// EnabledScans controls which scans to run.
	EnabledScans []string // "homoglyph", "patterns", "depth", "length"
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		MaxParamDepth:      20,
		MaxParamLength:     100000,
		MinConfusableScore: 0.3,
		EnabledScans:       []string{"homoglyph", "patterns", "depth", "length"},
	}
}

// Detector performs injection detection on tool call parameters.
type Detector struct {
	cfg Config
}

// NewDetector creates a detector with the given config.
func NewDetector(cfg Config) *Detector {
	return &Detector{cfg: cfg}
}

// ScanParams runs the full injection pipeline on tool call parameters.
// The params argument is the raw JSON params from the tool call.
func (d *Detector) ScanParams(toolName string, params any) ScanResult {
	var sr ScanResult

	if params == nil {
		return sr
	}

	// Convert params to JSON string for scanning
	raw, err := json.Marshal(params)
	if err != nil || len(raw) == 0 {
		return sr
	}

	paramStr := string(raw)

	// Depth check (JSON bomb protection)
	if d.isEnabled("depth") {
		depth := jsonDepth(paramStr)
		if depth > d.cfg.MaxParamDepth {
			sr.Injection = Result{
				Detected: true,
				Severity: 0.8,
				Reason:   fmt.Sprintf("excessive JSON depth: %d > %d", depth, d.cfg.MaxParamDepth),
			}
			return sr // short-circuit — depth bomb is high severity
		}
	}

	// Length check
	if d.isEnabled("length") {
		if len(paramStr) > d.cfg.MaxParamLength {
			sr.Injection = Result{
				Detected: true,
				Severity: 0.7,
				Reason:   fmt.Sprintf("param payload too large: %d bytes", len(paramStr)),
			}
			return sr
		}
	}

	// Flatten params to searchable text (lowercase, key:value pairs)
	flat := flattenParams(params)
	lower := strings.ToLower(flat)

	// Pattern-based injection scan
	if d.isEnabled("patterns") {
		if r := d.scanPatterns(lower); r.Detected {
			sr.Injection = r
			return sr
		}
	}

	// Homoglyph / confusable scan (cross-script character substitution)
	if d.isEnabled("homoglyph") {
		if r := d.scanConfusables(paramStr); r.Detected {
			sr.Injection = r
			return sr
		}
	}

	return sr
}

// scanPatterns checks known prompt injection patterns.
func (d *Detector) scanPatterns(lowerText string) Result {
	for _, pat := range KnownPromptPatterns {
		for _, pattern := range pat.Patterns {
			if strings.Contains(lowerText, pattern) {
				return Result{
					Detected:  true,
					Severity:  pat.Severity,
					Reason:    fmt.Sprintf("prompt injection detected: %s", pat.Description),
					PatternID: pat.ID,
				}
			}
		}
	}
	return Result{}
}

// scanConfusables detects homoglyph / confusable Unicode characters
// using a cross-script confusable character map (Cyrillic, Greek, etc.).
func (d *Detector) scanConfusables(text string) Result {
	var confusables []string

	for i, r := range text {
		if ascii, ok := isConfusableRune(r); ok {
			// Found a confusable — report the original text around it
			context := extractContext(text, i)
			confusables = append(confusables, fmt.Sprintf("'%c'(U+%04X)→'%c' at %s", r, r, ascii, context))
		}
	}

	if len(confusables) == 0 {
		return Result{}
	}

	return Result{
		Detected:    true,
		Severity:    0.7,
		Reason:      fmt.Sprintf("homoglyph/confusable characters detected (%d occurrences)", len(confusables)),
		Confusables: confusables,
	}
}

// isEnabled checks if a scan type is enabled.
func (d *Detector) isEnabled(name string) bool {
	for _, s := range d.cfg.EnabledScans {
		if s == name {
			return true
		}
	}
	return false
}

// flattenParams converts tool params to a flat searchable text string.
func flattenParams(v any) string {
	var parts []string
	flatten("", v, &parts, 0)
	return strings.Join(parts, " ")
}

func flatten(prefix string, v any, parts *[]string, depth int) {
	if depth > 100 {
		return
	}
	switch val := v.(type) {
	case map[string]any:
		for k, vv := range val {
			flatten(k, vv, parts, depth+1)
		}
	case []any:
		for _, vv := range val {
			flatten(prefix, vv, parts, depth+1)
		}
	case string:
		*parts = append(*parts, val)
	case float64, bool:
		*parts = append(*parts, fmt.Sprintf("%v", val))
	}
}

// jsonDepth computes the maximum nesting depth of a JSON string.
func jsonDepth(s string) int {
	var maxDepth, depth int
	for _, r := range s {
		switch r {
		case '{', '[':
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		case '}', ']':
			depth--
		}
	}
	return maxDepth
}

// extractContext returns a short context window around position i.
func extractContext(s string, i int) string {
	start := i - 10
	if start < 0 {
		start = 0
	}
	end := i + 10
	if end > len(s) {
		end = len(s)
	}
	return fmt.Sprintf("...%s...", s[start:end])
}


