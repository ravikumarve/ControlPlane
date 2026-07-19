package admin

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/audit"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/policy"
	"github.com/ravikumarve/ControlPlane/mcp-guard/internal/proxy"
)

// Server serves the admin HTTP API for the workspace UI.
type Server struct {
	mu       sync.RWMutex
	apiKey   string
	engine   *policy.Engine
	audit    *audit.Logger
	stats    *proxy.Stats
	configFn func() interface{}
	mux      *http.ServeMux
	srv      *http.Server
}

// New creates an admin API server.
func New(apiKey, listen string, engine *policy.Engine, auditLog *audit.Logger, stats *proxy.Stats, configFn func() interface{}) *Server {
	s := &Server{
		apiKey:   apiKey,
		engine:   engine,
		audit:    auditLog,
		stats:    stats,
		configFn: configFn,
		mux:      http.NewServeMux(),
	}

	s.mux.HandleFunc("POST /api/login", s.handleLogin)
	s.mux.HandleFunc("GET /api/status", s.authWrap(s.handleStatus))
	s.mux.HandleFunc("GET /api/policies", s.authWrap(s.handleListPolicies))
	s.mux.HandleFunc("POST /api/policies", s.authWrap(s.handleSavePolicies))
	s.mux.HandleFunc("GET /api/audit", s.authWrap(s.handleAudit))
	s.mux.HandleFunc("GET /api/config", s.authWrap(s.handleConfig))

	s.srv = &http.Server{
		Addr:    listen,
		Handler: s.corsWrap(s.mux),
	}

	return s
}

// Start begins the admin HTTP server.
func (s *Server) Start() error {
	log.Info().Str("addr", s.srv.Addr).Msg("admin API server starting")
	return s.srv.ListenAndServe()
}

// Stop shuts down the admin server.
func (s *Server) Stop() error {
	return s.srv.Close()
}

// corsWrap adds CORS headers for the workspace UI.
func (s *Server) corsWrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// authWrap checks the Authorization header against the configured API key.
func (s *Server) authWrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if key == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Strip "Bearer " prefix
		if len(key) > 7 && key[:7] == "Bearer " {
			key = key[7:]
		}

		if subtle.ConstantTimeCompare([]byte(key), []byte(s.apiKey)) != 1 {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// json writes a JSON response.
func jsonResp(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// --- Handlers ---

// handleLogin verifies the API key and returns a token.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		APIKey string `json:"api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	if subtle.ConstantTimeCompare([]byte(body.APIKey), []byte(s.apiKey)) != 1 {
		jsonResp(w, http.StatusUnauthorized, map[string]string{"error": "invalid api key"})
		return
	}

	jsonResp(w, http.StatusOK, map[string]string{
		"token": body.APIKey,
		"status": "ok",
	})
}

// handleStatus returns proxy statistics.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	snapshot := s.stats.Snapshot()
	jsonResp(w, http.StatusOK, map[string]interface{}{
		"status": "running",
		"stats":  snapshot,
	})
}

// handleListPolicies returns all configured policies.
func (s *Server) handleListPolicies(w http.ResponseWriter, r *http.Request) {
	policies := s.engine.List()
	jsonResp(w, http.StatusOK, map[string]interface{}{
		"policies": policies,
	})
}

// handleSavePolicies replaces all policies.
func (s *Server) handleSavePolicies(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Policies []policy.RawPolicy `json:"policies"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "invalid request: " + err.Error()})
		return
	}

	if err := s.engine.Replace(body.Policies); err != nil {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	jsonResp(w, http.StatusOK, map[string]string{"status": "saved"})
}

// handleAudit returns recent audit log entries.
func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	limit := 50
	entries, err := s.audit.Recent(limit)
	if err != nil {
		jsonResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsonResp(w, http.StatusOK, map[string]interface{}{
		"entries": entries,
	})
}

// handleConfig returns the current config.
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	cfg := s.configFn()
	jsonResp(w, http.StatusOK, map[string]interface{}{
		"config": cfg,
	})
}
