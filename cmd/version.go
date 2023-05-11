package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

const version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Release version of Unify",
	Long:  `Print the version number of Unify`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Unify " + version)
	},
}
