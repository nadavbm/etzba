package worker

import (
	"encoding/json"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/sqlclient"
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
func SetAPIAssignmentsToWorkers(data []byte) ([]Assignment, error) {
	var assignments []Assignment
	if err := json.Unmarshal(data, &assignments); err != nil {
		return nil, err
	}
	return assignments, nil
}

//
// ----------------------------------------------------------------- helpers ------------------------------------------------------------------------
//

func translateAssignmentToQueryBuilder(assignment *Assignment) *sqlclient.QueryBuilder {
	return &sqlclient.QueryBuilder{
		Command:    assignment.SqlQuery.Command,
		Table:      assignment.SqlQuery.Table,
		Constraint: assignment.SqlQuery.Constraint,
		ColumnsRef: assignment.SqlQuery.ColumnsRef,
		Values:     assignment.SqlQuery.Values,
	}
}

func translateAssignmentToAPIRequest(assignment *Assignment) *apiclient.ApiRequest {
	return &apiclient.ApiRequest{
		Url:             assignment.ApiRequest.Url,
		Method:          assignment.ApiRequest.Method,
		Payload:         assignment.ApiRequest.Payload,
		EndpointFile:    assignment.ApiRequest.EndpointFile,
		EndpointPattern: assignment.ApiRequest.EndpointPattern,
		Weight:          assignment.ApiRequest.Weight,
	}
}
