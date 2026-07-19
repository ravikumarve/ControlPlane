package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Event represents an alert event to dispatch.
type Event struct {
	Type      string    // "block", "hitl", "rate_limit", "injection"
	Tool      string
	Identity  string
	Reason    string
	Timestamp time.Time
	Policy    string
}

// Dispatcher sends alerts to configured webhook URLs.
type Dispatcher struct {
	webhookURL string
	channel    string // slack, discord, generic
	client     *http.Client
}

// NewDispatcher creates an alert dispatcher.
func NewDispatcher(webhookURL, channel string) *Dispatcher {
	return &Dispatcher{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Send dispatches an alert event to the configured webhook.
func (d *Dispatcher) Send(event Event) error {
	if d.webhookURL == "" {
		return nil
	}

	var payload []byte
	var err error

	switch d.channel {
	case "slack":
		payload, err = d.buildSlack(event)
	case "discord":
		payload, err = d.buildDiscord(event)
	default:
		payload, err = d.buildGeneric(event)
	}
	if err != nil {
		return fmt.Errorf("build payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("post webhook: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	log.Debug().Str("channel", d.channel).Str("type", event.Type).Msg("alert dispatched")
	return nil
}

func (d *Dispatcher) buildSlack(event Event) ([]byte, error) {
	ts := event.Timestamp.Unix()
	msg := slackPayload{
		Text: "MCP Guard Alert",
		Attachments: []slackAttachment{
			{
				Color: "#ff0000",
				Fields: []slackField{
					{Title: "Type", Value: event.Type, Short: true},
					{Title: "Tool", Value: event.Tool, Short: true},
					{Title: "Identity", Value: event.Identity, Short: true},
					{Title: "Policy", Value: event.Policy, Short: true},
					{Title: "Reason", Value: event.Reason, Short: false},
				},
				Footer: "MCP Guard",
				Ts:     ts,
			},
		},
	}
	return json.Marshal(msg)
}

func (d *Dispatcher) buildDiscord(event Event) ([]byte, error) {
	fields := []discordField{
		{Name: "Type", Value: event.Type, Inline: true},
		{Name: "Tool", Value: event.Tool, Inline: true},
		{Name: "Identity", Value: event.Identity, Inline: true},
		{Name: "Policy", Value: event.Policy, Inline: true},
		{Name: "Reason", Value: event.Reason, Inline: false},
	}
	msg := discordPayload{
		Embeds: []discordEmbed{
			{
				Title:     "MCP Guard Alert",
				Color:     16711680,
				Fields:    fields,
				Timestamp: event.Timestamp.Format(time.RFC3339),
			},
		},
	}
	return json.Marshal(msg)
}

func (d *Dispatcher) buildGeneric(event Event) ([]byte, error) {
	return json.Marshal(event)
}

type slackPayload struct {
	Text        string            `json:"text"`
	Attachments []slackAttachment `json:"attachments"`
}

type slackAttachment struct {
	Color   string       `json:"color"`
	Fields  []slackField `json:"fields"`
	Footer  string       `json:"footer"`
	Ts      int64        `json:"ts"`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type discordPayload struct {
	Embeds []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title     string         `json:"title"`
	Color     int            `json:"color"`
	Fields    []discordField `json:"fields"`
	Timestamp string         `json:"timestamp"`
}

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
