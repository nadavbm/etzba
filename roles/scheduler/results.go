package scheduler

import (
	"time"

	"github.com/nadavbm/etzba/roles/apiclient"
)

// Result record all task durations as duration slice and use later calculator to provide the following:
// X amount of tasks processed ,total processing time across all tasks ,the minimum task time (for a single task),
// the median task time ,the average task time ,and the maximum task time.
type Result struct {
	Assignments map[string][]time.Duration
	Responses   map[string][]*apiclient.Response
	Durations   []time.Duration
}
