package services

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"github.com/knightazura/domain"
	"log"
	"os"
)

// A service to load and manage model that need
// to be formatted into common type, GeneralDocument
// Note: This service only for testing purpose
type Seeder struct {}

func InitSeeder() *Seeder {
	return &Seeder{}
}

// Load advertisement model data from file
func (s *Seeder) LoadData(path string) (out *domain.Advertisements){
	// Open file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("unable to open source file due: %v", err)
		return nil
	}
	defer file.Close()

	// read as gzip
	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatalf("unable to initialize gzip reader due: %v", err)
		return nil
	}

	// read the reader using scanner to contstruct records
	cs := bufio.NewScanner(reader)
	var o domain.Advertisements
	for cs.Scan() {
		var model domain.Advertisement
		err = json.Unmarshal(cs.Bytes(), &model)
		if err != nil {
			continue
		}
		o = append(o, model)
	}

	return &o
}