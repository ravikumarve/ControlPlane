package policy

import (
	"github.com/matrix/mcp-guard/internal/config"
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

// Engine evaluates tool calls against configured policies.
type Engine struct {
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

// Evaluate checks a tool call against all policies in order.
// First matching policy wins (deny by default).
func (e *Engine) Evaluate(identity, tool string) Decision {
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
