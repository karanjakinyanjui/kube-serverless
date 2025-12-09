package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiURL string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ksls",
		Short: "Kube-Serverless CLI - Manage serverless functions on Kubernetes",
		Long: `A command-line interface for deploying and managing serverless functions
on the Kube-Serverless platform.`,
	}

	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "http://localhost:8080", "API server URL")

	// Add commands
	rootCmd.AddCommand(newDeployCommand())
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newGetCommand())
	rootCmd.AddCommand(newDeleteCommand())
	rootCmd.AddCommand(newInvokeCommand())
	rootCmd.AddCommand(newLogsCommand())
	rootCmd.AddCommand(newMetricsCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
