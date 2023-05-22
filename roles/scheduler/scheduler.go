package scheduler

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/nadavbm/etzba/pkg/reader"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/worker"
	"github.com/nadavbm/zlog"
	"gopkg.in/yaml.v2"
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
	jobRps int64
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
func NewScheduler(logger *zlog.Logger, duration time.Duration, executionType, configFile, helperFile string, rps, workers int, verbose bool) (*Scheduler, error) {
	return &Scheduler{
		Logger:          logger,
		tasksChan:       make(workerChannel, workers),
		resultsChan:     make(chan time.Duration),
		jobDuration:     duration,
		jobRps:          int64(rps),
		ExecutionType:   executionType,
		ConfigFile:      configFile,
		HelpersFile:     helperFile,
		numberOfWorkers: workers,
		Verbose:         verbose,
	}, nil
}

//
// --------------------------------------------------------------------------------------------- worker assignments -----------------------------------------------------------------------------------------------------------------
//

// setAssignmentsToWorkers will create a slice of assignment from helpers file
func (s *Scheduler) setAssignmentsToWorkers() ([]worker.Assignment, error) {
	switch {
	case s.ExecutionType == "api":
		data, err := s.reader.ReadFile(s.HelpersFile)
		if err != nil {
			s.Logger.Fatal("could not read json file")
			return nil, err
		}

		assignments, err := s.setAPIAssignmentsToWorkers(data)
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

		return s.setSQLAssignmentsToWorkers(data), nil
	}

	return nil, errors.New("could not create assignment")
}

// setAPIAssignmentsToWorkers takes a json helpers file config and prepare worker assignments from config
func (s *Scheduler) setAPIAssignmentsToWorkers(data []byte) ([]worker.Assignment, error) {
	var requests []apiclient.ApiRequest
	switch {
	case strings.HasSuffix(s.HelpersFile, ".json"):
		if err := json.Unmarshal(data, &requests); err != nil {
			return nil, err
		}
	case strings.HasSuffix(s.HelpersFile, ".yaml"):
		if err := yaml.Unmarshal(data, &requests); err != nil {
			return nil, err
		}
	}
	var assignments []worker.Assignment
	for _, r := range requests {
		var assignment worker.Assignment
		assignment.ApiRequest = r
		assignments = append(assignments, assignment)
	}
	return assignments, nil
}

// setSQLAssignmentsToWorkers will take csv output from helpers file and create assignments for all workers
func (s *Scheduler) setSQLAssignmentsToWorkers(data [][]string) []worker.Assignment {
	var assignments []worker.Assignment
	for i, line := range data {
		if i > 0 {
			var a worker.Assignment
			for c, field := range line {
				switch {
				case c == 0:
					{
						a.SqlQuery.Command = field
					}
				case c == 1:
					{
						a.SqlQuery.Table = field
					}
				case c == 2:
					{
						a.SqlQuery.Constraint = field
					}
				case c == 3:
					{
						a.SqlQuery.ColumnsRef = field
					}
				case c == 4:
					{
						a.SqlQuery.Values = field
					}
				}
			}
			assignments = append(assignments, a)
		}
	}
	return assignments
}

//
// --------------------------------------------------------------------------------------------- worker executions -----------------------------------------------------------------------------------------------------------------
//

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
