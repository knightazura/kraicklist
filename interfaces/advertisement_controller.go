package interfaces

import (
	"bytes"
	"encoding/json"
	"github.com/knightazura/domain"
	"github.com/knightazura/services"
	"github.com/knightazura/usecases"
	"net/http"

	"github.com/knightazura/contracts"
)

type Advertisement struct {
	Seeder *services.Seeder
	AdvertisementInteractor usecases.AdvertisementInteractor
}

// TO DO: to pass config value as parameter
func InitAdvertisementController(se contracts.SearchEngine, seeder *services.Seeder) *Advertisement {
	return &Advertisement{
		Seeder: seeder,
		AdvertisementInteractor: usecases.AdvertisementInteractor{
			AdvertisementRepository: &AdvertisementRepository{},
			IndexedDocumentRepository: &IndexedDocumentRepository{
				SearchEngine: se,
			},
		},
	}
}

func (controller *Advertisement) Store(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Method not allowed"))
		return
	}

	// Parse the request data
	var payload domain.Advertisement
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &payload)

	controller.AdvertisementInteractor.Store(payload)

	writer.WriteHeader(http.StatusOK)
	writer.Write(buf.Bytes())
}

// This should be a route handler function
// but for challenge purpose, it handle loaded data from fs
func (controller *Advertisement) Upload() {
	// Load ads data from file
	adDocs := controller.Seeder.LoadData("./data.gz")
	if adDocs != nil {
		controller.AdvertisementInteractor.Upload(*adDocs)
	}

}

func (controller *Advertisement) Search(writer http.ResponseWriter, req *http.Request) {
	//context := req.Context()

	// Process the request
	query := req.URL.Query().Get("q")

	if len(query) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("missing search query in query params"))
		return
	}

	// Get relevant records
	records := controller.AdvertisementInteractor.Search(query)
	// if err != nil {
	// 	writer.WriteHeader(http.StatusInternalServerError)
	// 	writer.Write([]byte(err.Error()))
	// 	return
	// }

	// output success response
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(records)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(buf.Bytes())
}