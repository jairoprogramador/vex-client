package cmd

import (
	"errors"
	"fmt"
	"os"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/factories"
	"github.com/spf13/cobra"
	"strings"
	"github.com/jairoprogramador/fastdeploy/internal/domain/logger/vos"
)

var (
	version = "0.1.0"
	withTtyFlag bool
	colorFlag string
)

var rootCmd = &cobra.Command {
	Use:   "fd",
	Short: "fastdeploy is a CLI tool for managing and deploying projects",
	Long:  `fastdeploy is a powerful and flexible CLI tool designed to streamline your development and deployment workflows.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if cmd.HasSubCommands() && cmd.CalledAs() == "fd" {
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

		orderService, err := factory.BuildOrderService()
		if err != nil {
			return err
		}

		environmentColor := fmt.Sprintf("%s --color=%s", strings.TrimSpace(environment), colorFlag)

		logger, err := orderService.Run(cmd.Context(), command, environmentColor, withTtyFlag)
		
		if logger.Status() != vos.Success {
			presenter := factory.BuildPresenter()
			presenter.Render(logger)
		}

		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&withTtyFlag, "with-tty", false, "Enable pseudo-TTY allocation.")
	rootCmd.PersistentFlags().StringVar(&colorFlag, "color", "always", "control color output (auto, always, never)")
	rootCmd.Version = fmt.Sprintf("v%s\n", version)
	rootCmd.SetVersionTemplate(`{{.Version}}`)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}
