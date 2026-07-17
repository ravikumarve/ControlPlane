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

	// Version comes from config
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
