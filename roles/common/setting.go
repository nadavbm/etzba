package common

import "time"

type Settings struct {
	// NumberOfWorkers use the "--workers=x" to set the amount of workers while running a load test job
	NumberOfWorkers int
	// Duration is how long the command should run (30s, 1m etc)
	Duration time.Duration
	// Rps define the frequency for api requests or sql queries during the job execution
	Rps int64
	// TasksOrder defined by api request or sql query weight. The weight cann be defined in the config file and order tasks in the worker assignment channel by calculating the weight of each task
	TasksOrder []int
	// ExecutionType from command line arg can be sql, api or other type of executions
	ExecutionType string
	// AuthFile used for authentication for api server or sql server
	AuthFile string
	// ConfigFile contains all assignments for workers during the etz run
	ConfigFile string
	// OutputFile is the file path to export results (json or yaml)
	OutputFile string
	// Verbose shows worker executions in terminal
	Verbose bool
}

func GetSettings(jobDuration time.Duration, exectionType, authFile, configFile, outputFile string, rps, workersCount int, verbose bool) *Settings {
	s := Settings{
		Duration:      jobDuration,
		ExecutionType: exectionType,
		AuthFile:      authFile,
		ConfigFile:    configFile,
		OutputFile:    outputFile,
		Verbose:       verbose,
	}
	if rps != 0 {
		s.Rps = int64(rps)
	}
	if workersCount != 0 {
		s.NumberOfWorkers = workersCount
	}
	return &s
}

// setRps returns a duration to sleep (in a second timeframe) that will set the amount of requests per seconds later during job execution
func (s *Settings) SetRps() time.Duration {
	var rpsSleepInDurationLoop time.Duration
	if s.Rps == 0 {
		return rpsSleepInDurationLoop
	}
	rpsSleepInDurationLoop = time.Duration(int64(1000/s.Rps) * 1000000)
	return rpsSleepInDurationLoop
}
