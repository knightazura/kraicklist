package vendors

import (
	"fmt"
	"github.com/knightazura/domain"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"os"
)

type Meilisearch struct {
	Client meilisearch.ClientInterface
}

func InitMeilisearch(mode string) *Meilisearch {
	port := os.Getenv("SEARCH_ENGINE_PORT")
	if mode == "test" {
		port = os.Getenv("SEARCH_ENGINE_TEST_PORT")
	}
	config := meilisearch.Config{
		Host: os.Getenv("SEARCH_ENGINE_HOST") + ":" + port,
		//APIKey: "masterkey",
	}
	client := meilisearch.NewClient(config)

	return &Meilisearch{
		Client: client,
	}
}

func MSAddDocuments(client meilisearch.ClientInterface, docs domain.GeneralDocuments, indexName string) {
	get, _ := client.Indexes().Get(indexName)

	// Create the index if it's not there
	if get == nil {
		_, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID: indexName,
		})

		if err != nil {
			log.Fatalf("Failed to create index of %s: %v", indexName, err)
			return
		}
	}

	var documents domain.GeneralDocuments
	for _, doc := range docs {
		documents = append(documents, doc)
	}

	_, err := client.Documents(indexName).AddOrUpdate(documents)
	if err != nil {
		log.Fatalf("Failed to add %s documents: %v", indexName, err)
		return
	}
	fmt.Printf("%s index created successfully\n", indexName)
}

func MSSearch(client meilisearch.ClientInterface, indexName string, query string) (result domain.SearchedDocument) {
	limit := int64(10)
	res, _ := client.Search(indexName).Search(meilisearch.SearchRequest{
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
