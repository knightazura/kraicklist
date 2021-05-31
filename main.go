package main

import (
	"github.com/knightazura/infrastructure"
	"github.com/knightazura/utils"
)

func main() {
	logger := utils.InitLogger()

	// Load configuration
	infrastructure.Bootstrap(logger)

	// Dispatch the app
	infrastructure.Dispatch(logger)
}