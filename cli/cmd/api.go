package cmd

import (
	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/roles/scheduler"
	"github.com/nadavbm/zlog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	apiCmd = &cobra.Command{
		Use:       "api",
		Short:     "Start benchmarking api server",
		Long:      `Start benchmarking api, defining the number of workers`,
		ValidArgs: validArgs,
		Run:       benchmarkAPI,
	}
)

func benchmarkAPI(cmd *cobra.Command, args []string) {
	logger := zlog.New()

	jobDuration, err := setDurationFromString(duration)
	if err != nil {
		logger.Fatal("could set job duration")
	}

	s, err := scheduler.NewScheduler(logger, jobDuration, "api", configFile, helpersFile, workersCount, Verbose)
	if err != nil {
		logger.Fatal("could not create a scheduler instance")
	}
	debug.Debug("args", duration, configFile, helpersFile, workersCount)

	if duration != "" {
		result, err := s.ExecuteTaskByDuration()
		if err != nil {
			s.Logger.Fatal("could not start execution", zap.Error(err))
		}
		debug.Debug("result", result)
	} else {
		if err := s.ExecuteJobUntilCompletion(); err != nil {
			s.Logger.Fatal("could not start execution")
		}
	}
}
