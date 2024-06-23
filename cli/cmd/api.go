package cmd

import (
	"github.com/nadavbm/etzba/pkg/filer"
	"github.com/nadavbm/etzba/pkg/printer"
	"github.com/nadavbm/etzba/roles/common"
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
		logger.Fatal("could set job duration", zap.Error(err))
	}

	settings := common.GetSettings(jobDuration, "api", authFile, configFile, outputFile, rps, workersCount, Verbose)
	s, err := scheduler.NewScheduler(logger, settings)
	if err != nil {
		logger.Fatal("could not create a scheduler instance", zap.Error(err))
	}

	var result *common.Result
	if duration != "" {
		result, err = s.ExecuteJobByDuration()
		if err != nil {
			s.Logger.Fatal("could not execute job by duration", zap.Error(err))
		}
	} else {
		if result, err = s.ExecuteJobUntilCompletion(); err != nil {
			s.Logger.Fatal("could not execute job until completion", zap.Error(err))
		}
	}

	if outputFile != "" {
		w := filer.NewWriter(logger)
		if err := w.WriteFile(outputFile, result); err != nil {
			logger.Error("could not write result to file", zap.Any("result", result), zap.Error(err))
		}
	}

	printer.PrintToTerminal(result, true)
}
