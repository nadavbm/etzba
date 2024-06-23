package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// global vars for command args
var (
	duration     string
	rps          int
	workersCount int
	authFile     string
	configFile   string
	outputFile   string
	Verbose      bool
	validArgs    = []string{"duration", "rps", "workers", "auth", "config", "output", "verbose"}
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
	rootCmd.PersistentFlags().StringVar(&duration, "duration", "", "job duration")
	rootCmd.PersistentFlags().IntVar(&rps, "rps", 0, "how many requests per second should be executed during the job")
	rootCmd.PersistentFlags().IntVar(&workersCount, "workers", 1, "workers to run the job")
	rootCmd.PersistentFlags().StringVar(&authFile, "auth", "../../files/secret.json", "authentication file location")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "../../files/config.csv", "config file location")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "results.json", "results file output location")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(sqlCmd)
	rootCmd.AddCommand(netCmd)
	rootCmd.AddCommand(apiCmd)
}

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
