package scheduler

import (
	"errors"
	"fmt"
	"time"

	"github.com/nadavbm/etzba/pkg/reader"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/etzba/roles/worker"
	"github.com/nadavbm/zlog"
)

type workerChannel chan worker.Assignment
type resultsChannel chan time.Duration

// Schduler will schedule workers to work on tasks (sql queries, net icmp or api calls) in a given duration or tasks queue
type Scheduler struct {
	Logger *zlog.Logger
	reader *reader.Reader
	//tasksChan is a channel for worker assignment. The scheduler will use this channel to scedule the amount \ frequency or weight of the worker assignment
	tasksChan chan worker.Assignment
	// resultsChan is a channel that collects all tasks results (by time duration - how long the query or request took in milliseconds or seconds) after the worker execute the request \ query and got a response
	resultsChan chan time.Duration
	// numberOfWorkers use the "--workers=x" to set the amount of workers while running a load test job
	numberOfWorkers int
	// jobDuration is how long the command should run (30s, 1m etc)
	jobDuration time.Duration
	// jobRps define the frequency for api requests or sql queries during the job execution
	jobRps int
	// tasksOrder defined by api request or sql query weight. The weight cann be defined in the config file and order tasks in the worker assignment channel by calculating the weight of each task
	tasksOrder []int
	// ExecutionType from command line arg can be sql, api or other type of executions
	ExecutionType string
	// ConfigFile used for authentication for api server or sql server
	ConfigFile string
	// HelpersFile provided via command line tool and contains all assignments for workers
	HelpersFile string
	// Verbose shows worker executions in terminal
	Verbose bool
}

// NewScheduler creates an instance of a Scheduler
func NewScheduler(logger *zlog.Logger, duration time.Duration, executionType, configFile, helperFile string, workers int, verbose bool) (*Scheduler, error) {
	return &Scheduler{
		Logger:          logger,
		tasksChan:       make(workerChannel, workers),
		resultsChan:     make(chan time.Duration),
		jobDuration:     duration,
		ExecutionType:   executionType,
		ConfigFile:      configFile,
		HelpersFile:     helperFile,
		numberOfWorkers: workers,
	}, nil
}

// setAssignmentsToWorkers will create a slice of assignment from helpers file
func (s *Scheduler) setAssignmentsToWorkers() ([]worker.Assignment, error) {
	switch {
	case s.ExecutionType == "api":
		data, err := s.reader.ReadJSONFile(s.HelpersFile)
		if err != nil {
			s.Logger.Fatal("could not read json file")
			return nil, err
		}

		assignments, err := worker.SetAPIAssignmentsToWorkers(data)
		if err != nil {
			s.Logger.Fatal("could not set api worker assignments")
			return nil, err
		}
		return assignments, nil
	case s.ExecutionType == "sql":
		data, err := s.reader.ReadCSVFile(s.HelpersFile)
		if err != nil {
			s.Logger.Fatal("could not read csv file")
			return nil, err
		}

		return worker.SetSQLAssignmentsToWorkers(data), nil
	}

	return nil, errors.New("could not create assignment")
}

// prepareAssignmentsForResultCollection set maps to prepare the result output with assignment durations and responses
func (s *Scheduler) prepareAssignmentsForResultCollection(assignments []worker.Assignment) (map[string][]time.Duration, map[string][]*apiclient.Response, error) {
	allAssignmentsExecutionsDurations := make(map[string][]time.Duration)
	allAssignmentsExecutionsResponses := make(map[string][]*apiclient.Response)
	var allAPIResponses []*apiclient.Response
	var allDurations []time.Duration
	for _, a := range assignments {
		allAssignmentsExecutionsDurations[getAssignmentAsString(a, s.ExecutionType)] = allDurations
		allAssignmentsExecutionsResponses[getAssignmentAsString(a, s.ExecutionType)] = allAPIResponses
	}

	return allAssignmentsExecutionsDurations, allAssignmentsExecutionsResponses, nil
}

// executeTaskFromAssignment will execute sql query or api request from the given worker assignment
func (s *Scheduler) executeTaskFromAssignment(assignment *worker.Assignment) (time.Duration, *apiclient.Response, error) {
	switch {
	case s.ExecutionType == "sql":
		dur, err := s.executeSQLQueryFromAssignment(assignment)
		return dur, nil, err
	case s.ExecutionType == "api":
		dur, res := s.executeAPIRequestFromAssignment(assignment)
		return dur, res, nil
	}
	return 0, nil, nil
}

// executeSQLQueryFromAssignment
func (s *Scheduler) executeSQLQueryFromAssignment(assignment *worker.Assignment) (time.Duration, error) {
	worker, err := worker.NewSQLWorker(s.Logger, s.ConfigFile)
	if err != nil {
		s.Logger.Fatal("could not create worker")
	}
	return worker.GetSQLQueryDuration(assignment)
}

// executeAPIRequestFromAssignment
func (s *Scheduler) executeAPIRequestFromAssignment(assigment *worker.Assignment) (time.Duration, *apiclient.Response) {
	worker, err := worker.NewAPIWorker(s.Logger, s.ConfigFile)
	if err != nil {
		s.Logger.Fatal("could not create worker")
	}
	return worker.GetAPIRequestDuration(assigment)

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

/*
     ## #         ###
        #         #
        #         #
##### #####  #### ######   #####
#  #    #       # #     # #     #
####    #      #  #     # #     #
#       #     #   #     # #     #
#       #    #    #     # #     #
####    #   ####  ######   ##### #
*/
