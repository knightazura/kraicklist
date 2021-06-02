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
	return id.SearchEngine.PerformSearch(query, indexName)
}

// Convert general document to meilisearch document
func (id *IndexedDocumentRepository) IndexDocs(doc *domain.GeneralDocument, indexName string) {
	id.SearchEngine.IndexDocuments(doc, indexName)

	return
}

func (id *IndexedDocumentRepository) BulkIndexDocs(docs *domain.GeneralDocuments, indexName string) {
	id.SearchEngine.Client.BulkInsert(docs, indexName)

	return
}

func (id *IndexedDocumentRepository) GetTotalDocuments(indexName string) int64 {
	return id.SearchEngine.TotalDocuments(indexName)
}
