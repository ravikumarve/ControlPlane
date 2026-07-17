package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config is the root configuration structure.
type Config struct {
	Version       string       `mapstructure:"version"`
	Proxy         ProxyConfig  `mapstructure:"proxy"`
	Policies      []PolicyCfg  `mapstructure:"policies"`
	SchemaPinning SchemaConfig `mapstructure:"schema_pinning"`
	Audit         AuditConfig  `mapstructure:"audit"`
	HITL          *HITLConfig  `mapstructure:"hitl"`
}

type ProxyConfig struct {
	Mode     string `mapstructure:"mode"`
	Listen   string `mapstructure:"listen"`
	Upstream string `mapstructure:"upstream"`
}

type PolicyCfg struct {
	Name        string       `mapstructure:"name"`
	Match       PolicyMatch  `mapstructure:"match"`
	Action      string       `mapstructure:"action"`
	Constraints *Constraints `mapstructure:"constraints,omitempty"`
	RateLimit   string       `mapstructure:"rate_limit,omitempty"`
	Alert       bool         `mapstructure:"alert,omitempty"`
}

type PolicyMatch struct {
	Identity string   `mapstructure:"identity"`
	Tools    []string `mapstructure:"tools"`
}

type Constraints struct {
	MaxAmount float64 `mapstructure:"max_amount,omitempty"`
}

type SchemaConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Mode    string `mapstructure:"mode"`
	Store   string `mapstructure:"store"`
}

type AuditConfig struct {
	Path     string `mapstructure:"path"`
	HMACKey  string `mapstructure:"hmac_key"`
	Rotation string `mapstructure:"rotation"`
}

type HITLConfig struct {
	WebhookURL string   `mapstructure:"webhook_url"`
	Timeout    string   `mapstructure:"timeout"`
	Channels   []string `mapstructure:"channels"`
}

// Load reads configuration from viper and expands env vars.
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}

	// Expand ${VAR} patterns in string fields
	cfg.Audit.HMACKey = expandEnv(cfg.Audit.HMACKey)
	if cfg.HITL != nil {
		cfg.HITL.WebhookURL = expandEnv(cfg.HITL.WebhookURL)
	}

	// Defaults
	if cfg.Proxy.Mode == "" {
		cfg.Proxy.Mode = "stdio"
	}
	if cfg.Audit.Path == "" {
		cfg.Audit.Path = "/var/log/mcp-guard/audit.jsonl"
	}
	if cfg.SchemaPinning.Store == "" {
		cfg.SchemaPinning.Store = ".mcp-guard/pins.json"
	}

	// Set defaults for HITL
	if cfg.HITL == nil {
		cfg.HITL = &HITLConfig{}
	}
	if cfg.HITL.Timeout == "" {
		cfg.HITL.Timeout = "5m"
	}

	return &cfg, nil
}

// Validate checks the configuration for errors and returns all issues found.
func (c *Config) Validate() []error {
	var errs []error

	// Proxy mode
	switch c.Proxy.Mode {
	case "stdio", "tcp":
		// valid
	case "":
		errs = append(errs, fmt.Errorf("proxy.mode is required (stdio or tcp)"))
	default:
		errs = append(errs, fmt.Errorf("unsupported proxy mode: %q (use stdio or tcp)", c.Proxy.Mode))
	}

	// TCP mode requires listen + upstream
	if c.Proxy.Mode == "tcp" {
		if c.Proxy.Listen == "" {
			errs = append(errs, fmt.Errorf("proxy.listen is required in tcp mode"))
		}
		if c.Proxy.Upstream == "" {
			errs = append(errs, fmt.Errorf("proxy.upstream is required in tcp mode"))
		}
	}

	// Policy validation
	seenNames := make(map[string]int)
	for i, p := range c.Policies {
		if p.Name == "" {
			errs = append(errs, fmt.Errorf("policies[%d].name is required", i))
		} else if j, ok := seenNames[p.Name]; ok {
			errs = append(errs, fmt.Errorf("duplicate policy name %q at indices %d and %d", p.Name, j, i))
		} else {
			seenNames[p.Name] = i
		}

		switch p.Action {
		case "allow", "block", "hitl":
			// valid
		case "":
			errs = append(errs, fmt.Errorf("policies[%d].action is required (allow/block/hitl)", i))
		default:
			errs = append(errs, fmt.Errorf("policies[%d].action %q is invalid (use allow/block/hitl)", i, p.Action))
		}

		// Validate rate limit string
		if p.RateLimit != "" {
			_, _, err := parseRateSimple(p.RateLimit)
			if err != nil {
				errs = append(errs, fmt.Errorf("policies[%d].rate_limit %q: %v", i, p.RateLimit, err))
			}
		}

		// Identity wildcard should be explicit
		if p.Match.Identity == "" && len(p.Match.Tools) == 0 {
			errs = append(errs, fmt.Errorf("policies[%d] has no match criteria (empty identity and empty tools)", i))
		}
	}

	// Audit HMAC key warning
	if c.Audit.HMACKey == "" {
		errs = append(errs, fmt.Errorf("audit.hmac_key is empty — audit log will not be tamper-proof"))
	}

	// HITL config validation
	if c.HITL != nil && c.HITL.WebhookURL != "" {
		if !strings.HasPrefix(c.HITL.WebhookURL, "http://") && !strings.HasPrefix(c.HITL.WebhookURL, "https://") {
			errs = append(errs, fmt.Errorf("hitl.webhook_url must start with http:// or https://"))
		}
	}

	return errs
}

// expandEnv replaces ${VAR} or $VAR patterns with environment variable values.
func expandEnv(s string) string {
	if !strings.Contains(s, "${") && !strings.Contains(s, "$") {
		return s
	}
	return os.Expand(s, func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			return "${" + key + "}"
		}
		return val
	})
}

// parseRateSimple does a lightweight check on rate limit format without importing ratelimit.
func parseRateSimple(s string) (int, string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, "", fmt.Errorf("empty rate string")
	}
	split := -1
	for i, r := range s {
		if r < '0' || r > '9' {
			split = i
			break
		}
	}
	if split <= 0 || split >= len(s) {
		return 0, "", fmt.Errorf("invalid format (expected e.g. 100/m)")
	}
	unit := s[split:]
	switch unit {
	case "/s", "/sec", "/second", "/m", "/min", "/minute", "/h", "/hr", "/hour":
		return 0, unit, nil
	default:
		return 0, "", fmt.Errorf("unsupported unit %q (use /s, /m, or /h)", unit)
	}
}
