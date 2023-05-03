package scheduler

import (
	"time"

	"github.com/nadavbm/etzba/roles/worker"
	"github.com/nadavbm/zlog"
)

type workerChannel chan worker.Assignment
type resultsChannel chan time.Duration

// Schduler will schedule workers to work on tasks (sql queries, net icmp or api calls) in a given duration or tasks queue
type Scheduler struct {
	Logger *zlog.Logger
	//tasksChan is a channel for worker assignment. The scheduler will use this channel to scedule the amount \ frequency or weight of the worker assignment
	tasksChan chan worker.Assignment
	// resultsChan is a channel that collects all tasks results (by time duration - how long the query or request took in milliseconds or seconds) after the worker execute the request \ query and got a response
	resultsChan     chan time.Duration
	numberOfWorkers int
	// jobDuration is how long the command should run (30s, 1m etc)
	jobDuration time.Duration
	// jobRps define the frequency for api requests or sql queries during the job execution
	jobRps int
	// tasksOrder defined by api request or sql query weight. The weight cann be defined in the config file and order tasks in the worker assignment channel by calculating the weight of each task
	tasksOrder    []int
	ExecutionType string
	ConfigFile    string
	HelperFile    string
	Verbose       bool
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
		HelperFile:      helperFile,
		numberOfWorkers: workers,
	}, nil
}

func (s *Scheduler) setAssignmentsToWorkers() ([]worker.Assignment, error) {
	var assignments []worker.Assignment
	switch {
	case s.ExecutionType == "sql":
		data, err := worker.ReadCSVFile(s.HelperFile)
		if err != nil {
			s.Logger.Fatal("could not read csv file")
			return nil, err
		}

		assignments = worker.SetSQLAssignmentsToWorkers(data)
	case s.ExecutionType == "api":
		data, err := worker.ReadJSONFile(s.HelperFile)
		if err != nil {
			s.Logger.Fatal("could not read json file")
			return nil, err
		}
		assignments, err = worker.SetAPIAssignmentsToWorkers(data)
		if err != nil {
			s.Logger.Fatal("could not set api assignments")
			return nil, err
		}
	}
	return assignments, nil
}

func (s *Scheduler) executeTaskFromAssignment(assignment *worker.Assignment) (time.Duration, error) {
	switch {
	case s.ExecutionType == "sql":
		return s.executeSQLQueriesFromAssignment(assignment)
	case s.ExecutionType == "api":
		return s.executeAPIRequestFromAssignment(assignment)
	}
	return 0, nil
}

func (s *Scheduler) executeSQLQueriesFromAssignment(assignment *worker.Assignment) (time.Duration, error) {
	worker, err := worker.NewSQLWorker(s.Logger, s.ConfigFile)
	if err != nil {
		s.Logger.Fatal("could not create worker")
	}
	return worker.GetSQLQueryDuration(assignment)
}

func (s *Scheduler) executeAPIRequestFromAssignment(assigment *worker.Assignment) (time.Duration, error) {
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
