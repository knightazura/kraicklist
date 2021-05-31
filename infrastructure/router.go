package infrastructure

import (
	"fmt"
	"github.com/knightazura/interfaces"
	"github.com/knightazura/services"
	"github.com/knightazura/utils"
	"net/http"
	"os"
)

type Services struct{
	Seeder *services.Seeder
	SearchEngine *services.SearchEngineHandler
}

func Dispatch(logger *utils.Logger) {
	services := setupServices()
	setupServer(logger, services)
}

func setupServices() *Services {
	// Search engine
	searchEngine, _ := services.InitSearchEngine()

	return &Services{
		Seeder: &services.Seeder{},
		SearchEngine: searchEngine,
	}
}

func setupServer(logger *utils.Logger, services *Services) {
	// Handle static files for frontend
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	adController := interfaces.InitAdvertisementController(logger, services.SearchEngine, services.Seeder)

	// Advertisement routes
	http.HandleFunc("/advertisement/search", adController.Search)
	http.HandleFunc("/advertisement", adController.Store)
	// Challenge purpose: mock of /advertisement/upload endpoint
	if os.Getenv("SEARCH_ENGINE_ACTIVE") == "meilisearch" {
		adController.Upload()
	}

	// Setup and start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3001"
	}
	fmt.Printf("Server is listening on %s...", port)
	
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		logger.LogError("unable to start server due: %s", err.Error())
	}
}