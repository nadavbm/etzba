package common

import (
	"time"

	"github.com/nadavbm/etzba/roles/apiclient"
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

// Durations compose of total # processed tasks, total processing time for the job, the minimum task time,
// the median  time, the average  time, and the maximum  time.
type Durations struct {
	Total                     int     `json:"total"`
	TotalJobTime              float64 `json:"job_time"`
	MinimumTime               float64 `json:"min_duration"`
	MedianTime                float64 `json:"med_duration"`
	AverageTime               float64 `json:"avg_duration"`
	MaximumTime               float64 `json:"max_duration"`
	TotalJobsOfAllWorkersTime float64
}
