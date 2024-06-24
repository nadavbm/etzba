package common

import (
	"math"
	"time"

	"github.com/nadavbm/etzba/pkg/calculator"
	"github.com/nadavbm/etzba/pkg/filer"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/zlog"
	"go.uber.org/zap"
)

// Result record all task durations as duration slice and use later calculator to provide the following:
// X amount of tasks processed ,total processing time across all tasks ,the minimum task time (for a single task),
// the median task time ,the average task time ,and the maximum task time.
type Result struct {
	Title         string       `json:"title"`
	StartTime     time.Time    `json:"startTime"`
	Tags          []string     `json:"tags"`
	ExecutionType string       `json:"executionType"`
	General       General      `json:"general"`
	Assignments   []Assignment `json:"assignments"`
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

// TODO: Collect api responses for output
type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func PrepareResultOuput(title, executionType string, jobDuration time.Duration, assignmentsDurations map[string][]time.Duration, allAssignmentsExecutionsResponses map[string][]*apiclient.Response) *Result {
	allDurations := concatAllDurations(assignmentsDurations)
	general := General{
		JobDuration:     jobDuration,
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
		// TODO: improve api response collection
		//var apiResponses []*apiclient.Response
		//for _, responses := range allAssignmentsExecutionsResponses {
		//	if executionType == "api" {
		//		apiResponses = append(apiResponses, responses...)
		//		assigment.ApiResponses = apiResponses
		//	}
		//}
		assignments = append(assignments, assigment)
	}

	return &Result{
		Title:         title,
		StartTime:     time.Now().Local().Add(jobDuration),
		ExecutionType: executionType,
		Tags:          []string{executionType},
		General:       general,
		Assignments:   assignments,
	}
}

func (r *Result) ParseResult(logger *zlog.Logger, outputFile string) error {
	r.General.JobDuration = time.Duration(time.Duration(r.General.JobDuration).Seconds())
	w := filer.NewWriter(logger)
	if err := w.WriteFile(outputFile, r); err != nil {
		logger.Error("could not write result to file", zap.Error(err))
		return err
	}

	return nil
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
