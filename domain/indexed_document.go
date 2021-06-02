package domain

/** Constructor for any formatted documents by any vendors */

// Response format of searched documents
type SearchedDocument struct {
	Hits []interface{} `json:"hits"`
	Offset int64 `json:"offset"`
	Limit int64 `json:"limit"`
	TotalHits int64 `json:"total_hits"`
	Query string `json:"query"`
}

// Common format for search engine documents
type GeneralDocument struct {
	ID        int64    `json:"id"`
	Data interface{} `json:"data"`
}

type GeneralDocuments []GeneralDocument

// Meilisearch document
type MeilisearchDocument []map[string]interface{}

// Algolia document struct
type AlgoliaDocument struct {
	ObjectID string `json:"objectID"`
	Data interface{} `json:"data"`
}