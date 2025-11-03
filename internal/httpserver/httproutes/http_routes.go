package httproutes

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	_ "github.com/rogerwesterbo/godns/docs" // swagger docs
	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// Router holds the zone service and provides HTTP routing
type Router struct {
	mux         *http.ServeMux
	zoneService *v1zoneservice.V1ZoneService
}

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(zoneService *v1zoneservice.V1ZoneService) *http.ServeMux {
	r := &Router{
		mux:         http.NewServeMux(),
		zoneService: zoneService,
	}

	r.registerRoutes()
	return r.mux
}

// registerRoutes sets up all HTTP routes
func (r *Router) registerRoutes() {
	// Swagger documentation
	r.mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Health check
	r.mux.HandleFunc("/health", r.handleHealth)
	r.mux.HandleFunc("/ready", r.handleReady)

	// API v1 routes
	r.mux.HandleFunc("/api/v1/zones", r.handleZones)
	r.mux.HandleFunc("/api/v1/zones/", r.handleZoneOperations)
}

// @Summary Health check
// @Description Check if the API server is healthy
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "status: healthy"
// @Router /health [get]
func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// @Summary Readiness check
// @Description Check if the API server is ready to accept requests
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "status: ready"
// @Router /ready [get]
func (r *Router) handleReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// Handle zones list and create
func (r *Router) handleZones(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.listZones(w, req)
	case http.MethodPost:
		r.createZone(w, req)
	default:
		r.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Handle individual zone operations and records
func (r *Router) handleZoneOperations(w http.ResponseWriter, req *http.Request) {
	// Parse path: /api/v1/zones/{domain}[/records[/{name}/{type}]]
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/zones/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		r.sendError(w, http.StatusBadRequest, "Domain is required")
		return
	}

	domain := parts[0]

	// Check if this is a record operation
	if len(parts) >= 2 && parts[1] == "records" {
		r.handleRecordOperations(w, req, domain, parts[2:])
		return
	}

	// Zone operations
	switch req.Method {
	case http.MethodGet:
		r.getZone(w, req, domain)
	case http.MethodPut:
		r.updateZone(w, req, domain)
	case http.MethodDelete:
		r.deleteZone(w, req, domain)
	default:
		r.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Handle record operations
func (r *Router) handleRecordOperations(w http.ResponseWriter, req *http.Request, domain string, parts []string) {
	// GET /api/v1/zones/{domain}/records - List all records (via zone)
	if len(parts) == 0 {
		switch req.Method {
		case http.MethodPost:
			r.createRecord(w, req, domain)
		default:
			r.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// Operations on specific record: /api/v1/zones/{domain}/records/{name}/{type}
	if len(parts) < 2 {
		r.sendError(w, http.StatusBadRequest, "Record name and type are required")
		return
	}

	name := parts[0]
	recordType := strings.ToUpper(parts[1])

	switch req.Method {
	case http.MethodGet:
		r.getRecord(w, req, domain, name, recordType)
	case http.MethodPut:
		r.updateRecord(w, req, domain, name, recordType)
	case http.MethodDelete:
		r.deleteRecord(w, req, domain, name, recordType)
	default:
		r.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Zone handlers

// @Summary List all DNS zones
// @Description Get a list of all DNS zones
// @Tags Zones
// @Produce json
// @Success 200 {array} models.DNSZone "List of zones"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones [get]
func (r *Router) listZones(w http.ResponseWriter, req *http.Request) {
	zones, err := r.zoneService.ListZones(req.Context())
	if err != nil {
		vlog.Errorf("Failed to list zones: %v", err)
		r.sendError(w, http.StatusInternalServerError, "Failed to list zones")
		return
	}

	r.sendJSON(w, http.StatusOK, zones)
}

// @Summary Create a new DNS zone
// @Description Create a new DNS zone with optional records
// @Tags Zones
// @Accept json
// @Produce json
// @Param zone body models.DNSZone true "Zone to create"
// @Success 201 {object} models.DNSZone "Zone created"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 409 {object} map[string]string "Zone already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones [post]
func (r *Router) createZone(w http.ResponseWriter, req *http.Request) {
	var zone models.DNSZone
	if err := r.decodeJSON(req.Body, &zone); err != nil {
		r.sendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.zoneService.CreateZone(req.Context(), &zone); err != nil {
		vlog.Errorf("Failed to create zone: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			r.sendError(w, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "invalid") {
			r.sendError(w, http.StatusBadRequest, err.Error())
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to create zone")
		}
		return
	}

	r.sendJSON(w, http.StatusCreated, zone)
}

// @Summary Get a DNS zone
// @Description Get a specific DNS zone by domain
// @Tags Zones
// @Produce json
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Success 200 {object} models.DNSZone "Zone details"
// @Failure 404 {object} map[string]string "Zone not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain} [get]
func (r *Router) getZone(w http.ResponseWriter, req *http.Request, domain string) {
	zone, err := r.zoneService.GetZone(req.Context(), domain)
	if err != nil {
		vlog.Errorf("Failed to get zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			r.sendError(w, http.StatusNotFound, "Zone not found")
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to get zone")
		}
		return
	}

	r.sendJSON(w, http.StatusOK, zone)
}

// @Summary Update a DNS zone
// @Description Update an existing DNS zone (replaces all records)
// @Tags Zones
// @Accept json
// @Produce json
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Param zone body models.DNSZone true "Updated zone data"
// @Success 200 {object} models.DNSZone "Zone updated"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Zone not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain} [put]
func (r *Router) updateZone(w http.ResponseWriter, req *http.Request, domain string) {
	var zone models.DNSZone
	if err := r.decodeJSON(req.Body, &zone); err != nil {
		r.sendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.zoneService.UpdateZone(req.Context(), domain, &zone); err != nil {
		vlog.Errorf("Failed to update zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			r.sendError(w, http.StatusNotFound, "Zone not found")
		} else if strings.Contains(err.Error(), "invalid") {
			r.sendError(w, http.StatusBadRequest, err.Error())
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to update zone")
		}
		return
	}

	r.sendJSON(w, http.StatusOK, zone)
}

// @Summary Delete a DNS zone
// @Description Delete a DNS zone and all its records
// @Tags Zones
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Success 204 "Zone deleted"
// @Failure 404 {object} map[string]string "Zone not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain} [delete]
func (r *Router) deleteZone(w http.ResponseWriter, req *http.Request, domain string) {
	if err := r.zoneService.DeleteZone(req.Context(), domain); err != nil {
		vlog.Errorf("Failed to delete zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			r.sendError(w, http.StatusNotFound, "Zone not found")
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to delete zone")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Record handlers

// @Summary Create a DNS record
// @Description Add a new DNS record to an existing zone
// @Tags Records
// @Accept json
// @Produce json
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Param record body models.DNSRecord true "Record to create"
// @Success 201 {object} models.DNSRecord "Record created"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Zone not found"
// @Failure 409 {object} map[string]string "Record already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain}/records [post]
func (r *Router) createRecord(w http.ResponseWriter, req *http.Request, domain string) {
	var record models.DNSRecord
	if err := r.decodeJSON(req.Body, &record); err != nil {
		r.sendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.zoneService.CreateRecord(req.Context(), domain, &record); err != nil {
		vlog.Errorf("Failed to create record in zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			r.sendError(w, http.StatusNotFound, "Zone not found")
		} else if strings.Contains(err.Error(), "already exists") {
			r.sendError(w, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "invalid") {
			r.sendError(w, http.StatusBadRequest, err.Error())
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to create record")
		}
		return
	}

	r.sendJSON(w, http.StatusCreated, record)
}

// @Summary Get a DNS record
// @Description Get a specific DNS record by name and type
// @Tags Records
// @Produce json
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Param name path string true "Record name (e.g., www.example.lan.)"
// @Param type path string true "Record type (e.g., A, AAAA, CNAME)"
// @Success 200 {object} models.DNSRecord "Record details"
// @Failure 404 {object} map[string]string "Zone or record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain}/records/{name}/{type} [get]
func (r *Router) getRecord(w http.ResponseWriter, req *http.Request, domain, name, recordType string) {
	record, err := r.zoneService.GetRecord(req.Context(), domain, name, recordType)
	if err != nil {
		vlog.Errorf("Failed to get record %s/%s in zone %s: %v", name, recordType, domain, err)
		if strings.Contains(err.Error(), "not found") {
			r.sendError(w, http.StatusNotFound, "Record not found")
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to get record")
		}
		return
	}

	r.sendJSON(w, http.StatusOK, record)
}

// @Summary Update a DNS record
// @Description Update an existing DNS record
// @Tags Records
// @Accept json
// @Produce json
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Param name path string true "Record name (e.g., www.example.lan.)"
// @Param type path string true "Record type (e.g., A, AAAA, CNAME)"
// @Param record body models.DNSRecord true "Updated record data"
// @Success 200 {object} models.DNSRecord "Record updated"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Zone or record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain}/records/{name}/{type} [put]
func (r *Router) updateRecord(w http.ResponseWriter, req *http.Request, domain, name, recordType string) {
	var record models.DNSRecord
	if err := r.decodeJSON(req.Body, &record); err != nil {
		r.sendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := r.zoneService.UpdateRecord(req.Context(), domain, name, recordType, &record); err != nil {
		vlog.Errorf("Failed to update record %s/%s in zone %s: %v", name, recordType, domain, err)
		if strings.Contains(err.Error(), "not found") {
			r.sendError(w, http.StatusNotFound, "Record not found")
		} else if strings.Contains(err.Error(), "invalid") {
			r.sendError(w, http.StatusBadRequest, err.Error())
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to update record")
		}
		return
	}

	r.sendJSON(w, http.StatusOK, record)
}

// @Summary Delete a DNS record
// @Description Delete a specific DNS record from a zone
// @Tags Records
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Param name path string true "Record name (e.g., www.example.lan.)"
// @Param type path string true "Record type (e.g., A, AAAA, CNAME)"
// @Success 204 "Record deleted"
// @Failure 404 {object} map[string]string "Zone or record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain}/records/{name}/{type} [delete]
func (r *Router) deleteRecord(w http.ResponseWriter, req *http.Request, domain, name, recordType string) {
	if err := r.zoneService.DeleteRecord(req.Context(), domain, name, recordType); err != nil {
		vlog.Errorf("Failed to delete record %s/%s in zone %s: %v", name, recordType, domain, err)
		if strings.Contains(err.Error(), "not found") {
			r.sendError(w, http.StatusNotFound, "Record not found")
		} else {
			r.sendError(w, http.StatusInternalServerError, "Failed to delete record")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func (r *Router) decodeJSON(body io.Reader, v interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func (r *Router) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		vlog.Errorf("Failed to encode JSON response: %v", err)
	}
}

func (r *Router) sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
