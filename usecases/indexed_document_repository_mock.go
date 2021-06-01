package usecases

import (
	"github.com/knightazura/domain"
	"github.com/stretchr/testify/mock"
	"log"
)

type IndexedDocumentRepositoryMock struct {
	Mock mock.Mock
}

func (repo *IndexedDocumentRepositoryMock) SearchDocs(query string, indexName string) domain.SearchedDocument {
	log.Println("implement me")
	return domain.SearchedDocument{}
}

func (repo *IndexedDocumentRepositoryMock) IndexDocs(docs *domain.GeneralDocuments, indexName string) {
	repo.Mock.Called(docs, indexName)
}

func (repo *IndexedDocumentRepositoryMock) GetTotalDocuments(indexName string) int64 {
	repo.Mock.Called(indexName)
	return int64(0)
}
