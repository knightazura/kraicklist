package usecases

import (
	"github.com/knightazura/domain"
	"github.com/stretchr/testify/mock"
)

type IndexedDocumentRepositoryMock struct {
	Mock mock.Mock
}

func (repo *IndexedDocumentRepositoryMock) SearchDocs(query string, indexName string) domain.SearchedDocument {
	repo.Mock.Called(query, indexName)
	return domain.SearchedDocument{}
}

func (repo *IndexedDocumentRepositoryMock) IndexDocs(doc *domain.GeneralDocument, indexName string) {
	repo.Mock.Called(doc, indexName)
}

func (repo *IndexedDocumentRepositoryMock) BulkIndexDocs(docs *domain.GeneralDocuments, indexName string) {
	repo.Mock.Called(docs, indexName)
}

func (repo *IndexedDocumentRepositoryMock) GetTotalDocuments(indexName string) int64 {
	repo.Mock.Called(indexName)
	return int64(0)
}
