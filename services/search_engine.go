package services

import (
	"os"

	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
	"github.com/knightazura/vendors"
)

type SearchEngineHandler struct {
	IndexName string
	Client    contracts.SearchEngine
	// another client can be added here
}

func InitSearchEngine() (*SearchEngineHandler, error) {
	seHandler := &SearchEngineHandler{}

	// Create (any) search engine instances here
	active := os.Getenv("SEARCH_ENGINE_ACTIVE")
	var engine contracts.SearchEngine
	switch active {
	case "meilisearch":
		engine = vendors.InitMeilisearch()
	case "algolia":
		engine = vendors.InitAlgolia()
	default:
		engine = vendors.InitRiot()
	}
	seHandler.Client = engine

	return seHandler, nil
}

func (se *SearchEngineHandler) IndexDocuments(doc *domain.GeneralDocument, indexName string) {
	se.Client.Add(doc, indexName)
}

func (se *SearchEngineHandler) PerformSearch(query string, indexName string) (result domain.SearchedDocument) {
	result = se.Client.Search(indexName, query)
	return
}

func (se *SearchEngineHandler) DeleteDocument(docID string, indexName string) {
	se.Client.DeleteDocument(docID, indexName)
	return
}

func (se *SearchEngineHandler) DeleteIndex(indexName string) {
	se.Client.DeleteIndex(indexName)
	return
}

func (se *SearchEngineHandler) TotalDocuments(indexName string) int64 {
	return se.Client.TotalDocuments(indexName)
}