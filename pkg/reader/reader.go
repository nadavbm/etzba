package reader

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
)

// ReadCSVFile get a csv file, use csv reader and retrun byte
func ReadCSVFile(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("failed to close json file", err.Error())
		}
	}()

	csvReader := csv.NewReader(f)

	return csvReader.ReadAll()
}

// ReadJSONFile get a json file and return byte slice
func ReadJSONFile(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Println("failed to close json file", err.Error())
		}
	}()

	return ioutil.ReadAll(jsonFile)
}
