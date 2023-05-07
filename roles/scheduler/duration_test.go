package scheduler

import (
	"testing"
	"time"

	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/etzba/roles/worker"
)

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
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 28.495934959349594 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[1], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 26.528455284552845 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[2], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 27.13241234 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[0], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 28.123 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[1], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 26.12344123 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[2], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 27.123214 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[0], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 28.41234 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[1], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

	avg_duration = 26.12341234 * time.Millisecond.Seconds()
	title = getAssignmentAsString(assignments[2], "sql")
	allAssignmentsExecutions = appendDurationToAssignmentResults(title, allAssignmentsExecutions, time.Duration(avg_duration))

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

	allDurations = concatAllDurations(allAssignmentsExecutions)
	if allDurations[0] != time.Duration(27.975609756097562*time.Millisecond.Seconds()) {
		t.Error("expected 27.975609756097562ns got", allDurations[0])
	}

}
