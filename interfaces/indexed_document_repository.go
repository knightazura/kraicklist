package interfaces

import (
	"github.com/knightazura/domain"
	"github.com/knightazura/services"
)

type IndexedDocumentRepository struct {
	SearchEngine *services.SearchEngineHandler
}

func (id *IndexedDocumentRepository) SearchDocs(query string, indexName string) domain.SearchedDocument {
	// Deciding search engine vendor happened here
	docs := id.SearchEngine.PerformSearch(query, indexName)

	return docs
}

// Convert general document to meilisearch document
func (id *IndexedDocumentRepository) IndexDocs(docs *domain.GeneralDocuments, indexName string) {
	id.SearchEngine.IndexDocuments(docs, indexName)

	return
}