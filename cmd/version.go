package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of vAuth",
	Long:  `Show this help output, or the help for a specified subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vAuth v0.0.1")
	},
}
