package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// global vars for command args
var (
	workers  int
	csvFile  string
	reqFile  string
	authFile string
	seedFile string
	Verbose  bool
)

var (
	rootCmd = &cobra.Command{
		Use:   "etz",
		Short: "etz root command",
		Long: `Root command for etzba.
				  Complete documentation is available at https://github.com/nadavbm/etzba`,
		Run: func(cmd *cobra.Command, args []string) {
			// Empty
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	sqlCmd.PersistentFlags().StringVar(&csvFile, "file", "../../files/somefile.csv", "read from csv file")
	sqlCmd.PersistentFlags().IntVar(&workers, "workers", 1, "define how many workers should run the job")

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(sqlCmd)
}
