package common

import (
	"math"
	"time"

	"github.com/nadavbm/etzba/pkg/calculator"
	"github.com/nadavbm/etzba/roles/apiclient"
)

// Result record all task durations as duration slice and use later calculator to provide the following:
// X amount of tasks processed ,total processing time across all tasks ,the minimum task time (for a single task),
// the median task time ,the average task time ,and the maximum task time.
type Result struct {
	General     General      `json:"general"`
	Assignments []Assignment `json:"assignments"`
}

type General struct {
	JobDuration        time.Duration `json:"jobDuration"`
	TotalExeuctions    int           `json:"totalExecutions"`
	RequestRate        float64       `json:"requestRate"`
	ProcessedDurations Durations     `json:"processedDurations"`
}

type Assignment struct {
	Title              string                `json:"title"`
	TotalExeuctions    int                   `json:"totalExecutions"`
	ProcessedDurations Durations             `json:"processedDurations"`
	ApiResponses       []*apiclient.Response `json:"apiResponses"`
}

// Durations compose of total # processed tasks, total processing time for the job, the minimum task time,
// the median  time, the average  time, and the maximum  time.
type Durations struct {
	MinimumTime float64 `json:"minimumTime"`
	MedianTime  float64 `json:"medianTime"`
	AverageTime float64 `json:"averageTime"`
	MaximumTime float64 `json:"maximumTime"`
}

func PrepareResultOuput(jobDuration time.Duration, assignmentsDurations map[string][]time.Duration, allAssignmentsExecutionsResponses map[string][]*apiclient.Response) *Result {
	allDurations := concatAllDurations(assignmentsDurations)
	general := General{
		JobDuration:     time.Duration(jobDuration.Seconds()),
		TotalExeuctions: len(allDurations),
		RequestRate:     calculateRequestRate(jobDuration, len(allDurations)),
		ProcessedDurations: Durations{
			MinimumTime: calculator.GetMinimumTime(allDurations),
			MedianTime:  calculator.GetMedianTime(allDurations),
			AverageTime: calculator.GetAverageTime(allDurations),
			MaximumTime: calculator.GetMaximumTime(allDurations),
		},
	}

	var assignments []Assignment
	for title, durations := range assignmentsDurations {
		assigment := Assignment{
			Title:           title,
			TotalExeuctions: len(durations),
			ProcessedDurations: Durations{
				MinimumTime: calculator.GetMinimumTime(durations),
				MedianTime:  calculator.GetMedianTime(durations),
				AverageTime: calculator.GetAverageTime(durations),
				MaximumTime: calculator.GetMaximumTime(durations),
			},
		}
		var apiResponses []*apiclient.Response
		for t, responses := range allAssignmentsExecutionsResponses {
			if t == title && responses[0] != nil {
				apiResponses = append(apiResponses, responses...)
				assigment.ApiResponses = apiResponses
			}
		}
		assignments = append(assignments, assigment)
	}

	return &Result{
		General:     general,
		Assignments: assignments,
	}
}

// concatAllDurations from assignment results to return durations from all assignments
func concatAllDurations(assignmentsDurations map[string][]time.Duration) []time.Duration {
	var allDurations []time.Duration
	for _, val := range assignmentsDurations {
		allDurations = append(allDurations, val...)
	}
	return allDurations
}

// calculateRequestRate return the request per second value
func calculateRequestRate(jobDuration time.Duration, totalExecutions int) float64 {
	return math.Round(float64(totalExecutions*1000000000) / (float64(jobDuration)))
}
