package interfaces

import (
	"encoding/json"
	"github.com/knightazura/domain"
	"github.com/knightazura/services"
	"github.com/knightazura/usecases"
	"github.com/knightazura/utils"
	"net/http"
)

type Advertisement struct {
	Seeder *services.Seeder
	AdvertisementInteractor usecases.AdvertisementInteractor
	Logger *utils.Logger
}

// TO DO: to pass config value as parameter
func InitAdvertisementController(logger *utils.Logger, se *services.SearchEngineHandler, seeder *services.Seeder) *Advertisement {
	return &Advertisement{
		Seeder: seeder,
		AdvertisementInteractor: usecases.AdvertisementInteractor{
			AdvertisementRepository: &AdvertisementRepository{},
			IndexedDocumentRepository: &IndexedDocumentRepository{
				SearchEngine: se,
			},
		},
		Logger: logger,
	}
}

func (controller *Advertisement) Store(writer http.ResponseWriter, req *http.Request) {
	response := utils.InitResponse(controller.Logger, writer)

	if req.Method != http.MethodPost {
		response.MethodNotAllowedResponse("Method not allowed")
	}

	// Parse the request data
	var payload domain.Advertisement
	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		controller.Logger.LogError("%s", "Failed to decode payload request")
	}

	controller.AdvertisementInteractor.Store(payload)

	response.OkResponse("Success store new ad", payload)
}

// This should be a route handler function
// but for challenge purpose, it handle loaded data from fs
func (controller *Advertisement) Upload() {
	// Load ads data from file
	adDocs := controller.Seeder.LoadData("./data.gz")
	if adDocs != nil {
		total := controller.AdvertisementInteractor.IndexedDocumentRepository.GetTotalDocuments("advertisement")
		// To avoid usage overload in algolia service
		if total < int64(len(*adDocs)) {
			controller.AdvertisementInteractor.Upload(*adDocs)
		}
	}
}

func (controller *Advertisement) Search(writer http.ResponseWriter, req *http.Request) {
	//context := req.Context()

	response := utils.InitResponse(controller.Logger, writer)

	// Process the request
	query := req.URL.Query().Get("q")

	if len(query) == 0 {
		response.BadRequestResponse("Missing search query in query params")
	}

	// Get relevant records
	records := controller.AdvertisementInteractor.Search(query)
	// if err != nil {
	// 	response.InternalServerErrorResponse("Cannot do search query: " + err.Error())
	// }

	// output success response
	response.OkResponse("Search success", records)
}