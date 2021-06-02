package main

import (
	"bytes"
	"encoding/json"
	"github.com/knightazura/domain"
	"github.com/knightazura/infrastructure"
	"github.com/knightazura/interfaces"
	"github.com/knightazura/services"
	"github.com/knightazura/usecases"
	"github.com/knightazura/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var advertisementRepository = &usecases.AdvertisementRepositoryMock{Mock: mock.Mock{}}
var indexedDocumentRepository = &usecases.IndexedDocumentRepositoryMock{Mock: mock.Mock{}}
var	advertisementInteractor = &usecases.AdvertisementInteractor{
	AdvertisementRepository:   advertisementRepository,
	IndexedDocumentRepository: indexedDocumentRepository,
}

type IntegrationTestSuite struct {
	suite.Suite
	EntityName string
	SearchEngine *services.SearchEngineHandler
	Logger *utils.Logger
}

func (suite *IntegrationTestSuite) SetupTest() {
	logger := utils.InitLogger()
	infrastructure.Bootstrap(logger)
	suite.EntityName = "advertisement"
	suite.SearchEngine, _ = services.InitSearchEngine()
	suite.Logger = logger
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateNewAd() {
	payload := domain.Advertisement{
		ID: int64(1),
		Title: randomString(20),
		Content: randomString(100),
		Tags: []string{randomString(3), randomString(4)},
	}
	gd := domain.GeneralDocument{
		ID: payload.ID,
		Data: &payload,
	}

	// Use case if payload is correct
	advertisementRepository.Mock.On("Store", &payload).Return(&payload, &gd)
	indexedDocumentRepository.Mock.On("IndexDocs", &gd, suite.EntityName).Return(nil)

	// Test
	// If payload is correct
	newAd := advertisementInteractor.Store(payload)
	suite.NotNil(newAd, "New add should be added")
	suite.Equal(payload.ID, newAd.ID)
	suite.Equal(payload.Title, newAd.Title)
}

func (suite *IntegrationTestSuite) TestCreateNewAdViaRequest() {
	payload := domain.Advertisement{
		ID: int64(2),
		Title: randomString(20),
		Content: randomString(100),
		Tags: []string{randomString(3), randomString(4)},
	}

	newAd := suite.hitStoreApi(payload)
	if newAd == nil {
		suite.Fail("Fail create new ad via API request")
	} else {
		suite.Equal(payload.ID, newAd.ID)
		suite.Equal(payload.Title, newAd.Title)
	}
}

func (suite *IntegrationTestSuite) TestSearchAd() {
	// Init server
	server := suite.createTestServer("search")
	defer server.Close()

	payload := domain.Advertisement{
		ID: int64(3),
		Title: randomString(20),
		Content: randomString(100),
		Tags: []string{randomString(3), randomString(4)},
	}

	// Create new ad that will be search in next step
	newAd := suite.hitStoreApi(payload)

	time.Sleep(2 * time.Second)

	// Search the new ad
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		suite.Fail("Cannot create search request to API")
	} else {
		q := req.URL.Query()
		q.Add("q", newAd.Title)

		req.URL.RawQuery = q.Encode()
		req.Header.Set("Content-Type", "application/json")

		res, err := server.Client().Do(req)
		if err != nil {
			suite.Fail("Cannot make search request to API")
		} else {
			var successResponse utils.SuccessResponse
			json.NewDecoder(res.Body).Decode(&successResponse)

			log.Println(successResponse)

			var searchResponse domain.SearchedDocument
			sr, _ := json.Marshal(successResponse.Data)
			json.Unmarshal(sr, &searchResponse)

			var doc domain.GeneralDocument
			d, err := json.Marshal(searchResponse.Hits[0])
			if err != nil {
				suite.Logger.LogError("Integration test: %s", err.Error())
			}
			json.Unmarshal(d, &doc)

			var docData domain.Advertisement
			dd, _ := json.Marshal(doc.Data)
			json.Unmarshal(dd, &docData)

			// Assertions
			suite.Equal(int64(1), searchResponse.TotalHits)
			suite.Equal(newAd.ID, docData.ID)
			suite.Equal(newAd.Title, docData.Title)
		}
	}
}

func (suite *IntegrationTestSuite) hitStoreApi(payload domain.Advertisement) *domain.Advertisement {
	// Cleanup search engine documents
	suite.resetTestDocument(payload.ID)

	server := suite.createTestServer("store")
	defer server.Close()

	p, _ := json.Marshal(payload)
	res, err := http.Post(server.URL, "application/json", bytes.NewBuffer(p))
	if err != nil {
		suite.FailNow("Failed do request: " + err.Error())
	}
	var successResponse utils.SuccessResponse
	json.NewDecoder(res.Body).Decode(&successResponse)

	var searchResponse domain.Advertisement
	srByte, _ := json.Marshal(successResponse.Data)
	json.Unmarshal(srByte, &searchResponse)

	return &searchResponse
}

func (suite *IntegrationTestSuite) createTestServer(endpoint string) *httptest.Server {
	seeder := &services.Seeder{}
	adController := interfaces.InitAdvertisementController(suite.Logger, suite.SearchEngine, seeder)

	if endpoint == "store" {
		return httptest.NewServer(http.HandlerFunc(adController.Store))
	} else {
		return httptest.NewServer(http.HandlerFunc(adController.Search))
	}
}

func (suite *IntegrationTestSuite) resetTestDocument(docID int64) {
	suite.SearchEngine.DeleteDocument(strconv.FormatInt(docID, 10), suite.EntityName)
}

func randomString(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	word := make([]rune, n)
	for i := range word {
		word[i] = chars[rand.Intn(len(chars))]
	}

	return string(word)
}
