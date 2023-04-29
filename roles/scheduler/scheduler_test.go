package scheduler

import (
	"testing"
	"time"

	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/etzba/roles/worker"
	"github.com/nadavbm/zlog"
)

func TestScheduler(t *testing.T) {
	logger := zlog.New()
	duration := 3 * time.Second

	configFile := "auth.json"
	helperFile := "queries.csv"
	workers := 40
	assignments := []worker.Assignment{
		{
			SqlQuery: sqlclient.QueryBuilder{
				Command:    "SELECT",
				Table:      "results",
				Constraint: "avg_query_duration BETWEEN 13.0 AND 15.0",
			},
		}, {
			SqlQuery: sqlclient.QueryBuilder{
				Command:    "SELECT",
				Table:      "results",
				Constraint: "min_query_duration BETWEEN 50.0 AND 60.0",
			},
		}, {SqlQuery: sqlclient.QueryBuilder{
			Command:    "SELECT",
			Table:      "results",
			Constraint: "total_queries BETWEEN 100 AND 200",
		},
		},
	}

	scheduler, err := NewScheduler(logger, duration, configFile, helperFile, workers, true)
	if err != nil {
		t.Fatal("could not create an instance of scheduler")
	}

	workCh := scheduler.setWorkChannelForDuration(assignments)
	for a := range workCh {
		if a.SqlQuery.Command != "SELECT" {
			t.Error("expected command to be select instead gota", a.SqlQuery.Command)
		}
	}
}
