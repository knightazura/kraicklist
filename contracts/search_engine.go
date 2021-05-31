package contracts

import (
	"github.com/knightazura/domain"
)

type SearchEngine interface {
	Search(query string, indexName string) domain.SearchedDocument
	Add(docs *domain.GeneralDocuments, indexName string)
	DeleteIndex(indexName string)
}
