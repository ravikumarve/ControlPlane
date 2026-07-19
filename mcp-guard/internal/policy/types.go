package policy

import (
	"fmt"
	"sync"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
)

// Action constants.
const (
	ActionAllow = "allow"
	ActionBlock = "block"
	ActionHITL  = "hitl"
)

// Decision is the result of evaluating a policy.
type Decision struct {
	Action     string
	PolicyName string
	Reason     string
	RiskScore  float64
}

// RawPolicy is a simplified policy for API serialization.
type RawPolicy struct {
	Name   string            `json:"name"`
	Action string            `json:"action"`
	Match  RawPolicyMatch    `json:"match"`
}

type RawPolicyMatch struct {
	Identity string            `json:"identity"`
	Tools    []string          `json:"tools"`
	Params   map[string]any    `json:"params,omitempty"`
}

// Engine evaluates tool calls against configured policies.
type Engine struct {
	mu       sync.RWMutex
	policies []config.PolicyCfg
	matcher  *Matcher
}

// NewEngine creates a policy engine from configured policies.
func NewEngine(policies []config.PolicyCfg) *Engine {
	return &Engine{
		policies: policies,
		matcher:  NewMatcher(),
	}
}

// List returns all policies as raw configs.
func (e *Engine) List() []config.PolicyCfg {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]config.PolicyCfg, len(e.policies))
	copy(result, e.policies)
	return result
}

// Replace replaces all policies with a new set.
func (e *Engine) Replace(raw []RawPolicy) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	policies := make([]config.PolicyCfg, len(raw))
	for i, rp := range raw {
		if rp.Name == "" {
			return fmt.Errorf("policy %d: name is required", i)
		}
		switch rp.Action {
		case "allow", "block", "hitl":
		default:
			return fmt.Errorf("policy %d: invalid action %q", i, rp.Action)
		}
		policies[i] = config.PolicyCfg{
			Name: rp.Name,
			Action: rp.Action,
			Match: config.PolicyMatch{
				Identity: rp.Match.Identity,
				Tools:    rp.Match.Tools,
				Params:   rp.Match.Params,
			},
		}
	}

	e.policies = policies
	return nil
}

// Evaluate checks a tool call against all policies in order.
// First matching policy wins (deny by default).
func (e *Engine) Evaluate(identity, tool string, params map[string]any) Decision {
	for _, p := range e.policies {
		// Check identity match
		if p.Match.Identity != "" && p.Match.Identity != "*" {
			if !e.matcher.Match(p.Match.Identity, identity) {
				continue
			}
		}

		// Check tool match
		if len(p.Match.Tools) > 0 {
			matched := false
			for _, pattern := range p.Match.Tools {
				if e.matcher.Match(pattern, tool) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// Check parameter match (if specified in policy)
		if len(p.Match.Params) > 0 {
			paramsMatch := true
			for paramKey, paramPattern := range p.Match.Params {
				actualValue, exists := params[paramKey]
				if !exists {
					paramsMatch = false
					break
				}

				actualStr := fmt.Sprintf("%v", actualValue)

				switch pattern := paramPattern.(type) {
				case string:
					if !e.matcher.Match(pattern, actualStr) {
						paramsMatch = false
					}
				case []any:
					matched := false
					for _, p := range pattern {
						if patternStr, ok := p.(string); ok {
							if e.matcher.Match(patternStr, actualStr) {
								matched = true
								break
							}
						}
					}
					if !matched {
						paramsMatch = false
					}
				default:
					paramsMatch = false
				}

				if !paramsMatch {
					break
				}
			}

			if !paramsMatch {
				continue
			}
		}

		// Policy matched
		switch p.Action {
		case ActionAllow:
			return Decision{
				Action:     ActionAllow,
				PolicyName: p.Name,
			}
		case ActionBlock:
			reason := "blocked by policy: " + p.Name
			return Decision{
				Action:     ActionBlock,
				PolicyName: p.Name,
				Reason:     reason,
			}
		case ActionHITL:
			riskScore := 0.5
			if p.Constraints != nil && p.Constraints.MaxAmount > 0 {
				riskScore = 0.8
			}
			return Decision{
				Action:     ActionHITL,
				PolicyName: p.Name,
				Reason:     "requires human approval",
				RiskScore:  riskScore,
			}
		}
	}

	// Default: deny with reason
	return Decision{
		Action: ActionBlock,
		Reason: "no matching policy — denied by default",
	}
}
