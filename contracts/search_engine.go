package contracts

import (
	"github.com/knightazura/domain"
)

type SearchEngine interface {
	Search(query string, indexName string) domain.SearchedDocument
	Add(docs *domain.GeneralDocuments, indexName string)
	DeleteDocument(docID string, indexName string)
	DeleteIndex(indexName string)
	TotalDocuments(indexName string) int64
}
