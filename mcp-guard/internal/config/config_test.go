package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func setupTestConfig(t *testing.T, yamlContent string) {
	t.Helper()
	viper.Reset()
	dir := t.TempDir()
	configPath := filepath.Join(dir, "mcp-guard.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0600); err != nil {
		t.Fatalf("Write config failed: %v", err)
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("viper.ReadInConfig failed: %v", err)
	}
}

func TestExpandEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test-value")
	os.Setenv("EMPTY_KEY", "")
	defer os.Unsetenv("TEST_KEY")
	defer os.Unsetenv("EMPTY_KEY")

	tests := []struct {
		input string
		want  string
	}{
		{"${TEST_KEY}", "test-value"},
		{"prefix_${TEST_KEY}_suffix", "prefix_test-value_suffix"},
		{"no-env-here", "no-env-here"},
		{"${NONEXISTENT}", "${NONEXISTENT}"},
		{"${EMPTY_KEY}", "${EMPTY_KEY}"},
	}

	for _, tc := range tests {
		got := expandEnv(tc.input)
		if got != tc.want {
			t.Errorf("expandEnv(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestLoad_Defaults(t *testing.T) {
	setupTestConfig(t, `
version: "1"
proxy:
  mode: stdio
audit:
  path: /tmp/test-audit.jsonl
`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Version != "1" {
		t.Errorf("Version = %q; want 1", cfg.Version)
	}
	if cfg.Proxy.Mode != "stdio" {
		t.Errorf("Proxy.Mode = %q; want stdio", cfg.Proxy.Mode)
	}
	if cfg.Audit.Path != "/tmp/test-audit.jsonl" {
		t.Errorf("Audit.Path = %q; want /tmp/test-audit.jsonl", cfg.Audit.Path)
	}
}

func TestLoad_WithPolicies(t *testing.T) {
	setupTestConfig(t, `
version: "1"
proxy:
  mode: tcp
  listen: ":8443"
  upstream: "localhost:3000"
policies:
  - name: allow-read
    action: allow
    match:
      identity: "*"
      tools: ["read_*", "get_*"]
  - name: block-dangerous
    action: block
    match:
      tools: ["drop_table"]
    alert: true
schema_pinning:
  enabled: true
  mode: warn
audit:
  path: /tmp/audit.jsonl
  hmac_key: my-secret-key
  rotation: 100MB
hitl:
  webhook_url: "${SLACK_URL}"
  timeout: 10m
`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.Policies) != 2 {
		t.Fatalf("Expected 2 policies, got %d", len(cfg.Policies))
	}

	p1 := cfg.Policies[0]
	if p1.Name != "allow-read" || p1.Action != "allow" || len(p1.Match.Tools) != 2 {
		t.Errorf("Policy 1 mismatch: %+v", p1)
	}

	p2 := cfg.Policies[1]
	if p2.Name != "block-dangerous" || p2.Action != "block" || !p2.Alert {
		t.Errorf("Policy 2 mismatch: %+v", p2)
	}

	if !cfg.SchemaPinning.Enabled {
		t.Error("SchemaPinning.Enabled should be true")
	}
	if cfg.SchemaPinning.Mode != "warn" {
		t.Errorf("SchemaPinning.Mode = %q; want warn", cfg.SchemaPinning.Mode)
	}

	if cfg.HITL == nil {
		t.Fatal("HITL config should not be nil")
	}
	if cfg.HITL.Timeout != "10m" {
		t.Errorf("HITL.Timeout = %q; want 10m", cfg.HITL.Timeout)
	}
}

func TestLoad_HITLDefaults(t *testing.T) {
	setupTestConfig(t, `version: "1"`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.HITL == nil {
		t.Fatal("HITL config should have defaults")
	}
	if cfg.HITL.Timeout != "5m" {
		t.Errorf("Default HITL.Timeout = %q; want 5m", cfg.HITL.Timeout)
	}
}

func TestLoad_MinimalConfig(t *testing.T) {
	setupTestConfig(t, `version: "1"`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Proxy.Mode != "stdio" {
		t.Errorf("Default Proxy.Mode = %q; want stdio", cfg.Proxy.Mode)
	}
	if cfg.Audit.Path != "/var/log/mcp-guard/audit.jsonl" {
		t.Errorf("Default Audit.Path = %q; want /var/log/mcp-guard/audit.jsonl", cfg.Audit.Path)
	}
	if cfg.SchemaPinning.Store != ".mcp-guard/pins.json" {
		t.Errorf("Default SchemaPinning.Store = %q; want .mcp-guard/pins.json", cfg.SchemaPinning.Store)
	}
}

// --- Validate tests ---

func TestValidate_Valid(t *testing.T) {
	cfg := &Config{
		Version: "1",
		Proxy: ProxyConfig{
			Mode:     "tcp",
			Listen:   ":8443",
			Upstream: "localhost:3000",
		},
		Policies: []PolicyCfg{
			{Name: "allow-all", Action: "allow", Match: PolicyMatch{Identity: "*"}},
		},
		Audit: AuditConfig{
			Path:    "/tmp/audit.jsonl",
			HMACKey: "secret-123",
		},
	}
	errs := cfg.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidate_InvalidMode(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "http"},
	}
	errs := cfg.Validate()
	if len(errs) == 0 {
		t.Fatal("expected errors for invalid mode")
	}
	found := false
	for _, e := range errs {
		if contains(e.Error(), "unsupported proxy mode") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'unsupported proxy mode' error, got: %v", errs)
	}
}

func TestValidate_TCPMissingFields(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "tcp"},
		Audit: AuditConfig{HMACKey: "key"},
	}
	errs := cfg.Validate()
	if len(errs) == 0 {
		t.Fatal("expected errors for tcp missing listen/upstream")
	}
}

func TestValidate_DuplicatePolicyName(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "stdio"},
		Policies: []PolicyCfg{
			{Name: "dup", Action: "allow", Match: PolicyMatch{Identity: "*"}},
			{Name: "dup", Action: "block", Match: PolicyMatch{Tools: []string{"x"}}},
		},
		Audit: AuditConfig{HMACKey: "key"},
	}
	errs := cfg.Validate()
	found := false
	for _, e := range errs {
		if contains(e.Error(), "duplicate policy name") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected duplicate policy name error, got: %v", errs)
	}
}

func TestValidate_InvalidAction(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "stdio"},
		Policies: []PolicyCfg{
			{Name: "bad", Action: "delete", Match: PolicyMatch{Identity: "*"}},
		},
		Audit: AuditConfig{HMACKey: "key"},
	}
	errs := cfg.Validate()
	found := false
	for _, e := range errs {
		if contains(e.Error(), "action") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected action error, got: %v", errs)
	}
}

func TestValidate_InvalidRateLimit(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "stdio"},
		Policies: []PolicyCfg{
			{Name: "rl", Action: "allow", Match: PolicyMatch{Identity: "*"}, RateLimit: "xyz"},
		},
		Audit: AuditConfig{HMACKey: "key"},
	}
	errs := cfg.Validate()
	found := false
	for _, e := range errs {
		if contains(e.Error(), "rate_limit") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected rate_limit error, got: %v", errs)
	}
}

func TestValidate_NoMatch(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "stdio"},
		Policies: []PolicyCfg{
			{Name: "empty", Action: "allow"},
		},
		Audit: AuditConfig{HMACKey: "key"},
	}
	errs := cfg.Validate()
	found := false
	for _, e := range errs {
		if contains(e.Error(), "no match criteria") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'no match criteria' error, got: %v", errs)
	}
}

func TestValidate_EmptyHMACKey(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "stdio"},
		Policies: []PolicyCfg{
			{Name: "p", Action: "allow", Match: PolicyMatch{Identity: "*"}},
		},
		Audit: AuditConfig{Path: "/tmp/audit.jsonl"},
	}
	errs := cfg.Validate()
	found := false
	for _, e := range errs {
		if contains(e.Error(), "hmac_key is empty") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected hmac_key warning, got: %v", errs)
	}
}

func TestValidate_BadHITLWebhook(t *testing.T) {
	cfg := &Config{
		Proxy: ProxyConfig{Mode: "stdio"},
		Policies: []PolicyCfg{
			{Name: "p", Action: "allow", Match: PolicyMatch{Identity: "*"}},
		},
		Audit: AuditConfig{HMACKey: "key"},
		HITL:  &HITLConfig{WebhookURL: "not-a-url"},
	}
	errs := cfg.Validate()
	found := false
	for _, e := range errs {
		if contains(e.Error(), "webhook_url") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected webhook_url error, got: %v", errs)
	}
}

func TestParseRateSimple(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"100/m", false},
		{"10/s", false},
		{"1000/h", false},
		{"", true},
		{"abc", true},
		{"100/x", true},
	}
	for _, tt := range tests {
		_, _, err := parseRateSimple(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("expected error for %q", tt.input)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("unexpected error for %q: %v", tt.input, err)
		}
	}
}

func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
