package scheduler

import (
	"fmt"
	"time"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/etzba/roles/worker"
)

// Result record all task durations as duration slice and use later calculator to provide the following:
// X amount of tasks processed ,total processing time across all tasks ,the minimum task time (for a single task),
// the median task time ,the average task time ,and the maximum task time.
type Result struct {
	JobDuration time.Duration
	RequestRate int
	Assignments map[string][]time.Duration
	Responses   map[string][]*apiclient.Response
	Durations   []time.Duration
}

// prepareAssignmentsForResultCollection set maps to prepare the result output with assignment durations and responses
func (s *Scheduler) prepareAssignmentsForResultCollection(assignments []worker.Assignment) (map[string][]time.Duration, map[string][]*apiclient.Response, error) {
	allAssignmentsExecutionsDurations := make(map[string][]time.Duration)
	allAssignmentsExecutionsResponses := make(map[string][]*apiclient.Response)
	var allAPIResponses []*apiclient.Response
	var allDurations []time.Duration
	for _, a := range assignments {
		allAssignmentsExecutionsDurations[getAssignmentAsString(a, s.Settings.ExecutionType)] = allDurations
		allAssignmentsExecutionsResponses[getAssignmentAsString(a, s.Settings.ExecutionType)] = allAPIResponses
	}

	return allAssignmentsExecutionsDurations, allAssignmentsExecutionsResponses, nil
}

//
// ------------------------------------------------------------------------------- helpers ------------------------------------------------------------------
//

// getAssignmentAsString prepare the assignment title for stdout
func getAssignmentAsString(a worker.Assignment, command string) string {
	switch {
	case command == "api":
		return fmt.Sprintf("URL: %s, Method: %s", a.ApiRequest.Url, a.ApiRequest.Method)
	case command == "sql":
		return fmt.Sprintf("%s", sqlclient.ToSQL(&a.SqlQuery))
	}
	return ""
}
