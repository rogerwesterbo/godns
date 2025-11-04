package v1searchservice

import (
	"context"
	"slices"
	"strings"

	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/rogerwesterbo/godns/internal/services/v1zoneservice"
)

// SearchResultType represents the type of search result
type SearchResultType string

const (
	SearchResultTypeZone   SearchResultType = "zone"
	SearchResultTypeRecord SearchResultType = "record"
)

// SearchResult represents a single search result
type SearchResult struct {
	Type   SearchResultType  `json:"type" example:"zone"`                  // Type of result (zone, record)
	Zone   string            `json:"zone,omitempty" example:"example.lan"` // Zone name
	Record *models.DNSRecord `json:"record,omitempty"`                     // Record details (if type is record)
}

// SearchResponse represents the search API response
type SearchResponse struct {
	Query   string         `json:"query" example:"example"` // The search query
	Results []SearchResult `json:"results"`                 // List of search results
	Count   int            `json:"count" example:"5"`       // Number of results found
}

// V1SearchService provides search functionality across DNS zones and records
type V1SearchService struct {
	zoneService *v1zoneservice.V1ZoneService
}

// NewV1SearchService creates a new search service
func NewV1SearchService(zoneService *v1zoneservice.V1ZoneService) *V1SearchService {
	return &V1SearchService{
		zoneService: zoneService,
	}
}

// Search performs a search across zones and records
// The query is case-insensitive and searches in:
// - Zone names (domain)
// - Record names
// - Record values
// - Record types
func (s *V1SearchService) Search(ctx context.Context, query string, types []SearchResultType) (*SearchResponse, error) {
	if query == "" {
		return &SearchResponse{
			Query:   query,
			Results: []SearchResult{},
			Count:   0,
		}, nil
	}

	// Normalize query to lowercase for case-insensitive search
	normalizedQuery := strings.ToLower(query)

	// Determine which types to search
	searchZones := len(types) == 0 || contains(types, SearchResultTypeZone)
	searchRecords := len(types) == 0 || contains(types, SearchResultTypeRecord)

	var results []SearchResult

	// Get all zones
	zones, err := s.zoneService.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	// Search through zones and records
	for _, zone := range zones {
		// Search in zone names
		if searchZones && strings.Contains(strings.ToLower(zone.Domain), normalizedQuery) {
			results = append(results, SearchResult{
				Type: SearchResultTypeZone,
				Zone: zone.Domain,
			})
		}

		// Search in records
		if searchRecords {
			for _, record := range zone.Records {
				// Check if query matches record name, type, or value
				if strings.Contains(strings.ToLower(record.Name), normalizedQuery) ||
					strings.Contains(strings.ToLower(record.Type), normalizedQuery) ||
					strings.Contains(strings.ToLower(record.Value), normalizedQuery) {

					results = append(results, SearchResult{
						Type:   SearchResultTypeRecord,
						Zone:   zone.Domain,
						Record: &record,
					})
				}
			}
		}
	}

	return &SearchResponse{
		Query:   query,
		Results: results,
		Count:   len(results),
	}, nil
}

// contains checks if a slice contains a specific SearchResultType
func contains(slice []SearchResultType, item SearchResultType) bool {
	return slices.Contains(slice, item)
}
