package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fd",
	Long:  `All software has versions. This is fd's.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.SetArgs([]string{"--version"})
		rootCmd.Execute()
	},
}
