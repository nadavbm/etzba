package worker

import (
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/sqlclient"
)

// Assignment for a Worker in order to create a db query and measure the time it takes
type Assignment struct {
	ApiRequest apiclient.ApiRequest   `json:"apiRequest" yaml:"apiRequests"`
	SqlQuery   sqlclient.QueryBuilder `json:"sqlQuery" yaml:"sqlQueries"`
}

// translateAssignmentToQueryBuilder takes a worker assignment and prepare sql query from it
func translateAssignmentToQueryBuilder(assignment *Assignment) *sqlclient.QueryBuilder {
	return &sqlclient.QueryBuilder{
		Command:    assignment.SqlQuery.Command,
		Table:      assignment.SqlQuery.Table,
		Constraint: assignment.SqlQuery.Constraint,
		ColumnsRef: assignment.SqlQuery.ColumnsRef,
		Values:     assignment.SqlQuery.Values,
	}
}

// translateAssignmentToAPIRequest takes a worker assignment and prepare api request from it
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
