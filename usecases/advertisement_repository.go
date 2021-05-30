package usecases

import "github.com/knightazura/domain"

type AdvertisementRepository interface {
	Store(payload *domain.Advertisement) (newAd *domain.Advertisement, newDoc *domain.GeneralDocument)
	BulkStore(ads *domain.Advertisements) (newAds *domain.Advertisements, newDocs *domain.GeneralDocuments)
}