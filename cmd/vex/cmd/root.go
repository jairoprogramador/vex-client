package cmd

import (
	"errors"
	"fmt"
	"github.com/jairoprogramador/vex-client/internal/infrastructure/factories"
	"github.com/spf13/cobra"
	"os"
)

var (
	withTtyFlag bool
	colorFlag   string
	version     string
)

var rootCmd = &cobra.Command{
	Use:   "vex",
	Short: "vex is a CLI tool for managing and deploying projects",
	Long:  `vex is a powerful and flexible CLI tool designed to streamline your development and deployment workflows.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if cmd.HasSubCommands() && cmd.CalledAs() == "vex" {
				return nil
			}
			return errors.New("a step argument is required")
		}
		if len(args) > 2 {
			return errors.New("a maximum of two arguments are allowed: a step and an optional environment")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if len(args) == 0 {
			return cmd.Help()
		}

		command := args[0]
		environment := ""
		if len(args) > 1 {
			environment = args[1]
		}

		factory := factories.NewServiceFactory()

		executorService, err := factory.BuildExecutor()
		if err != nil {
			return err
		}

		return executorService.Run(cmd.Context(), command, environment)
	},
}

func Execute(versionMain string) {
	version = versionMain
	rootCmd.Version = fmt.Sprintf("v%s\n", version)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&withTtyFlag, "with-tty", false, "Enable pseudo-TTY allocation.")
	rootCmd.PersistentFlags().StringVar(&colorFlag, "color", "always", "control color output (auto, always, never)")
	rootCmd.SetVersionTemplate(`{{.Version}}`)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}
