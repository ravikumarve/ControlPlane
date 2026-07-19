package hitl

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/config"
)

// Request represents a human-in-the-loop approval request.
type Request struct {
	ID        string      `json:"id"`
	Identity  string      `json:"identity"`
	Tool      string      `json:"tool"`
	Params    any         `json:"params"`
	RawData   string      `json:"raw_data,omitempty"`
	RiskScore float64     `json:"risk_score"`
	Status    string      `json:"status"` // pending | approved | denied | expired
	CreatedAt time.Time   `json:"created_at"`
	ExpiresAt time.Time   `json:"expires_at"`
	ApprovedBy string     `json:"approved_by,omitempty"`
}

// Handler manages the lifecycle of HITL approval requests.
type Handler struct {
	cfg        *config.HITLConfig
	mu         sync.RWMutex
	pending    map[string]*Request
	webhookURL string
}

// NewHandler creates an HITL handler.
func NewHandler(cfg *config.HITLConfig) *Handler {
	return &Handler{
		cfg:        cfg,
		pending:    make(map[string]*Request),
		webhookURL: cfg.WebhookURL,
	}
}

// Submit creates a new approval request and dispatches notifications.
func (h *Handler) Submit(req Request) {
	req.ID = fmt.Sprintf("req-%s-%d", req.Tool, rand.Intn(99999))
	req.Status = "pending"
	req.CreatedAt = time.Now()

	// Default timeout
	timeout := 5 * time.Minute
	if h.cfg.Timeout != "" {
		if d, err := time.ParseDuration(h.cfg.Timeout); err == nil {
			timeout = d
		}
	}
	req.ExpiresAt = req.CreatedAt.Add(timeout)

	h.mu.Lock()
	h.pending[req.ID] = &req
	h.mu.Unlock()

	// Auto-expire
	go func() {
		<-time.After(timeout)
		h.mu.Lock()
		if r, ok := h.pending[req.ID]; ok && r.Status == "pending" {
			r.Status = "expired"
			log.Warn().Str("id", req.ID).Msg("HITL request expired")
		}
		h.mu.Unlock()
	}()

	// Dispatch webhook
	if h.webhookURL != "" {
		h.dispatchWebhook(req)
	}

	log.Info().
		Str("id", req.ID).
		Str("tool", req.Tool).
		Str("identity", req.Identity).
		Float64("risk", req.RiskScore).
		Msg("HITL request submitted")
}

// Approve marks a request as approved.
func (h *Handler) Approve(id string, by string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	req, ok := h.pending[id]
	if !ok {
		return fmt.Errorf("request %s not found", id)
	}
	if req.Status != "pending" {
		return fmt.Errorf("request %s is not pending (status: %s)", id, req.Status)
	}

	req.Status = "approved"
	req.ApprovedBy = by

	log.Info().Str("id", id).Str("by", by).Msg("HITL request approved")
	return nil
}

// Deny marks a request as denied.
func (h *Handler) Deny(id string, by string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	req, ok := h.pending[id]
	if !ok {
		return fmt.Errorf("request %s not found", id)
	}
	if req.Status != "pending" {
		return fmt.Errorf("request %s is not pending", id)
	}

	req.Status = "denied"
	req.ApprovedBy = by

	log.Info().Str("id", id).Str("by", by).Msg("HITL request denied")
	return nil
}

// ListPending returns all pending approval requests.
func (h *Handler) ListPending() []*Request {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []*Request
	for _, req := range h.pending {
		if req.Status == "pending" {
			result = append(result, req)
		}
	}
	return result
}

// dispatchWebhook sends the approval request to the configured webhook URL.
func (h *Handler) dispatchWebhook(req Request) {
	payload := map[string]any{
		"type":    "mcp_guard_approval",
		"id":      req.ID,
		"tool":    req.Tool,
		"identity": req.Identity,
		"risk_score": req.RiskScore,
		"expires_at": req.ExpiresAt.Format(time.RFC3339),
		"approve_url": fmt.Sprintf("https://localhost:8443/hitl/%s/approve", req.ID),
		"deny_url":    fmt.Sprintf("https://localhost:8443/hitl/%s/deny", req.ID),
	}

	data, _ := json.Marshal(payload)
	log.Debug().
		Str("url", h.webhookURL).
		RawJSON("payload", data).
		Msg("dispatching HITL webhook")

	// In production, use net/http to POST the webhook.
	// For MVP, we log the approval URL.
	log.Info().
		Str("approve", fmt.Sprintf("mcp-guard approve %s", req.ID)).
		Str("deny", fmt.Sprintf("mcp-guard approve %s --deny", req.ID)).
		Msg("HITL approval URLs")
}
