package v1healthhandler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// @Summary Health check
// @Description Check if the API server is healthy
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "status: healthy"
// @Router /health [get]
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// @Summary Readiness check
// @Description Check if the API server is ready to accept requests
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "status: ready"
// @Router /ready [get]
func (h *HealthHandler) HandleReady(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
