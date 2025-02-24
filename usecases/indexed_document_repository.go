package usecases

import "github.com/knightazura/domain"

// An indexed document repository belong to the usecases layer
type IndexedDocumentRepository interface {
	SearchDocs(query string, indexName string) domain.SearchedDocument
	IndexDocs(doc *domain.GeneralDocument, indexName string)
	BulkIndexDocs(docs *domain.GeneralDocuments, indexName string)
	GetTotalDocuments(indexName string) int64
}
