package policy

import (
	"testing"

	"github.com/matrix/mcp-guard/internal/config"
)

func TestEngine_Evaluate_Allow(t *testing.T) {
	engine := NewEngine([]config.PolicyCfg{
		{
			Name:   "allow-read",
			Action: ActionAllow,
			Match: config.PolicyMatch{
				Identity: "*",
				Tools:    []string{"read_*", "get_*", "search_*"},
			},
		},
	})

	tests := []struct {
		identity string
		tool     string
		want     string
	}{
		{"agent1", "read_database", ActionAllow},
		{"agent2", "get_user", ActionAllow},
		{"anyone", "search_records", ActionAllow},
		{"agent1", "delete_user", ActionBlock}, // blocked by default
	}

	for _, tc := range tests {
		got := engine.Evaluate(tc.identity, tc.tool)
		if got.Action != tc.want {
			t.Errorf("Evaluate(%q, %q) = %q; want %q (policy: %s)",
				tc.identity, tc.tool, got.Action, tc.want, got.PolicyName)
		}
	}
}

func TestEngine_Evaluate_Block(t *testing.T) {
	engine := NewEngine([]config.PolicyCfg{
		{
			Name:   "block-dangerous",
			Action: ActionBlock,
			Match: config.PolicyMatch{
				Tools: []string{"drop_table", "rm_rf", "exec_shell"},
			},
			Alert: true,
		},
	})

	tests := []struct {
		tool string
		want string
	}{
		{"drop_table", ActionBlock},
		{"rm_rf", ActionBlock},
		{"exec_shell", ActionBlock},
		{"read_database", ActionBlock}, // blocked by default (no matching allow)
	}

	for _, tc := range tests {
		got := engine.Evaluate("any", tc.tool)
		if got.Action != tc.want {
			t.Errorf("Evaluate(_, %q) = %q; want %q (policy: %s, reason: %s)",
				tc.tool, got.Action, tc.want, got.PolicyName, got.Reason)
		}
	}
}

func TestEngine_Evaluate_HITL(t *testing.T) {
	engine := NewEngine([]config.PolicyCfg{
		{
			Name:   "payment-approval",
			Action: ActionHITL,
			Match: config.PolicyMatch{
				Tools: []string{"execute_payout", "refund", "transfer"},
			},
			Constraints: &config.Constraints{
				MaxAmount: 500,
			},
		},
	})

	tests := []struct {
		tool        string
		wantAction  string
		wantRisk    float64
	}{
		{"execute_payout", ActionHITL, 0.8}, // has constraint → higher risk
		{"refund", ActionHITL, 0.8},
		{"read_user", ActionBlock, 0}, // not matched → blocked by default
	}

	for _, tc := range tests {
		got := engine.Evaluate("agent", tc.tool)
		if got.Action != tc.wantAction {
			t.Errorf("Evaluate(_, %q) = %q; want %q", tc.tool, got.Action, tc.wantAction)
		}
		if tc.wantRisk > 0 && got.RiskScore != tc.wantRisk {
			t.Errorf("Evaluate(_, %q) RiskScore = %f; want %f", tc.tool, got.RiskScore, tc.wantRisk)
		}
	}
}

func TestEngine_Evaluate_FirstMatchWins(t *testing.T) {
	engine := NewEngine([]config.PolicyCfg{
		{
			Name:   "block-first",
			Action: ActionBlock,
			Match: config.PolicyMatch{
				Tools: []string{"*"},
			},
		},
		{
			Name:   "allow-second",
			Action: ActionAllow,
			Match: config.PolicyMatch{
				Tools: []string{"read_*"},
			},
		},
	})

	// First policy matches everything with block, so even "read_*" should be blocked
	got := engine.Evaluate("agent", "read_database")
	if got.Action != ActionBlock {
		t.Errorf("First-match-wins: got %q; want %q (matched: %s)", got.Action, ActionBlock, got.PolicyName)
	}
	if got.PolicyName != "block-first" {
		t.Errorf("Expected policy 'block-first', got %q", got.PolicyName)
	}
}

func TestEngine_Evaluate_IdentityMatch(t *testing.T) {
	engine := NewEngine([]config.PolicyCfg{
		{
			Name:   "payment-agent",
			Action: ActionAllow,
			Match: config.PolicyMatch{
				Identity: "payment-bot",
				Tools:    []string{"execute_payout"},
			},
		},
		{
			Name:   "block-payment",
			Action: ActionBlock,
			Match: config.PolicyMatch{
				Tools: []string{"execute_payout"},
			},
		},
	})

	// payment-bot should be allowed
	got := engine.Evaluate("payment-bot", "execute_payout")
	if got.Action != ActionAllow {
		t.Errorf("payment-bot should be allowed; got %q", got.Action)
	}

	// other agents should be blocked by the second policy
	got = engine.Evaluate("random-agent", "execute_payout")
	if got.Action != ActionBlock {
		t.Errorf("random-agent should be blocked; got %q", got.Action)
	}
}

func TestEngine_Evaluate_DefaultDeny(t *testing.T) {
	engine := NewEngine([]config.PolicyCfg{})

	got := engine.Evaluate("any", "anything")
	if got.Action != ActionBlock {
		t.Errorf("Empty policies should default-deny; got %q", got.Action)
	}
	if got.Reason == "" {
		t.Errorf("Block reason should not be empty")
	}
}

func TestGlobPatterns(t *testing.T) {
	engine := NewEngine([]config.PolicyCfg{
		{
			Name:   "all-db",
			Action: ActionAllow,
			Match: config.PolicyMatch{
				Tools: []string{"db_*", "query_*"},
			},
		},
	})

	tests := []struct {
		tool string
		want string
	}{
		{"db_select", ActionAllow},
		{"db_insert", ActionAllow},
		{"db_delete", ActionAllow},
		{"query_users", ActionAllow},
		{"http_get", ActionBlock},
	}

	for _, tc := range tests {
		got := engine.Evaluate("agent", tc.tool)
		if got.Action != tc.want {
			t.Errorf("Evaluate(_, %q) = %q; want %q", tc.tool, got.Action, tc.want)
		}
	}
}

func TestMatcher_Glob(t *testing.T) {
	m := NewMatcher()

	tests := []struct {
		pattern string
		value   string
		want    bool
	}{
		{"*", "anything", true},
		{"read_*", "read_database", true},
		{"read_*", "write_database", false},
		{"*_db", "read_db", true},
		{"*_db", "read_database", false},
		{"exact", "exact", true},
		{"exact", "not-exact", false},
		{"agent-?", "agent-1", true},
		{"agent-?", "agent-12", false},
	}

	for _, tc := range tests {
		got := m.Match(tc.pattern, tc.value)
		if got != tc.want {
			t.Errorf("Match(%q, %q) = %v; want %v", tc.pattern, tc.value, got, tc.want)
		}
	}

	// Test cache hit
	got := m.Match("read_*", "read_something")
	if !got {
		t.Error("Cache: Match should return true for cached pattern")
	}
}

func TestMatcher_MatchAny(t *testing.T) {
	m := NewMatcher()

	got := m.MatchAny([]string{"read_*", "get_*", "search_*"}, "read_users")
	if !got {
		t.Error("MatchAny should match 'read_users' against 'read_*'")
	}

	got = m.MatchAny([]string{"read_*", "get_*"}, "delete_users")
	if got {
		t.Error("MatchAny should NOT match 'delete_users'")
	}

	got = m.MatchAny([]string{}, "anything")
	if got {
		t.Error("MatchAny with empty patterns should return false")
	}
}
