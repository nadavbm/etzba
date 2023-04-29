package cmd

import "github.com/spf13/cobra"

var (
	netCmd = &cobra.Command{
		Use:       "net",
		Short:     "Start benchmarking your network",
		Long:      `Start benchmarking network, defining the number of workers and csv file input`,
		ValidArgs: validArgs,
		Run:       benchmarkSql,
	}
)

func benchmarkNet(cmd *cobra.Command, args []string) {
	//logger := zlog.New()
	//
	//s := scheduler.NewScheduler(logger, csvFile, workers)
	//if err := s.StartBenchmarkingWithWorkers(Verbose); err != nil {
	//	s.Logger.Fatal("could not start workers")
	//}
}
