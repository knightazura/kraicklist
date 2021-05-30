package interfaces

import (
	"fmt"
	"github.com/knightazura/domain"
)

type AdvertisementRepository struct {}

func(ar *AdvertisementRepository) Store(payload *domain.Advertisement) (newAd *domain.Advertisement, newDoc *domain.GeneralDocument) {
	return payload, &domain.GeneralDocument{
		ID: payload.ID,
		Data: payload,
	}
}

func (ar *AdvertisementRepository) BulkStore(ads *domain.Advertisements) (newAds *domain.Advertisements, newDocs *domain.GeneralDocuments) {
	if len(*ads) == 0 {
		_ = fmt.Errorf("There's no ads that need to be converted to docs")
		return
	}

	var na domain.Advertisements
	var docs domain.GeneralDocuments
	for _, ad := range *ads {
		na = append(na, ad)
		docs = append(docs, domain.GeneralDocument{
			ID:   ad.ID,
			Data: ad,
		})
	}
	return &na, &docs
}