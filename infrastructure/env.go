package infrastructure

import (
	"bufio"
	"github.com/knightazura/utils"
	"os"
	"strings"
)

// Load app configuration variables from an env file.
func Bootstrap(logger *utils.Logger) {
	filePath := ".env"

	envFile, err := os.Open(filePath)
	if err != nil {
		logger.LogError("Failed to open env file %s", err.Error())
	}
	defer envFile.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(envFile)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.LogError("Error reading env file %s", err.Error())
	}

	for _, line := range lines {
		pair := strings.Split(line, "=")
		os.Setenv(pair[0], pair[1])
	}
}