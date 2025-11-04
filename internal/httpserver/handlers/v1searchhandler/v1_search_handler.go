package v1searchhandler

import (
	"net/http"
	"strings"

	"github.com/rogerwesterbo/godns/internal/httpserver/helpers"
	"github.com/rogerwesterbo/godns/internal/services/v1searchservice"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// SearchHandler handles search endpoints
type SearchHandler struct {
	searchService *v1searchservice.V1SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService *v1searchservice.V1SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// @Summary Search DNS zones and records
// @Description Search across DNS zones and records with optional type filtering
// @Tags Search
// @Produce json
// @Param q query string true "Search query (case-insensitive)"
// @Param type query []string false "Filter by result type (zone, record). Can specify multiple types." Enums(zone, record)
// @Success 200 {object} v1searchservice.SearchResponse "Search results"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/search [get]
func (h *SearchHandler) Search(w http.ResponseWriter, req *http.Request) {
	// Get query parameter
	query := req.URL.Query().Get("q")
	if query == "" {
		helpers.SendError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	// Get type filter (optional)
	typeParams := req.URL.Query()["type"]
	var types []v1searchservice.SearchResultType
	
	for _, t := range typeParams {
		switch strings.ToLower(t) {
		case "zone":
			types = append(types, v1searchservice.SearchResultTypeZone)
		case "record":
			types = append(types, v1searchservice.SearchResultTypeRecord)
		default:
			helpers.SendError(w, http.StatusBadRequest, "Invalid type parameter. Allowed values: zone, record")
			return
		}
	}

	// Perform search
	results, err := h.searchService.Search(req.Context(), query, types)
	if err != nil {
		vlog.Errorf("Failed to perform search for query '%s': %v", query, err)
		helpers.SendError(w, http.StatusInternalServerError, "Failed to perform search")
		return
	}

	helpers.SendJSON(w, http.StatusOK, results)
}
