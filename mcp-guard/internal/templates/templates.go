package templates

import (
	_ "embed"
	"strings"
)

//go:embed github.yaml
var githubYAML string

//go:embed postgres.yaml
var postgresYAML string

//go:embed slack.yaml
var slackYAML string

//go:embed filesystem.yaml
var filesystemYAML string

//go:embed generic.yaml
var genericYAML string

var templates = map[string]string{
	"github":     githubYAML,
	"postgres":   postgresYAML,
	"slack":      slackYAML,
	"filesystem": filesystemYAML,
	"generic":    genericYAML,
}

func Get(name string) (string, bool) {
	content, ok := templates[strings.ToLower(name)]
	return content, ok
}

func List() []string {
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	return names
}
