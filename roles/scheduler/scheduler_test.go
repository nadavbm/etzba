package scheduler

import (
	"testing"
	"time"

	"github.com/nadavbm/etzba/roles/common"
	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/etzba/roles/worker"
	"github.com/nadavbm/zlog"
)

func TestRpsSet(t *testing.T) {
	logger := zlog.New()
	settings := common.Settings{
		Duration:        time.Duration(3 * time.Second),
		ExecutionType:   "api",
		ConfigFile:      "secret.json",
		HelpersFile:     "api.yaml",
		Rps:             10,
		NumberOfWorkers: 2,
		Verbose:         true,
	}
	s, err := NewScheduler(logger, &settings)
	if err != nil {
		t.Fatal("could not create scheduler instance")
	}

	rps := s.Settings.SetRps()
	if rps != time.Duration(100*time.Millisecond) {
		t.Errorf("expected rps to be 100ms but instead got %v", rps)
	}

	s.Settings.Rps = 200
	rps = s.Settings.SetRps()
	if rps != time.Duration(5*time.Millisecond) {
		t.Errorf("expected rps to be 100ms but instead got %v", rps)
	}

	s.Settings.Rps = 50
	rps = s.Settings.SetRps()
	if rps != time.Duration(20*time.Millisecond) {
		t.Errorf("expected rps to be 100ms but instead got %v", rps)
	}
}

func TestAppendAssignmentDurationsToConcatDurations(t *testing.T) {
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

	assignments := []worker.Assignment{}

	for _, q := range sqlQueries {
		a := worker.Assignment{
			SqlQuery: q,
		}
		assignments = append(assignments, a)
	}

	allAssignmentsExecutions := make(map[string][]time.Duration)
	var allDurations []time.Duration
	for _, a := range assignments {
		allAssignmentsExecutions[getAssignmentAsString(a, "sql")] = allDurations
	}

	avg_duration := 27.975609756097562 * time.Millisecond.Seconds()
	title := getAssignmentAsString(assignments[0], "sql")
	if title != "SELECT * FROM results WHERE avg_query_duration BETWEEN 13.0 AND 15.0" {
		t.Error("expected title to be SELECT * FROM results WHERE avg_query_duration BETWEEN 13.0 AND 15.0 but got", title)
	}

	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 28.495934959349594 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[1], "sql")
	if title != "SELECT * FROM results WHERE min_query_duration BETWEEN 50.0 AND 60.0" {
		t.Error("expected title to be SELECT * FROM results WHERE min_query_duration BETWEEN 50.0 AND 60.0 but got", title)
	}
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 26.528455284552845 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[2], "sql")
	if title != "SELECT * FROM results WHERE total_queries BETWEEN 100 AND 200" {
		t.Error("expected title to be SELECT * FROM results WHERE total_queries BETWEEN 100 AND 200 but got", title)
	}
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 27.13241234 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[0], "sql")
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 28.123 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[1], "sql")
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 26.12344123 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[2], "sql")
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 27.123214 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[0], "sql")
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 28.41234 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[1], "sql")
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 26.12341234 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[2], "sql")
	allAssignmentsExecutions = appendDurationsToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	val, _ := allAssignmentsExecutions["SELECT * FROM results WHERE avg_query_duration BETWEEN 13.0 AND 15.0"]
	if val[0] != time.Duration(27.975609756097562*time.Millisecond.Seconds()) {
		t.Error("expected 27.975609756097562ms got", val[0])
	}

	val, _ = allAssignmentsExecutions["SELECT * FROM results WHERE min_query_duration BETWEEN 50.0 AND 60.0"]
	if val[0] != time.Duration(28.495934959349594*time.Millisecond.Seconds()) {
		t.Error("expected 28.495934959349594ms got", val[0])
	}

	val, _ = allAssignmentsExecutions["SELECT * FROM results WHERE total_queries BETWEEN 100 AND 200"]
	if val[0] != time.Duration(26.528455284552845*time.Millisecond.Seconds()) {
		t.Error("expected 26.528455284552845ms got", val[0])
	}

	val, _ = allAssignmentsExecutions["SELECT * FROM results WHERE avg_query_duration BETWEEN 13.0 AND 15.0"]
	if val[1] != time.Duration(27.13241234*time.Millisecond.Seconds()) {
		t.Error("expected 27.13241234ms got", val[1])
	}

	val, _ = allAssignmentsExecutions["SELECT * FROM results WHERE min_query_duration BETWEEN 50.0 AND 60.0"]
	if val[1] != time.Duration(28.123*time.Millisecond.Seconds()) {
		t.Error("expected 28.123ms got", val[1])
	}

	val, _ = allAssignmentsExecutions["SELECT * FROM results WHERE total_queries BETWEEN 100 AND 200"]
	if val[1] != time.Duration(26.12344123*time.Millisecond.Seconds()) {
		t.Error("expected 26.12344123ms got", val[1])
	}
}
