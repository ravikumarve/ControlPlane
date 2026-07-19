package inject

// PromptPattern is a known prompt injection pattern to detect.
type PromptPattern struct {
	ID          string
	Description string
	Patterns    []string // substrings to match (lowercased)
	Severity    float64  // 0.0–1.0
}

// KnownPromptPatterns is the built-in registry of prompt injection signatures.
var KnownPromptPatterns = []PromptPattern{
	{
		ID:          "INJ-IGNORE-PREV",
		Description: "Attempt to ignore previous instructions",
		Patterns:    []string{"ignore previous", "ignore all previous", "disregard previous", "forget previous", "override previous"},
		Severity:    0.9,
	},
	{
		ID:          "INJ-NEW-SYSTEM",
		Description: "Attempt to redefine system persona",
		Patterns:    []string{"you are now", "act as", "pretend to be", "from now on you are", "your new role is", "you are an ai", "respond as"},
		Severity:    0.7,
	},
	{
		ID:          "INJ-PROMPT-LEAK",
		Description: "Attempt to leak system prompt",
		Patterns:    []string{"print your", "output your", "reveal your", "show your prompt", "what are your instructions", "system prompt", "initial prompt", "what instructions"},
		Severity:    0.8,
	},
	{
		ID:          "INJ-PROMPT-EXTRACT",
		Description: "Attempt to extract instructions via encoding",
		Patterns:    []string{"encoded", "base64", "rot13", "hex encoded", "translate to", "repeat after", "repeat every"},
		Severity:    0.6,
	},
	{
		ID:          "INJ-JAILBREAK",
		Description: "Known jailbreak prefixes",
		Patterns:    []string{"dan ", "jailbreak", "freedom", "unfiltered", "no restrictions", "uncensored", "unconstrained"},
		Severity:    0.9,
	},
	{
		ID:          "INJ-ROLE-HIJACK",
		Description: "Attempt to hijack assistant role",
		Patterns:    []string{"now you are", "you will act", "you will pretend", "you must roleplay", "forget your role", "stop being", "don't be"},
		Severity:    0.8,
	},
	{
		ID:          "INJ-SQL",
		Description: "SQL injection attempt",
		Patterns:    []string{"' or '1'='1", "' or 1=1", "drop table", "delete from", "union select", ";--", "'--"},
		Severity:    0.9,
	},
	{
		ID:          "INJ-COMMAND",
		Description: "Shell command injection attempt",
		Patterns:    []string{"; rm ", "; rm -rf ", "| bash", "`cat ", "$(cat ", "pipe to", "2>&1", "/dev/null"},
		Severity:    0.9,
	},
	{
		ID:          "INJ-DATA-EXFIL",
		Description: "Data exfiltration patterns",
		Patterns:    []string{"send to", "post to", "upload to", "exfiltrate", "fetch from", "http://", "https://", "curl ", "wget "},
		Severity:    0.6,
	},
	{
		ID:          "INJ-CONFABULATION",
		Description: "Attempt to force fabricated data",
		Patterns:    []string{"pretend", "imagine", "hypothetically", "falsify", "make up", "lie about", "hallucinate"},
		Severity:    0.5,
	},
}
