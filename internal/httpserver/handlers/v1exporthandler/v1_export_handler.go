package v1exporthandler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rogerwesterbo/godns/internal/httpserver/helpers"
	"github.com/rogerwesterbo/godns/internal/services/v1exportservice"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// ExportHandler handles DNS zone export endpoints
type ExportHandler struct {
	exportService *v1exportservice.V1ExportService
}

// NewExportHandler creates a new export handler
func NewExportHandler(exportService *v1exportservice.V1ExportService) *ExportHandler {
	return &ExportHandler{
		exportService: exportService,
	}
}

// @Summary Export all zones
// @Description Export all DNS zones in a specified format (coredns, powerdns, bind, zonefile)
// @Tags Export
// @Produce plain
// @Param format query string false "Export format: coredns, powerdns, bind, or zonefile" default(bind)
// @Success 200 {string} string "Exported zone configuration"
// @Failure 400 {object} map[string]string "Invalid format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/export [get]
func (h *ExportHandler) ExportAll(w http.ResponseWriter, req *http.Request) {
	// Get format from query parameter
	format := req.URL.Query().Get("format")
	if format == "" {
		format = "bind" // Default to BIND format
	}

	// Validate format
	if !v1exportservice.ValidateFormat(format) {
		helpers.SendError(w, http.StatusBadRequest, "Invalid format. Supported formats: coredns, powerdns, bind, zonefile")
		return
	}

	// Export all zones
	exported, err := h.exportService.ExportAllZones(req.Context(), v1exportservice.ExportFormat(format))
	if err != nil {
		vlog.Errorf("Failed to export zones: %v", err)
		helpers.SendError(w, http.StatusInternalServerError, "Failed to export zones")
		return
	}

	// Set content disposition for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"zones-%s.txt\"", format))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(exported))
}

// @Summary Export a specific zone
// @Description Export a specific DNS zone in a specified format
// @Tags Export
// @Produce plain
// @Param zone path string true "Zone name (e.g., example.lan)"
// @Param format query string false "Export format: coredns, powerdns, bind, or zonefile" default(bind)
// @Success 200 {string} string "Exported zone configuration"
// @Failure 400 {object} map[string]string "Invalid format"
// @Failure 404 {object} map[string]string "Zone not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/export/{zone} [get]
func (h *ExportHandler) ExportZone(w http.ResponseWriter, req *http.Request, domain string) {
	if domain == "" {
		helpers.SendError(w, http.StatusBadRequest, "Domain is required")
		return
	}

	// Get format from query parameter
	format := req.URL.Query().Get("format")
	if format == "" {
		format = "bind" // Default to BIND format
	}

	// Validate format
	if !v1exportservice.ValidateFormat(format) {
		helpers.SendError(w, http.StatusBadRequest, "Invalid format. Supported formats: coredns, powerdns, bind, zonefile")
		return
	}

	// Export the zone
	exported, err := h.exportService.ExportZone(req.Context(), domain, v1exportservice.ExportFormat(format))
	if err != nil {
		vlog.Errorf("Failed to export zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Zone not found")
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to export zone")
		}
		return
	}

	// Set content disposition for file download
	zoneName := strings.TrimSuffix(domain, ".")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s-%s.txt\"", zoneName, format))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(exported))
}
