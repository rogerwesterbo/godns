package httproutes

import (
	"net/http"
	"strings"

	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1adminhandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1exporthandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1recordhandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1searchhandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/handlers/v1zonehandler"
	"github.com/rogerwesterbo/godns/internal/httpserver/middleware"
	_ "github.com/rogerwesterbo/godns/internal/httpserver/swaggerdocs" // swagger docs
	"github.com/rogerwesterbo/godns/internal/services/v1cacheservice"
	"github.com/rogerwesterbo/godns/internal/services/v1exportservice"
	"github.com/rogerwesterbo/godns/internal/services/v1healthcheckservice"
	"github.com/rogerwesterbo/godns/internal/services/v1loadbalancerservice"
	"github.com/rogerwesterbo/godns/internal/services/v1querylogservice"
	"github.com/rogerwesterbo/godns/internal/services/v1ratelimitservice"
	"github.com/rogerwesterbo/godns/internal/services/v1recordservice"
	"github.com/rogerwesterbo/godns/internal/services/v1searchservice"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router holds the handlers and provides HTTP routing
type Router struct {
	mux            *http.ServeMux
	zoneHandler    *v1zonehandler.ZoneHandler
	recordHandler  *v1recordhandler.RecordHandler
	exportHandler  *v1exporthandler.ExportHandler
	searchHandler  *v1searchhandler.SearchHandler
	adminHandler   *v1adminhandler.AdminHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(
	zoneService *v1zoneservice.V1ZoneService,
	cacheService *v1cacheservice.DNSCache,
	rateLimiter *v1ratelimitservice.RateLimiter,
	loadBalancer *v1loadbalancerservice.LoadBalancer,
	healthCheck *v1healthcheckservice.HealthCheckService,
	queryLog *v1querylogservice.QueryLogService,
	authMiddleware *middleware.AuthMiddleware,
) *http.ServeMux {
	exportService := v1exportservice.NewV1ExportService(zoneService)
	searchService := v1searchservice.NewV1SearchService(zoneService)

	r := &Router{
		mux:            http.NewServeMux(),
		zoneHandler:    v1zonehandler.NewZoneHandler(zoneService),
		recordHandler:  v1recordhandler.NewRecordHandler(v1recordservice.NewV1RecordService(zoneService.GetClient())),
		exportHandler:  v1exporthandler.NewExportHandler(exportService),
		searchHandler:  v1searchhandler.NewSearchHandler(searchService),
		adminHandler:   v1adminhandler.NewAdminHandler(cacheService, rateLimiter, loadBalancer, healthCheck, queryLog),
		authMiddleware: authMiddleware,
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

	// Wrap handler with authentication middleware
	authenticatedHandler := r.authMiddleware.Authenticate(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		// Export endpoints use plain text middleware (exception to JSON default)
		if strings.HasPrefix(path, "/api/v1/export") {
			middleware.PlainTextContentType(r.handleAPIRoutes)(rw, request)
			return
		}

		// All other API routes use JSON middleware
		middleware.JSONContentType(r.handleAPIRoutes)(rw, request)
	}))

	authenticatedHandler.ServeHTTP(w, req)
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
	case path == "/api/v1/search":
		r.handleSearch(w, req)
	case path == "/api/v1/export":
		r.handleExport(w, req)
	case strings.HasPrefix(path, "/api/v1/export/"):
		r.handleExportZone(w, req)
	case strings.HasPrefix(path, "/api/v1/admin/"):
		r.handleAdmin(w, req)
	default:
		http.NotFound(w, req)
	}
}

// Handle search
func (r *Router) handleSearch(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.searchHandler.Search(w, req)
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

	// Check if this is a status operation
	if len(parts) >= 2 && parts[1] == "status" {
		if req.Method == http.MethodPatch {
			r.zoneHandler.SetZoneStatus(w, req, domain)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

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

	// Check if this is a status operation
	if len(parts) >= 3 && parts[2] == "status" {
		if req.Method == http.MethodPatch {
			r.recordHandler.SetRecordStatus(w, req, domain, name, recordType)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

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

// Handle admin operations
func (r *Router) handleAdmin(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/admin/")

	switch path {
	case "stats":
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.adminHandler.GetSystemStats(w, req)

	case "cache/stats":
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.adminHandler.GetCacheStats(w, req)

	case "cache/clear":
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.adminHandler.ClearCache(w, req)

	case "loadbalancer/stats":
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.adminHandler.GetLoadBalancerStats(w, req)

	case "healthcheck/stats":
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.adminHandler.GetHealthCheckStats(w, req)

	case "querylog/stats":
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.adminHandler.GetQueryLogStats(w, req)

	case "ratelimiter/stats":
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.adminHandler.GetRateLimiterStats(w, req)

	default:
		http.NotFound(w, req)
	}
}
