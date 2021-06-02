package vendors

import (
	"github.com/knightazura/utils"
	"os"

	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
	"github.com/meilisearch/meilisearch-go"
)

type Meilisearch struct {
	Logger   *utils.Logger
	Client meilisearch.ClientInterface
}

func InitMeilisearch() contracts.SearchEngine {
	host := os.Getenv("MEILISEARCH_HOST")
	port := os.Getenv("MEILISEARCH_PORT")

	return &Meilisearch{
		Logger: utils.InitLogger(),
		Client: meilisearch.NewClient(meilisearch.Config{
			Host: host + ":" + port,
			//APIKey: "masterkey",
		}),
	}
}

func (m *Meilisearch) Add(doc *domain.GeneralDocument, indexName string) {
	get, _ := m.Client.Indexes().Get(indexName)

	// Create the index if it's not there
	if get == nil {
		_, err := m.Client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID: indexName,
		})

		if err != nil {
			m.Logger.LogError("Meilisearch: Failed to create index of %s: %v", indexName, err)
			return
		}
	}
	meiliDoc := domain.MeilisearchDocument{
		{
			"id": doc.ID,
			"data": doc.Data,
		},
	}

	// Do the job
	_, err := m.Client.Documents(indexName).AddOrUpdate(meiliDoc)
	if err != nil {
		m.Logger.LogError("Meilisearch: Failed to add %s documents: %v", indexName, err)
		return
	}

	m.Logger.LogAccess("Meilisearch: %s index created successfully", indexName)
	return
}

func (m *Meilisearch) BulkInsert(docs *domain.GeneralDocuments, indexName string) {
	m.createIndex(indexName)

	var meiliDocs domain.MeilisearchDocument
	for _, doc := range *docs {
		meiliDocs = append(meiliDocs, map[string]interface{}{
			"id": doc.ID,
			"data": doc.Data,
		})
	}

	// Do the job
	_, err := m.Client.Documents(indexName).AddOrUpdate(meiliDocs)
	if err != nil {
		m.Logger.LogError("Meilisearch: Failed to add %s documents: %v", indexName, err)
		return
	}

	m.Logger.LogAccess("Meilisearch: %s index created successfully", indexName)
	return
}

func (m *Meilisearch) Search(indexName string, query string) (result domain.SearchedDocument) {
	limit := int64(10)
	res, err := m.Client.Search(indexName).Search(meilisearch.SearchRequest{
		Query:  query,
		Limit:  limit,
		Offset: 0,
	})
	if err != nil {
		m.Logger.LogError("Meilisearch: Failed to search document %v", err)
		return
	}

	result = domain.SearchedDocument{
		Hits:      res.Hits,
		Limit:     limit,
		Offset:    0,
		TotalHits: res.NbHits,
		Query:     query,
	}

	m.Logger.LogAccess("Meilisearch: Search %s in %s documents. Found %d", query, indexName, len(res.Hits))
	return
}

func (m *Meilisearch) DeleteDocument(docID string, indexName string) {
	_, err := m.Client.Documents(indexName).Delete(docID)
	if err != nil {
		m.Logger.LogError("Meilisearch: Failed to delete Meilisearch document: %s", err.Error())
	}
}

func (m *Meilisearch) DeleteIndex(indexName string) {
	_, err := m.Client.Indexes().Delete(indexName)
	if err != nil {
		m.Logger.LogError("Meilisearch: Failed to delete Meilisearch index: %s", err.Error())
	}
}

func (m *Meilisearch) TotalDocuments(indexName string) int64 {
	res, err := m.Client.Stats().Get(indexName)
	if err != nil {
		m.Logger.LogError("Meilisearch: Failed to get total documents of %s: %v", indexName, err)
		return int64(0)
	} else {
		return res.NumberOfDocuments
	}
}

func (m *Meilisearch) createIndex(indexName string) {
	get, _ := m.Client.Indexes().Get(indexName)

	// Create the index if it's not there
	if get == nil {
		_, err := m.Client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID: indexName,
		})

		if err != nil {
			m.Logger.LogError("Meilisearch: Failed to create index of %s: %v", indexName, err)
			return
		}
	}
}

func (e *Meilisearch) TotalDocuments(indexName string) int64 {
	res, _ := e.Client.Stats().Get(indexName)
	return res.NumberOfDocuments
}
