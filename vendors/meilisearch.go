package vendors

import (
	"fmt"
	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"os"
)

type Meilisearch struct {
	Client meilisearch.ClientInterface
}

func InitMeilisearch() contracts.SearchEngine {
	host := os.Getenv("MEILISEARCH_HOST")
	port := os.Getenv("MEILISEARCH_PORT")

	return &Meilisearch{
		Client: meilisearch.NewClient(meilisearch.Config{
			Host: host + ":" + port,
			//APIKey: "masterkey",
		}),
	}
}

func (m *Meilisearch) Add(docs *domain.GeneralDocuments, indexName string) {
	get, _ := m.Client.Indexes().Get(indexName)

	// Create the index if it's not there
	if get == nil {
		_, err := m.Client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID: indexName,
		})

		if err != nil {
			log.Fatalf("Failed to create index of %s: %v", indexName, err)
			return
		}
	}

	var documents domain.GeneralDocuments
	for _, doc := range *docs {
		documents = append(documents, doc)
	}

	_, err := m.Client.Documents(indexName).AddOrUpdate(documents)
	if err != nil {
		log.Fatalf("Failed to add %s documents: %v", indexName, err)
		return
	}
	fmt.Printf("%s index created successfully\n", indexName)
}

func (m *Meilisearch) Search(indexName string, query string) (result domain.SearchedDocument) {
	limit := int64(10)
	res, _ := m.Client.Search(indexName).Search(meilisearch.SearchRequest{
		Query:  query,
		Limit:  limit,
		Offset: 0,
	})

	result = domain.SearchedDocument{
		Hits: res.Hits,
		Limit: limit,
		Offset: 0,
		TotalHits: res.NbHits,
		Query: query,
	}
	return
}

func (m *Meilisearch) DeleteDocument(docID string, indexName string) {
	_, err := m.Client.Documents(indexName).Delete(docID)
	if err != nil {
		log.Printf("Failed to delete Meilisearch document: %s", err.Error())
	}
}

func (m *Meilisearch) DeleteIndex(indexName string) {
	_, err := m.Client.Indexes().Delete(indexName)
	if err != nil {
		log.Printf("Failed to delete Meilisearch index: %s", err.Error())
	}
}
