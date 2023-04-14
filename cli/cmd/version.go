package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of etzba",
	Long:  `Print the version number of the application from git tags`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 1.1.1")
	},
}
