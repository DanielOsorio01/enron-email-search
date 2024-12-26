package email

import (
	"time"

	"github.com/DanielOsorio01/enron-email-search/back/models"
)

// SearchResult represents the structure of the ZincSearch response
type SearchResult struct {
	Took     int         `json:"took"`
	TimedOut bool        `json:"timed_out"`
	MaxScore float64     `json:"max_score"`
	Hits     SearchHits  `json:"hits"`
	Buckets  interface{} `json:"buckets"`
	Error    string      `json:"error"`
}

// SearchHits contains the hits information from the search result
type SearchHits struct {
	Total Total           `json:"total"`
	Hits  []SearchHitItem `json:"hits"`
}

// Total represents the total hits count
type Total struct {
	Value int `json:"value"`
}

// SearchHitItem represents a single hit in the search results
type SearchHitItem struct {
	Index     string       `json:"_index"`
	Type      string       `json:"_type"`
	ID        string       `json:"_id"`
	Score     float64      `json:"_score"`
	Timestamp string       `json:"@timestamp"`
	Source    models.Email `json:"_source"`
}

// SearchParams contains the parameters for a search request
type SearchParams struct {
	Term         string
	Field        string
	SearchType   string
	From         int
	MaxResults   int
	StartTime    time.Time
	EndTime      time.Time
	SortFields   []string
	SourceFields []string
}
