package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nadavbm/etzba/pkg/filer"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/etzba/roles/common"
	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/etzba/roles/worker"
	"github.com/nadavbm/zlog"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var wg sync.WaitGroup
var mutex = &sync.Mutex{}

type workerChannel chan worker.Assignment

// Schduler will schedule workers to work on tasks (sql queries, net icmp or api calls) in a given duration or tasks queue
type Scheduler struct {
	Logger *zlog.Logger
	reader *filer.Reader
	//tasksChan is a channel for worker assignment. The scheduler will use this channel to scedule the amount \ frequency or weight of the worker assignment
	tasksChan chan worker.Assignment
	// resultsChan is a channel that collects all tasks results (by time duration - how long the query or request took in milliseconds or seconds) after the worker execute the request \ query and got a response
	resultsChan chan time.Duration
	// Settings
	Settings *common.Settings
	// authenticator
	Authenticator *authenticator.Authenticator
	// connectionPool
	ConnectionPool *pgxpool.Pool
}

// NewScheduler creates an instance of a Scheduler
func NewScheduler(logger *zlog.Logger, settings *common.Settings) (*Scheduler, error) {
	return &Scheduler{
		Logger:        logger,
		tasksChan:     make(workerChannel, settings.NumberOfWorkers),
		resultsChan:   make(chan time.Duration),
		Settings:      settings,
		Authenticator: authenticator.NewAuthenticator(logger, settings.AuthFile),
	}, nil
}

//
// --------------------------------------------------------------------------------------------- worker assignments -----------------------------------------------------------------------------------------------------------------
//

// setAssignmentsToWorkers will create a slice of assignment from config file
func (s *Scheduler) setAssignmentsToWorkers() ([]worker.Assignment, error) {
	switch {
	case s.Settings.ExecutionType == "api":
		data, err := s.reader.ReadFile(s.Settings.ConfigFile)
		if err != nil {
			s.Logger.Fatal("could not read config file", zap.Error(err))
			return nil, err
		}

		assignments, err := s.setAPIAssignmentsToWorkers(data)
		if err != nil {
			s.Logger.Fatal("could not set api worker assignments", zap.Error(err))
			return nil, err
		}
		return assignments, nil
	case s.Settings.ExecutionType == "sql":
		data, err := s.reader.ReadCSVFile(s.Settings.ConfigFile)
		if err != nil {
			s.Logger.Fatal("could not read config csv file", zap.Error(err))
			return nil, err
		}

		return s.setSQLAssignmentsToWorkers(data), nil
	}

	return nil, errors.New("could not create assignment")
}

// setAPIAssignmentsToWorkers takes a json file config and prepare worker assignments from config
func (s *Scheduler) setAPIAssignmentsToWorkers(data []byte) ([]worker.Assignment, error) {
	var requests []apiclient.ApiRequest
	switch {
	case strings.HasSuffix(s.Settings.ConfigFile, ".json"):
		if err := json.Unmarshal(data, &requests); err != nil {
			return nil, err
		}
	case strings.HasSuffix(s.Settings.ConfigFile, ".yaml"):
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

// setSQLAssignmentsToWorkers will take csv output from config file and create assignments for all workers
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
	case s.Settings.ExecutionType == "sql":
		dur, err := s.executeSQLQueryFromAssignment(assignment)
		return dur, nil, err
	case s.Settings.ExecutionType == "api":
		dur, res := s.executeAPIRequestFromAssignment(assignment)
		return dur, res, nil
	}
	return 0, nil, nil
}

// executeSQLQueryFromAssignment
func (s *Scheduler) executeSQLQueryFromAssignment(assignment *worker.Assignment) (time.Duration, error) {
	worker, err := worker.NewSQLWorker(s.Logger, s.Settings.ConfigFile, s.ConnectionPool)
	if err != nil {
		s.Logger.Fatal("could not create new sql worker", zap.Error(err))
	}
	return worker.GetSQLQueryDuration(assignment)
}

// executeAPIRequestFromAssignment
func (s *Scheduler) executeAPIRequestFromAssignment(assigment *worker.Assignment) (time.Duration, *apiclient.Response) {
	worker, err := worker.NewAPIWorker(s.Logger, s.Settings.AuthFile)
	if err != nil {
		s.Logger.Fatal("could not create new api worker", zap.Error(err))
	}
	return worker.GetAPIRequestDuration(assigment)

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
		return sqlclient.ToSQL(&a.SqlQuery)
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
