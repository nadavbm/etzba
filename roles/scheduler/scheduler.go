package scheduler

import (
	"strconv"
	"strings"
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
	tasksChan <-chan worker.Assignment
	// resultsChan is a channel that collects all tasks results (by time duration - how long the query or request took in milliseconds or seconds) after the worker execute the request \ query and got a response
	resultsChan     <-chan time.Duration
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
func NewScheduler(logger *zlog.Logger, duration time.Duration, configFile, helperFile string, workers int, verbose bool) (*Scheduler, error) {
	resultsCh := make(chan time.Duration)
	workCh := make(workerChannel)

	return &Scheduler{
		Logger:          logger,
		tasksChan:       workCh,
		resultsChan:     resultsCh,
		jobDuration:     duration,
		ConfigFile:      configFile,
		HelperFile:      helperFile,
		numberOfWorkers: workers,
	}, nil
}

//
// ----------------------------------------------------------------------------------------- helpers --------------------------------------------------------------------
//
// setDurationFromString get a string in a form of 30s (seconds) 12m (minutes) 1h (hours) and return the duration
func setDurationFromString(duration string) (time.Duration, error) {
	switch {
	case strings.HasSuffix(duration, "s"):
		strNum := duration[0 : len(duration)-1]
		num, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			return 1, err
		}
		return time.Duration(num) * time.Second, nil

	case strings.HasSuffix(duration, "m"):
		strNum := duration[0 : len(duration)-1]
		num, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			return 1, err
		}
		return time.Duration(num) * time.Minute, nil
	case strings.HasSuffix(duration, "h"):
		strNum := duration[0 : len(duration)-1]
		num, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			return 1, err
		}
		return time.Duration(num) * time.Hour, nil
	default:
		return time.Duration(1) * time.Second, nil
	}
}

/*
     ## #         ###
        #  		  #
	    #	 	  #
##### #####  #### ######   #####
#  #    #       # #     # #     #
####    #      #  #     # #     #
#       #     #   #     # #     #
#       #    #    #     # #     #
####    #   ####  ######   ##### #
*/
