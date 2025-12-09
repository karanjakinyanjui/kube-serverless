package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func newInvokeCommand() *cobra.Command {
	var payload string

	cmd := &cobra.Command{
		Use:   "invoke [function-name]",
		Short: "Invoke a function",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			url := fmt.Sprintf("%s/api/v1/functions/%s/invoke", apiURL, name)

			resp, err := http.Post(url, "application/json", bytes.NewBufferString(payload))
			if err != nil {
				return fmt.Errorf("failed to invoke function: %w", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to invoke function: %s", string(body))
			}

			fmt.Println(string(body))
			return nil
		},
	}

	cmd.Flags().StringVarP(&payload, "payload", "p", "{}", "Function payload (JSON)")

	return cmd
}
