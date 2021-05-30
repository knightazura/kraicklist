package usecases

import (
	"github.com/knightazura/domain"
	"github.com/stretchr/testify/mock"
	"log"
)

type AdvertisementRepositoryMock struct {
	Mock mock.Mock
}

func (repo *AdvertisementRepositoryMock) Store(payload *domain.Advertisement) (*domain.Advertisement, *domain.GeneralDocument) {
	args := repo.Mock.Called(payload)
	newAd := args.Get(0).(*domain.Advertisement)
	newDoc := &domain.GeneralDocument{
		ID: newAd.ID,
		Data: newAd,
	}
	return newAd, newDoc
}

func (repo *AdvertisementRepositoryMock) BulkStore(ads *domain.Advertisements) (newAds *domain.Advertisements, newDocs *domain.GeneralDocuments) {
	log.Println("implement me")
	return nil, nil
}