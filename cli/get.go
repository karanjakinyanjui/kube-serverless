package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get [function-name]",
		Short: "Get function details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			url := fmt.Sprintf("%s/api/v1/functions/%s", apiURL, name)

			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to get function: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				return fmt.Errorf("failed to get function: %s", string(body))
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			var function map[string]interface{}
			if err := json.Unmarshal(body, &function); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			// Pretty print JSON
			prettyJSON, err := json.MarshalIndent(function, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format response: %w", err)
			}

			fmt.Println(string(prettyJSON))
			return nil
		},
	}
}
