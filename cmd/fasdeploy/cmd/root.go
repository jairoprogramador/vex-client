// Package cmd contains the Cobra command definitions.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jairoprogramador/fastdeploy-auth/internal/application"
	"github.com/jairoprogramador/fastdeploy-auth/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-auth/internal/infrastructure/auth"
	"github.com/jairoprogramador/fastdeploy-auth/internal/infrastructure/config"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	//commit  = "unknown"
	//date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "fd [step] [environment]",
	Short: "fd is a tool for managing your deployment workflow.",
	Long:  `A flexible command-line tool to handle authentication and execution of various deployment steps.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return errors.New("a maximum of two arguments are allowed: a step and an optional environment")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		if args[0] == "version" {
			return nil
		}

		step := args[0]
		env := ""
		if len(args) > 1 {
			env = args[1]
		}

		workDir, err := os.Getwd()
		if err != nil {
			return err
		}

		providerRepository := config.NewYAMLProviderRepository()
		authService := auth.NewAuthService()
		authAppService := application.NewAuthenticateAppService(providerRepository, authService, workDir)

		result, err := authAppService.Authenticate(cmd.Context(), step, env)
		if err != nil {
			return err
		}

		printResult(result)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = fmt.Sprintf("fd client version: %s\n", version)
	rootCmd.SetVersionTemplate(`{{.Version}}`)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}

func printResult(result *dto.Result) {
	fmt.Printf("Execution Result:\n")
	fmt.Printf("Step: %s\n", result.Step)
	fmt.Printf("Environment: %s\n", result.Env)
	fmt.Printf("WorkDir: %s\n", result.WorkDir)
	fmt.Printf("Token: %s\n", result.Token.AccessToken)
}
