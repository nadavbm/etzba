package cmd

import "github.com/spf13/cobra"

var (
	apiCmd = &cobra.Command{
		Use:       "api",
		Short:     "Start benchmarking api server",
		Long:      `Start benchmarking api, defining the number of workers`,
		ValidArgs: validArgs,
		Run:       benchmarkSql,
	}
)

func benchmarkAPI(cmd *cobra.Command, args []string) {
	//logger := zlog.New()
	//
	//s := scheduler.NewScheduler(logger, csvFile, workers)
	//if err := s.StartBenchmarkingWithWorkers(Verbose); err != nil {
	//	s.Logger.Fatal("could not start workers")
	//}
}
