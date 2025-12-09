package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLogsCommand() *cobra.Command {
	var follow bool

	cmd := &cobra.Command{
		Use:   "logs [function-name]",
		Short: "Get function logs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			fmt.Printf("To view logs for function '%s', use:\n", name)
			fmt.Printf("kubectl logs -l function=%s -n kube-serverless", name)
			if follow {
				fmt.Printf(" -f")
			}
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")

	return cmd
}
