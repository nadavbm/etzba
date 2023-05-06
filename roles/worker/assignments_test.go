package worker

import (
	"testing"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/sqlclient"
)

func TestAssignmentTranslation(t *testing.T) {
	apiRequests := []apiclient.ApiRequest{
		{
			Url:    "http://yeah.yeah.yeahs",
			Method: "POST",
		},
		{
			Url:    "https://to.to.com",
			Method: "GET",
		},
		{
			Url:    "http://yeah.yeah.yeahs",
			Method: "GET",
		},
	}

	assignments := []Assignment{}

	for _, r := range apiRequests {
		a := Assignment{
			ApiRequest: r,
		}
		assignments = append(assignments, a)
	}

	req := translateAssignmentToAPIRequest(&assignments[0])
	if req.Method != "POST" {
		t.Errorf("expected method to be POST, instead got %s", req.Method)
	}

	if req.Url != "http://yeah.yeah.yeahs" {
		t.Errorf("expected url to be http://yeah.yeah.yeahs, instead got %s", req.Url)
	}

	req = translateAssignmentToAPIRequest(&assignments[1])
	if req.Method != "GET" {
		t.Errorf("expected method to be POST, instead got %s", req.Method)
	}

	if req.Url != "https://to.to.com" {
		t.Errorf("expected url to be https://to.to.com, instead got %s", req.Url)
	}

	sqlQueries := []sqlclient.QueryBuilder{
		{
			Command:    "SELECT",
			Table:      "results",
			Constraint: "avg_query_duration BETWEEN 13.0 AND 15.0",
		},
		{
			Command:    "SELECT",
			Table:      "results",
			Constraint: "min_query_duration BETWEEN 50.0 AND 60.0",
		},
		{
			Command:    "SELECT",
			Table:      "results",
			Constraint: "total_queries BETWEEN 100 AND 200",
		},
	}

	assignments = []Assignment{}

	for _, q := range sqlQueries {
		a := Assignment{
			SqlQuery: q,
		}
		assignments = append(assignments, a)
	}

	builder := translateAssignmentToQueryBuilder(&assignments[0])
	if builder.Command != "SELECT" {
		t.Errorf("expected command to be SELECT, instead got %s", builder.Command)
	}

	if builder.Constraint != "avg_query_duration BETWEEN 13.0 AND 15.0" {
		t.Errorf("expected constraint to be avg_query_duration BETWEEN 13.0 AND 15.0, instead got %s", builder.Constraint)
	}

	builder = translateAssignmentToQueryBuilder(&assignments[2])
	if builder.Command != "SELECT" {
		t.Errorf("expected command to be SELECT, instead got %s", builder.Command)
	}

	if builder.Constraint != "total_queries BETWEEN 100 AND 200" {
		t.Errorf("expected constraint to be total_queries BETWEEN 100 AND 200, instead got %s", builder.Constraint)
	}
}
