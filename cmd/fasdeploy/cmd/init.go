package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing 'init' command...")
		fmt.Println("Project initialized successfully.")
	},
}
