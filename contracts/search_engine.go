package contracts

import (
	"github.com/knightazura/domain"
)

type SearchEngine interface {
	PerformSearch(query string, indexName string) domain.SearchedDocument
	IndexDocuments(docs domain.GeneralDocuments, indexName string)
}
