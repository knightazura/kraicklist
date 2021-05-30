package usecases

import (
	"github.com/knightazura/domain"
	"log"
)

type AdvertisementInteractor struct {
	AdvertisementRepository AdvertisementRepository
	IndexedDocumentRepository IndexedDocumentRepository
}

const EntityName = "advertisement"

func (adInteractor *AdvertisementInteractor) Store(payload *domain.Advertisement) {
	// Add to database
	newAd, newDoc := adInteractor.AdvertisementRepository.Store(payload)
	log.Printf(`New ad: "%s" has been stored successfully`, newAd.Title)

	// Index the new entity
	adInteractor.ConvertToIndexedDocuments([]domain.GeneralDocument{newDoc})
}

func (adInteractor *AdvertisementInteractor) Search(query string) domain.SearchedDocument {
	docs := adInteractor.IndexedDocumentRepository.SearchDocs(query, EntityName)
	return docs
}

func (adInteractor *AdvertisementInteractor) Upload(ads domain.Advertisements) (newAds domain.Advertisements, docs domain.GeneralDocuments) {
	// Add to database
	newAds, docs = adInteractor.AdvertisementRepository.BulkStore(ads)

	// Index the new entities
	adInteractor.ConvertToIndexedDocuments(docs)
	return
}

// Convert advertisement data to search engine document
// Should add context as first parameter
func (adInteractor *AdvertisementInteractor) ConvertToIndexedDocuments(docs domain.GeneralDocuments) {
	adInteractor.IndexedDocumentRepository.IndexDocs(docs, EntityName)
	return
}