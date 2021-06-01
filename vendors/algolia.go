package vendors

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
	"github.com/knightazura/utils"
	"os"
)

type Algolia struct {
	Logger *utils.Logger
	Client *search.Client
	Settings *search.Settings
}

func InitAlgolia() contracts.SearchEngine {
	appId := os.Getenv("ALGOLIA_APP_ID")
	apiKey := os.Getenv("ALGOLIA_API_KEY")

	return &Algolia{
		Logger: utils.InitLogger(),
		Client: search.NewClient(appId, apiKey),
		Settings: &search.Settings{
			SearchableAttributes: opt.SearchableAttributes(
				"data.title",
				"data.content",
				"data.tags",
				),
		},
	}
}

func (a *Algolia) Add(docs *domain.GeneralDocuments, indexName string) {
	index := a.Client.InitIndex(indexName)

	// Do the job
	_, err := index.SaveObjects(
		*docs,
		opt.AutoGenerateObjectIDIfNotExist(true),
		opt.ExposeIntermediateNetworkErrors(true),
		)
	if err != nil {
		a.Logger.LogError("Failed to index documents to Algolia: %s", err.Error())
	} else {
		a.Logger.LogAccess("%s Algolia index created successfully", indexName)
	}
}

func (a *Algolia) Search(indexName string, query string) (result domain.SearchedDocument) {
	index := a.Client.InitIndex(indexName)
	_, err := index.SetSettings(*a.Settings)
	if err != nil {
		a.Logger.LogError("Failed to set Algolia index configuration: %s", err.Error())
	}

	res, err := index.Search(query)
	if err != nil {
		a.Logger.LogError("Failed to search Algolia objects: %s", err.Error())
	} else {
		var hits []interface{}
		for _, h := range res.Hits {
			hits = append(hits, h)
		}

		result = domain.SearchedDocument{
			Hits: hits,
			Limit: int64(res.HitsPerPage),
			Offset: int64(res.Page),
			TotalHits: int64(res.NbHits),
			Query: res.Query,
		}
	}
	return
}

func (a *Algolia) DeleteIndex(indexName string) {
	index := a.Client.InitIndex(indexName)
	_, err := index.Delete(indexName)
	if err != nil {
		a.Logger.LogError("Failed to delete Algolia index: %s", err.Error())
	}
}