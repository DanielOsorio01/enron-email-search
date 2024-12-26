package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanielOsorio01/enron-email-search/back/repository/email"
)

type Email struct {
	Repo *email.ZincsearchRepo
}

// Response represents the standard API response structure
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

func (h *Email) List(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// Get search parameters with defaults
	term := query.Get("term")
	if term == "" {
		sendError(w, "Search term is required", http.StatusBadRequest)
		return
	}

	from, _ := strconv.Atoi(query.Get("from"))
	maxResults, _ := strconv.Atoi(query.Get("max_results"))
	field := query.Get("field")
	searchType := query.Get("search_type")

	// Create search parameters
	params := email.DefaultSearchParams()
	params.Term = term

	// Override defaults if provided in query
	if from > 0 {
		params.From = from
	}
	if maxResults > 0 {
		params.MaxResults = maxResults
	}
	if field != "" {
		params.Field = field
	}
	if searchType != "" {
		params.SearchType = searchType
	}

	// Get sort fields if provided
	if sortFields := query["sort_fields"]; len(sortFields) > 0 {
		params.SortFields = sortFields
	}

	// Get source fields if provided
	if sourceFields := query["source_fields"]; len(sourceFields) > 0 {
		params.SourceFields = sourceFields
	}

	// Perform the search
	emails, err := h.Repo.Search(r.Context(), "enron_emails", params)
	if err != nil {
		sendError(w, "Failed to search emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get total count for the search term
	total, err := h.Repo.CountResults(r.Context(), "enron_emails", term)
	if err != nil {
		sendError(w, "Failed to get total count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := Response{
		Success: true,
		Data: map[string]interface{}{
			"emails": emails,
			"total":  total,
			"from":   params.From,
			"size":   len(emails),
		},
	}

	// Send response
	sendJSON(w, response, http.StatusOK)
}

// Helper function to send JSON response
func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function to send error response
func sendError(w http.ResponseWriter, message string, status int) {
	response := Response{
		Success: false,
		Error:   message,
	}
	sendJSON(w, response, status)
}
