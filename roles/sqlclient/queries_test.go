package sqlclient

import "testing"

func TestToSQL(t *testing.T) {
	builders := []QueryBuilder{
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

	query := ToSQL(&builders[0])
	if query != "SELECT * FROM results WHERE avg_query_duration BETWEEN 13.0 AND 15.0" {
		t.Errorf("expected command to be SELECT * FROM results WHERE avg_query_duration BETWEEN 13.0 AND 15.0, instead got %s", query)
	}

	query = ToSQL(&builders[1])
	if query != "SELECT * FROM results WHERE min_query_duration BETWEEN 50.0 AND 60.0" {
		t.Errorf("expected command to be SELECT * FROM results WHERE min_query_duration BETWEEN 50.0 AND 60.0, instead got %s", query)
	}

	query = ToSQL(&builders[2])
	if query != "SELECT * FROM results WHERE total_queries BETWEEN 100 AND 200" {
		t.Errorf("expected command to be SELECT * FROM results WHERE total_queries BETWEEN 100 AND 200, instead got %s", query)
	}
}
