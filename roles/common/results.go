package common

import (
	"math"
	"time"

	"github.com/nadavbm/etzba/pkg/calculator"
	"github.com/nadavbm/etzba/pkg/filer"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/prompusher"
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
	Title              string        `json:"title"`
	TotalExeuctions    int           `json:"totalExecutions"`
	ProcessedDurations Durations     `json:"processedDurations"`
	ApiResponses       []ApiResponse `json:"apiResponses,omitempty"`
}

// Durations compose of total # processed tasks, total processing time for the job, the minimum task time,
// the median  time, the average  time, and the maximum  time.
type Durations struct {
	MinimumTime float64 `json:"minimumTime"`
	MedianTime  float64 `json:"medianTime"`
	AverageTime float64 `json:"averageTime"`
	MaximumTime float64 `json:"maximumTime"`
}

// ApiResponse is a similar http response during the test, that collect the count of requests with similar status code and message from api server
type ApiResponse struct {
	// Code is http status code
	Code int `json:"code"`
	// Message is the api payload
	Message string `json:"message"`
	// RequestsCount with this status code
	RequestsCount int `json:"count"`
	ContentLength int `json:"length"`
}

// PrepareResultOuput collect and process all details and durations from a test and return a Result
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
		if executionType == "api" {
			assigment.ApiResponses = processApiResponsesPerAssignment(allAssignmentsExecutionsResponses[title])
		}
		assignments = append(assignments, assigment)
	}

	pushResultToPrometheus([]string{"result"}, general.ProcessedDurations.AverageTime, general.RequestRate, float64(general.TotalExeuctions), general.JobDuration.Seconds())

	return &Result{
		Title:         title,
		StartTime:     time.Now().Local().Add(jobDuration),
		ExecutionType: executionType,
		Tags:          []string{executionType},
		General:       general,
		Assignments:   assignments,
	}
}

// ParseResult output result into json file
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

// processApiResponsesPerAssignment return ApiResponse with similar response returned from api and the request count
func processApiResponsesPerAssignment(responses []*apiclient.Response) []ApiResponse {
	Codekinds := getAllStatusCodesKindsFromApiResponse(responses)
	codeApiRespMap := make(map[int][]*apiclient.Response)
	for _, r := range responses {
		for _, v := range Codekinds {
			if r.Code == v {
				codeApiRespMap[r.Code] = append(codeApiRespMap[r.Code], r)
			}
		}
	}

	// Process only with the first api response content length and message, based on status code
	var processedApiResponses []ApiResponse
	for _, code := range Codekinds {
		apiResp := ApiResponse{
			Code:          code,
			Message:       codeApiRespMap[code][0].Status,
			RequestsCount: len(codeApiRespMap[code]),
			ContentLength: codeApiRespMap[code][0].ContentLength,
		}
		processedApiResponses = append(processedApiResponses, apiResp)
	}
	return processedApiResponses
}

// getAllStatusCodesKindsFromApiResponse
func getAllStatusCodesKindsFromApiResponse(responses []*apiclient.Response) []int {
	statusCodes := make(map[int]bool)
	for _, r := range responses {
		statusCodes[r.Code] = true
	}
	var codeTypes []int
	for k := range statusCodes {
		codeTypes = append(codeTypes, k)
	}
	return codeTypes
}

// ---------------------------------------------------------- prometheus ------------------------------------------------------------------------
//
// pushResultToPrometheus push result data to prometheus
func pushResultToPrometheus(labels []string, averageRequestDuration, requestRate, totalExecutions, jobDuration float64) {
	pushAvgReqDurationToPrometheus(labels, averageRequestDuration)
	pushRequestRateToPrometheus(labels, requestRate)
	pushTotalExecutionsToPrometheus(labels, totalExecutions)
	pushJobDurationToPrometheus(labels, jobDuration)
}

func pushAvgReqDurationToPrometheus(labels []string, averageRequestDuration float64) {
	avgReqDuration := prompusher.PrometheusClient.NewHistogram("avg_req_duration", "general average request duration", labels)
	prompusher.PrometheusClient.PushHistogram(avgReqDuration, "etzba_result", "avg_req_duration", labels, averageRequestDuration)
}

func pushRequestRateToPrometheus(labels []string, requestRate float64) {
	requestRateVec := prompusher.PrometheusClient.NewGauge("request_rate", "request rate from result", labels)
	prompusher.PrometheusClient.PushGauge(requestRateVec, "etzba_result", "request_rate", labels, requestRate)
}

func pushTotalExecutionsToPrometheus(labels []string, totalExecutions float64) {
	totalExecutionsVec := prompusher.PrometheusClient.NewGauge("total_executions", "count total executions from result", labels)
	prompusher.PrometheusClient.PushGauge(totalExecutionsVec, "etzba_result", "total_executions", labels, totalExecutions)
}

func pushJobDurationToPrometheus(labels []string, jobDuration float64) {
	totalJobDuration := prompusher.PrometheusClient.NewHistogram("job_duration", "job duration", labels)
	prompusher.PrometheusClient.PushHistogram(totalJobDuration, "etzba_result", "job_duration", labels, jobDuration)
}
