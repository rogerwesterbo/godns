package httproutes

import (
	"net/http"
	"strings"

	_ "github.com/rogerwesterbo/godns/docs" // swagger docs
	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1exporthandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1recordhandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1zonehandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/middleware"
	"github.com/rogerwesterbo/godns/internal/services/v1exportservice"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router holds the handlers and provides HTTP routing
type Router struct {
	mux           *http.ServeMux
	zoneHandler   *v1zonehandler.ZoneHandler
	recordHandler *v1recordhandler.RecordHandler
	exportHandler *v1exporthandler.ExportHandler
}

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(zoneService *v1zoneservice.V1ZoneService) *http.ServeMux {
	exportService := v1exportservice.NewV1ExportService(zoneService)

	r := &Router{
		mux:           http.NewServeMux(),
		zoneHandler:   v1zonehandler.NewZoneHandler(zoneService),
		recordHandler: v1recordhandler.NewRecordHandler(zoneService),
		exportHandler: v1exporthandler.NewExportHandler(exportService),
	}

	r.registerRoutes()
	return r.mux
}

// registerRoutes sets up all HTTP routes
func (r *Router) registerRoutes() {
	// Swagger documentation
	r.mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// API v1 routes - wrap all API routes with a base handler that applies JSON middleware by default
	r.mux.HandleFunc("/api/", r.apiRouter)
}

// apiRouter routes all /api/* requests and applies appropriate middleware
func (r *Router) apiRouter(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Export endpoints use plain text middleware (exception to JSON default)
	if strings.HasPrefix(path, "/api/v1/export") {
		middleware.PlainTextContentType(r.handleAPIRoutes)(w, req)
		return
	}

	// All other API routes use JSON middleware
	middleware.JSONContentType(r.handleAPIRoutes)(w, req)
}

// handleAPIRoutes handles the actual routing logic for API endpoints
func (r *Router) handleAPIRoutes(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Route to appropriate handler
	switch {
	case path == "/api/v1/zones":
		r.handleZones(w, req)
	case strings.HasPrefix(path, "/api/v1/zones/"):
		r.handleZoneOperations(w, req)
	case path == "/api/v1/export":
		r.handleExport(w, req)
	case strings.HasPrefix(path, "/api/v1/export/"):
		r.handleExportZone(w, req)
	default:
		http.NotFound(w, req)
	}
}

// Handle export all zones
func (r *Router) handleExport(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.exportHandler.ExportAll(w, req)
}

// Handle export specific zone
func (r *Router) handleExportZone(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract domain from path
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/export/")
	domain := path

	r.exportHandler.ExportZone(w, req, domain)
}

// Handle zones list and create
func (r *Router) handleZones(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.zoneHandler.ListZones(w, req)
	case http.MethodPost:
		r.zoneHandler.CreateZone(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle individual zone operations and records
func (r *Router) handleZoneOperations(w http.ResponseWriter, req *http.Request) {
	// Parse path: /api/v1/zones/{domain}[/records[/{name}/{type}]]
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/zones/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
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
		r.zoneHandler.GetZone(w, req, domain)
	case http.MethodPut:
		r.zoneHandler.UpdateZone(w, req, domain)
	case http.MethodDelete:
		r.zoneHandler.DeleteZone(w, req, domain)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle record operations
func (r *Router) handleRecordOperations(w http.ResponseWriter, req *http.Request, domain string, parts []string) {
	// POST /api/v1/zones/{domain}/records - Create a record
	if len(parts) == 0 {
		switch req.Method {
		case http.MethodPost:
			r.recordHandler.CreateRecord(w, req, domain)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Operations on specific record: /api/v1/zones/{domain}/records/{name}/{type}
	if len(parts) < 2 {
		http.Error(w, "Record name and type are required", http.StatusBadRequest)
		return
	}

	name := parts[0]
	recordType := strings.ToUpper(parts[1])

	switch req.Method {
	case http.MethodGet:
		r.recordHandler.GetRecord(w, req, domain, name, recordType)
	case http.MethodPut:
		r.recordHandler.UpdateRecord(w, req, domain, name, recordType)
	case http.MethodDelete:
		r.recordHandler.DeleteRecord(w, req, domain, name, recordType)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
