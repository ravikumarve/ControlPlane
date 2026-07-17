package policy

import (
	"strings"

	"github.com/gobwas/glob"
)

// Matcher provides glob-based identity and tool matching.
type Matcher struct {
	cache map[string]glob.Glob
}

// NewMatcher creates a new matcher with a compiled glob cache.
func NewMatcher() *Matcher {
	return &Matcher{
		cache: make(map[string]glob.Glob),
	}
}

// Match checks if a value matches a glob pattern.
// Supports * (any sequence) and ? (single char).
func (m *Matcher) Match(pattern, value string) bool {
	// Fast path for exact match
	if pattern == value {
		return true
	}

	// Fast path for wildcard
	if pattern == "*" {
		return true
	}

	// Try cache
	g, ok := m.cache[pattern]
	if !ok {
		var err error
		g, err = glob.Compile(pattern)
		if err != nil {
			// Fall back to simple strings.Contains
			m.cache[pattern] = nil
			return strings.Contains(value, strings.Trim(pattern, "*"))
		}
		m.cache[pattern] = g
	}

	return g.Match(value)
}

// MatchAny checks if a value matches any of the given patterns.
func (m *Matcher) MatchAny(patterns []string, value string) bool {
	for _, p := range patterns {
		if m.Match(p, value) {
			return true
		}
	}
	return false
}
