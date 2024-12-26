package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ZincClient struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

// NewZincClient creates a new instance of ZincClient
func NewZincClient(baseURL, username, password string) *ZincClient {
	return &ZincClient{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (zc *ZincClient) newRequest(method, endpoint string, body []byte) (*http.Request, error) {
	url := zc.baseURL + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(zc.username, zc.password)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (zc *ZincClient) Ping() error {
	req, err := zc.newRequest("GET", "/", nil)
	if err != nil {
		return err
	}
	resp, err := zc.httpClient.Get(req.URL.String())
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// SearchRequest defines the structure for the search request body
type SearchRequest struct {
	SearchType string   `json:"search_type"`
	Query      Query    `json:"query"`
	SortFields []string `json:"sort_fields"`
	From       int      `json:"from"`
	MaxResults int      `json:"max_results"`
	Source     []string `json:"_source"`
}

// Query defines the query section of the search request body
type Query struct {
	Term string `json:"term"`
	// StartTime string `json:"start_time"`
	// EndTime   string `json:"end_time"`
}

// Search sends a search request to the Zinc database
// Search sends a search request to the Zinc database
func (zc *ZincClient) Search(index, term, searchType, field string, from, maxResults int, sourceFields, sortFields []string) ([]byte, error) {
	// Prepare the search request body
	searchRequest := SearchRequest{
		SearchType: searchType,
		Query: Query{
			Term: term,
		},
		SortFields: sortFields,
		From:       from,
		MaxResults: maxResults,
		Source:     sourceFields,
	}

	// Marshal the search request into JSON
	body, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search request: %w", err)
	}

	// Construct the endpoint for the search request
	endpoint := fmt.Sprintf("/api/%s/_search", index)

	// Use the newRequest method to create the HTTP request
	req, err := zc.newRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	// Send the request using the HTTP client
	resp, err := zc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send search request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search request failed with status %d: %s", resp.StatusCode, responseBody)
	}

	return responseBody, nil
}
