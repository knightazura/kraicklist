package contracts

import (
	"github.com/knightazura/domain"
)

type SearchEngine interface {
	Search(indexName string, query string) domain.SearchedDocument
	Add(doc *domain.GeneralDocument, indexName string)
	BulkInsert(docs *domain.GeneralDocuments, indexName string)
	DeleteDocument(docID string, indexName string)
	DeleteIndex(indexName string)
	TotalDocuments(indexName string) int64
}
