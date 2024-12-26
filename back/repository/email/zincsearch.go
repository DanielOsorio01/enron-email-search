package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DanielOsorio01/enron-email-search/back/db"
	"github.com/DanielOsorio01/enron-email-search/back/models"
)

type ZincsearchRepo struct {
	Client *db.ZincClient
}

// NewZincsearchRepo creates a new instance of ZincsearchRepo
func NewZincsearchRepo(client *db.ZincClient) *ZincsearchRepo {
	return &ZincsearchRepo{
		Client: client,
	}
}

// DefaultSearchParams returns default search parameters
func DefaultSearchParams() SearchParams {
	return SearchParams{
		SearchType:   SearchTypeMatch,
		From:         DefaultFrom,
		MaxResults:   DefaultMaxResults,
		SourceFields: []string{}, // Empty array returns all fields
		SortFields:   []string{"-@timestamp"},
	}
}

// Search performs a search operation and returns matching emails
func (r *ZincsearchRepo) Search(ctx context.Context, index string, params SearchParams) ([]models.Email, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	response, err := r.Client.Search(
		index,
		params.Term,
		params.SearchType,
		params.Field,
		params.From,
		params.MaxResults,
		params.SourceFields,
		params.SortFields,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}

	var searchResult SearchResult
	if err := json.Unmarshal(response, &searchResult); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	if searchResult.Error != "" {
		return nil, fmt.Errorf("search error: %s", searchResult.Error)
	}

	emails := make([]models.Email, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		emails = append(emails, hit.Source)
	}

	return emails, nil
}

// SearchWithDefaults performs a simple search with default parameters
func (r *ZincsearchRepo) SearchWithDefaults(ctx context.Context, index, term string) ([]models.Email, error) {
	params := DefaultSearchParams()
	params.Term = term
	return r.Search(ctx, index, params)
}

// CountResults returns the total number of matching documents
func (r *ZincsearchRepo) CountResults(ctx context.Context, index, term string) (int, error) {
	params := DefaultSearchParams()
	params.Term = term
	params.MaxResults = 0

	response, err := r.Client.Search(
		index,
		params.Term,
		params.SearchType,
		params.Field,
		params.From,
		params.MaxResults,
		params.SourceFields,
		params.SortFields,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to perform count search: %w", err)
	}

	var searchResult SearchResult
	if err := json.Unmarshal(response, &searchResult); err != nil {
		return 0, fmt.Errorf("failed to parse count results: %w", err)
	}

	return searchResult.Hits.Total.Value, nil
}
