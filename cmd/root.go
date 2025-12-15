package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"hub/cmd/migrate"
	"hub/cmd/serve"
)

var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "Backend MSA Hub - Radio streaming backend service",
	Long: `Backend MSA Hub is a microservice backend for radio streaming platform.
It provides REST API endpoints and database management capabilities.`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add subcommands
	rootCmd.AddCommand(serve.NewServeCommand())
	rootCmd.AddCommand(migrate.NewMigrateCommand())
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
