package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [function-name]",
		Short: "Delete a function",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			url := fmt.Sprintf("%s/api/v1/functions/%s", apiURL, name)

			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to delete function: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				body, _ := ioutil.ReadAll(resp.Body)
				return fmt.Errorf("failed to delete function: %s", string(body))
			}

			fmt.Printf("Function '%s' deleted successfully\n", name)
			return nil
		},
	}
}
