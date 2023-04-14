package cmd

import (
	"github.com/spf13/cobra"
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
	//logger := zlog.New()
	//
	//s := scheduler.NewScheduler(logger, csvFile, workers)
	//if err := s.StartBenchmarkingWithWorkers(Verbose); err != nil {
	//	s.Logger.Fatal("could not start workers")
	//}
}
