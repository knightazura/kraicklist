package vendors

import (
	"os"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
	"github.com/knightazura/utils"
)

type Algolia struct {
	Logger   *utils.Logger
	Client   *search.Client
	Settings *search.Settings
}

func InitAlgolia() contracts.SearchEngine {
	appId := os.Getenv("ALGOLIA_APP_ID")
	apiKey := os.Getenv("ALGOLIA_API_KEY")

	return &Algolia{
		Logger: utils.InitLogger(),
		Client: search.NewClient(appId, apiKey),
		Settings: &search.Settings{
			// Search can be improved with this option
			// CustomRanking: opt.CustomRanking(),
			SearchableAttributes: opt.SearchableAttributes(
				"data.title",
				"data.content",
				"data.tags",
			),
		},
	}
}

func (a *Algolia) Add(doc *domain.GeneralDocument, indexName string) {
	index := a.Client.InitIndex(indexName)

	// Assign to following algolia document structure
	algoDoc := domain.AlgoliaDocument{
		ObjectID: strconv.FormatInt(doc.ID, 10),
		Data:     doc.Data,
	}

	// Do the job
	_, err := index.SaveObject(
		algoDoc,
		opt.AutoGenerateObjectIDIfNotExist(true),
		opt.ExposeIntermediateNetworkErrors(true),
	)
	if err != nil {
		a.Logger.LogError("Algolia: Failed to index documents to Algolia: %s", err.Error())
		return
	}

	a.Logger.LogAccess("Algolia: %s index created successfully", indexName)
	return
}

func (a *Algolia) BulkInsert(docs *domain.GeneralDocuments, indexName string) {
	index := a.createIndex(indexName)

	if index == nil {
		a.Logger.LogError("Algolia: Failed to do bulk insert %s documents, due to fail to create index", indexName)
		return
	}

	// Assign to following algolia document structure
	var algoDocs []domain.AlgoliaDocument
	for _, doc := range *docs {
		algoDocs = append(algoDocs, domain.AlgoliaDocument{
			ObjectID: strconv.FormatInt(doc.ID, 10),
			Data:     doc.Data,
		})
	}

	// Do the job
	_, err := index.SaveObjects(
		algoDocs,
		opt.AutoGenerateObjectIDIfNotExist(true),
		opt.ExposeIntermediateNetworkErrors(true),
	)
	if err != nil {
		a.Logger.LogError("Algolia: Failed to bulk insert documents to Algolia: %v", err)
		return
	}

	a.Logger.LogAccess("Algolia: Successfully do bulk insert of %s index", indexName)
	return
}

func (a *Algolia) Search(indexName string, query string) (result domain.SearchedDocument) {
	index := a.Client.InitIndex(indexName)
	_, err := index.SetSettings(*a.Settings)
	if err != nil {
		a.Logger.LogError("Algolia: Failed to set Algolia index configuration: %s", err.Error())
	}

	res, err := index.Search(query)
	if err != nil {
		a.Logger.LogError("Algolia: Failed to search Algolia objects: %s", err.Error())
	}
	var hits []interface{}
	for _, h := range res.Hits {
		hits = append(hits, h)
	}

	result = domain.SearchedDocument{
		Hits:      hits,
		Limit:     int64(res.HitsPerPage),
		Offset:    int64(res.Page),
		TotalHits: int64(res.NbHits),
		Query:     res.Query,
	}

	a.Logger.LogAccess("Algolia: Search %s in %s documents. Found %d", query, indexName, len(hits))
	return
}

func (a *Algolia) DeleteDocument(docID string, indexName string) {
	index := a.Client.InitIndex(indexName)
	_, err := index.DeleteObject(docID)
	if err != nil {
		a.Logger.LogError("Algolia: Failed to delete Algolia document: %s", err.Error())
	}
}

func (a *Algolia) DeleteIndex(indexName string) {
	index := a.Client.InitIndex(indexName)
	_, err := index.Delete(indexName)
	if err != nil {
		a.Logger.LogError("Algolia: Failed to delete Algolia index: %s", err.Error())
	}
}

func (a *Algolia) TotalDocuments(indexName string) int64 {
	totalDocuments := int64(0)

	res, err := a.Client.ListIndices()
	if err != nil {
		a.Logger.LogError("Algolia: Failed to get total documents of %s: %v", indexName, err)
		return totalDocuments
	}

	for _, item := range res.Items {
		if item.Name == indexName {
			totalDocuments = item.Entries
			break
		}
	}

	return totalDocuments
}

func (a *Algolia) createIndex(indexName string) *search.Index {
	index := a.Client.InitIndex(indexName)
	ok, err := index.Exists()
	if err != nil {
		a.Logger.LogError("Algolia: Failed to check %s index existence or create it: %v", indexName, err)
		return nil
	}
	if ok {
		a.Logger.LogAccess("Algolia: Create %s index successfully", index)
		return index
	}

	return nil
}