package v1zonehandler

import (
	"net/http"
	"strings"

	"github.com/rogerwesterbo/godns/internal/httpserver/helpers"
	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// ZoneHandler handles DNS zone endpoints
type ZoneHandler struct {
	zoneService *v1zoneservice.V1ZoneService
}

// NewZoneHandler creates a new zone handler
func NewZoneHandler(zoneService *v1zoneservice.V1ZoneService) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
	}
}

// @Summary List all DNS zones
// @Description Get a list of all DNS zones
// @Tags Zones
// @Produce json
// @Success 200 {array} models.DNSZone "List of zones"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones [get]
func (h *ZoneHandler) ListZones(w http.ResponseWriter, req *http.Request) {
	zones, err := h.zoneService.ListZones(req.Context())
	if err != nil {
		vlog.Errorf("Failed to list zones: %v", err)
		helpers.SendError(w, http.StatusInternalServerError, "Failed to list zones")
		return
	}

	helpers.SendJSON(w, http.StatusOK, zones)
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
func (h *ZoneHandler) CreateZone(w http.ResponseWriter, req *http.Request) {
	var zone models.DNSZone
	if err := helpers.DecodeJSON(req.Body, &zone); err != nil {
		helpers.SendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.zoneService.CreateZone(req.Context(), &zone); err != nil {
		vlog.Errorf("Failed to create zone: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			helpers.SendError(w, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "invalid") {
			helpers.SendError(w, http.StatusBadRequest, err.Error())
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to create zone")
		}
		return
	}

	helpers.SendJSON(w, http.StatusCreated, zone)
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
func (h *ZoneHandler) GetZone(w http.ResponseWriter, req *http.Request, domain string) {
	zone, err := h.zoneService.GetZone(req.Context(), domain)
	if err != nil {
		vlog.Errorf("Failed to get zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Zone not found")
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to get zone")
		}
		return
	}

	helpers.SendJSON(w, http.StatusOK, zone)
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
func (h *ZoneHandler) UpdateZone(w http.ResponseWriter, req *http.Request, domain string) {
	var zone models.DNSZone
	if err := helpers.DecodeJSON(req.Body, &zone); err != nil {
		helpers.SendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.zoneService.UpdateZone(req.Context(), domain, &zone); err != nil {
		vlog.Errorf("Failed to update zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Zone not found")
		} else if strings.Contains(err.Error(), "invalid") {
			helpers.SendError(w, http.StatusBadRequest, err.Error())
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to update zone")
		}
		return
	}

	helpers.SendJSON(w, http.StatusOK, zone)
}

// @Summary Delete a DNS zone
// @Description Delete a DNS zone and all its records
// @Tags Zones
// @Param domain path string true "Domain name (e.g., example.lan)"
// @Success 204 "Zone deleted"
// @Failure 404 {object} map[string]string "Zone not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/zones/{domain} [delete]
func (h *ZoneHandler) DeleteZone(w http.ResponseWriter, req *http.Request, domain string) {
	if err := h.zoneService.DeleteZone(req.Context(), domain); err != nil {
		vlog.Errorf("Failed to delete zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Zone not found")
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to delete zone")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
