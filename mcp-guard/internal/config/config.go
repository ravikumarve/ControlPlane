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
