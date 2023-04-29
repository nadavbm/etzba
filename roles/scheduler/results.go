package scheduler

import "time"

// Results record all task durations as duration slice and use later calculator to provide the following:
// X amount of tasks processed ,total processing time across all tasks ,the minimum task time (for a single task),
// the median task time ,the average task time ,and the maximum task time.
type Results struct {
	ProcessingTime []time.Duration
}

// Result compose of total # processed tasks, total processing time for the job, the minimum task time,
// the median  time, the average  time, and the maximum  time.
type Result struct {
	Total                     int
	TotalProcessingTime       float64
	MinimumTime               float64
	MedianTime                float64
	AverageTime               float64
	MaximumTime               float64
	TotalJobsOfAllWorkersTime float64
}
