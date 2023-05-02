package worker

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/sqlclient"

	"go.uber.org/zap"
)

// Assignment for a Worker in order to create a db query and measure the time it takes
type Assignment struct {
	ApiRequest apiclient.ApiRequest   `json:"apiRequest"`
	SqlQuery   sqlclient.QueryBuilder `json:"sqlQuery"`
}

// SetSQLAssignmentsToWorkers will take the csv output and create assignments for worker
func SetSQLAssignmentsToWorkers(data [][]string) []Assignment {
	var assignments []Assignment
	for i, line := range data {
		if i > 0 {
			var a Assignment
			for c, field := range line {
				switch {
				case c == 0:
					{
						a.SqlQuery.Command = field
					}
				case c == 1:
					{
						a.SqlQuery.Table = field
					}
				case c == 2:
					{
						a.SqlQuery.Constraint = field
					}
				case c == 3:
					{
						a.SqlQuery.ColumnsRef = field
					}
				case c == 4:
					{
						a.SqlQuery.Values = field
					}
				}
			}
			assignments = append(assignments, a)
		}
	}
	return assignments
}

// SetAPIAssignmentsToWorkers will take the csv output and create assignments for worker
func SetAPIAssignmentsToWorkers(data [][]string) ([]Assignment, error) {
	// TODO: remove data and read directly json array
	var assignments []Assignment
	for i, line := range data {
		if i > 0 {
			var a Assignment
			for c, field := range line {
				switch {
				case c == 0:
					{
						a.ApiRequest.Url = field
					}
				}
			}
		}
	}
	return assignments, nil
}

//
// ---------------------------------------------------------------------------------------- csv reader ------------------------------------------------------------------------------
//

// ReadCSVFile get a csv file, use csv reader and retrun byte
func ReadCSVFile(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("failed to close json file", zap.Error(err))
		}
	}()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("failed to read data from csv file", zap.Error(err))
		return nil, err
	}

	return data, nil
}

// ReadJSONFile get a json file and return byte slice
func ReadJSONFile(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Fatal("failed to close json file", zap.Error(err))
		}
	}()

	return ioutil.ReadAll(jsonFile)
}
