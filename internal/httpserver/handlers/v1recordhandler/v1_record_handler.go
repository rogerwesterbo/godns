package v1recordhandler

import (
	"net/http"
	"strings"

	"github.com/rogerwesterbo/godns/internal/httpserver/helpers"
	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/internal/services/v1recordservice"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// RecordHandler handles DNS record endpoints
type RecordHandler struct {
	recordService *v1recordservice.V1RecordService
}

// NewRecordHandler creates a new record handler
func NewRecordHandler(recordService *v1recordservice.V1RecordService) *RecordHandler {
	return &RecordHandler{
		recordService: recordService,
	}
}

// @Summary Create a DNS record
// @Description Add a new DNS record to an existing zone
// @Tags Records
// @Accept json
// @Produce json
// @Param zone path string true "Zone name (e.g., example.lan)"
// @Param record body models.DNSRecord true "Record to create"
// @Success 201 {object} models.DNSRecord "Record created"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Zone not found"
// @Failure 409 {object} map[string]string "Record already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/zones/{zone}/records [post]
func (h *RecordHandler) CreateRecord(w http.ResponseWriter, req *http.Request, domain string) {
	var record models.DNSRecord
	if err := helpers.DecodeJSON(req.Body, &record); err != nil {
		helpers.SendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.recordService.CreateRecord(req.Context(), domain, &record); err != nil {
		vlog.Errorf("Failed to create record in zone %s: %v", domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Zone not found")
		} else if strings.Contains(err.Error(), "already exists") {
			helpers.SendError(w, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "invalid") {
			helpers.SendError(w, http.StatusBadRequest, err.Error())
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to create record")
		}
		return
	}

	helpers.SendJSON(w, http.StatusCreated, record)
}

// @Summary Get a DNS record
// @Description Get a specific DNS record by name and type
// @Tags Records
// @Produce json
// @Param zone path string true "Zone name (e.g., example.lan)"
// @Param name path string true "Record name (e.g., www.example.lan.)"
// @Param type path string true "Record type (e.g., A, AAAA, CNAME)"
// @Success 200 {object} models.DNSRecord "Record details"
// @Failure 404 {object} map[string]string "Zone or record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/zones/{zone}/records/{name}/{type} [get]
func (h *RecordHandler) GetRecord(w http.ResponseWriter, req *http.Request, domain, name, recordType string) {
	record, err := h.recordService.GetRecord(req.Context(), domain, name, recordType)
	if err != nil {
		vlog.Errorf("Failed to get record %s/%s in zone %s: %v", name, recordType, domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Record not found")
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to get record")
		}
		return
	}

	helpers.SendJSON(w, http.StatusOK, record)
}

// @Summary Update a DNS record
// @Description Update an existing DNS record
// @Tags Records
// @Accept json
// @Produce json
// @Param zone path string true "Zone name (e.g., example.lan)"
// @Param name path string true "Record name (e.g., www.example.lan.)"
// @Param type path string true "Record type (e.g., A, AAAA, CNAME)"
// @Param record body models.DNSRecord true "Updated record data"
// @Success 200 {object} models.DNSRecord "Record updated"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Zone or record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/zones/{zone}/records/{name}/{type} [put]
func (h *RecordHandler) UpdateRecord(w http.ResponseWriter, req *http.Request, domain, name, recordType string) {
	var record models.DNSRecord
	if err := helpers.DecodeJSON(req.Body, &record); err != nil {
		helpers.SendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.recordService.UpdateRecord(req.Context(), domain, name, recordType, &record); err != nil {
		vlog.Errorf("Failed to update record %s/%s in zone %s: %v", name, recordType, domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Record not found")
		} else if strings.Contains(err.Error(), "invalid") {
			helpers.SendError(w, http.StatusBadRequest, err.Error())
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to update record")
		}
		return
	}

	helpers.SendJSON(w, http.StatusOK, record)
}

// @Summary Delete a DNS record
// @Description Delete a specific DNS record from a zone
// @Tags Records
// @Param zone path string true "Zone name (e.g., example.lan)"
// @Param name path string true "Record name (e.g., www.example.lan.)"
// @Param type path string true "Record type (e.g., A, AAAA, CNAME)"
// @Success 204 "Record deleted"
// @Failure 404 {object} map[string]string "Zone or record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/zones/{zone}/records/{name}/{type} [delete]
func (h *RecordHandler) DeleteRecord(w http.ResponseWriter, req *http.Request, domain, name, recordType string) {
	if err := h.recordService.DeleteRecord(req.Context(), domain, name, recordType); err != nil {
		vlog.Errorf("Failed to delete record %s/%s in zone %s: %v", name, recordType, domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Record not found")
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to delete record")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Set DNS record status
// @Description Enable or disable a specific DNS record
// @Tags Records
// @Accept json
// @Param zone path string true "Zone name (e.g., example.lan)"
// @Param name path string true "Record name (e.g., www.example.lan.)"
// @Param type path string true "Record type (e.g., A, AAAA, CNAME)"
// @Param status body object{enabled=bool} true "Record status"
// @Success 204 "Record status updated"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Zone or record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/zones/{zone}/records/{name}/{type}/status [patch]
func (h *RecordHandler) SetRecordStatus(w http.ResponseWriter, req *http.Request, domain, name, recordType string) {
	var statusReq struct {
		Enabled bool `json:"enabled"`
	}
	
	if err := helpers.DecodeJSON(req.Body, &statusReq); err != nil {
		helpers.SendError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.recordService.SetRecordEnabled(req.Context(), domain, name, recordType, statusReq.Enabled); err != nil {
		vlog.Errorf("Failed to set record status for %s/%s in zone %s: %v", name, recordType, domain, err)
		if strings.Contains(err.Error(), "not found") {
			helpers.SendError(w, http.StatusNotFound, "Record not found")
		} else {
			helpers.SendError(w, http.StatusInternalServerError, "Failed to update record status")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
