package usecases

import "github.com/knightazura/domain"

// An indexed document repository belong to the usecases layer
type IndexedDocumentRepository interface {
	SearchDocs(query string, indexName string) domain.SearchedDocument
	ToIndexedDocument(docs domain.GeneralDocuments, indexName string)
}