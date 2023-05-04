package cmd

import (
	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/pkg/printer"
	"github.com/nadavbm/etzba/roles/scheduler"
	"github.com/nadavbm/zlog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var validArgs = []string{"file", "workers", "verbose"}

var (
	sqlCmd = &cobra.Command{
		Use:       "sql",
		Short:     "Start benchmarking your sql instance",
		Long:      `Start benchmarking db, defining the number of workers and csv file input`,
		ValidArgs: validArgs,
		Run:       benchmarkSql,
	}
)

func benchmarkSql(cmd *cobra.Command, args []string) {
	logger := zlog.New()

	jobDuration, err := setDurationFromString(duration)
	if err != nil {
		logger.Fatal("could set job duration")
	}

	s, err := scheduler.NewScheduler(logger, jobDuration, "sql", configFile, helpersFile, workersCount, Verbose)
	if err != nil {
		logger.Fatal("could not create a scheduler instance")
	}
	debug.Debug("args", duration, configFile, helpersFile, workersCount)

	var result *scheduler.Result
	if duration != "" {
		result, err = s.ExecuteTaskByDuration()
		if err != nil {
			s.Logger.Fatal("could not start execution", zap.Error(err))
		}
		debug.Debug("result", result)
	} else {
		if err := s.ExecuteJobUntilCompletion(); err != nil {
			s.Logger.Fatal("could not start execution")
		}
	}

	printer.PrintTaskDurations(result)
}
