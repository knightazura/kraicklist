package usecases

import (
	"log"

	"github.com/knightazura/domain"
)

type AdvertisementInteractor struct {
	AdvertisementRepository   AdvertisementRepository
	IndexedDocumentRepository IndexedDocumentRepository
}

const EntityName = "advertisement"

func (adInteractor *AdvertisementInteractor) Store(payload domain.Advertisement) *domain.Advertisement {
	// Add to database
	newAd, newDoc := adInteractor.AdvertisementRepository.Store(&payload)
	log.Printf(`New ad: "%s" has been stored successfully`, newAd.Title)

	// Index the new entity
	adInteractor.ConvertToIndexedDocument(*newDoc)

	return newAd
}

func (adInteractor *AdvertisementInteractor) Search(query string) domain.SearchedDocument {
	docs := adInteractor.IndexedDocumentRepository.SearchDocs(query, EntityName)
	return docs
}

func (adInteractor *AdvertisementInteractor) Upload(ads domain.Advertisements) (newAds *domain.Advertisements, docs *domain.GeneralDocuments) {
	// Add to database
	newAds, docs = adInteractor.AdvertisementRepository.BulkStore(&ads)

	// Index the new entities
	adInteractor.IndexedDocumentRepository.BulkIndexDocs(docs, EntityName)
	return
}

// Convert advertisement data to search engine document
// Should add context as first parameter
func (adInteractor *AdvertisementInteractor) ConvertToIndexedDocument(doc domain.GeneralDocument) {
	adInteractor.IndexedDocumentRepository.IndexDocs(&doc, EntityName)
}
