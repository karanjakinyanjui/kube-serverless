package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type FunctionList struct {
	Name     string `json:"name"`
	Runtime  string `json:"runtime"`
	Replicas int32  `json:"status.replicas"`
	Endpoint string `json:"status.endpoint"`
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all functions",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("%s/api/v1/functions", apiURL)
			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to list functions: %w", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			var functions []map[string]interface{}
			if err := json.Unmarshal(body, &functions); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tRUNTIME\tREPLICAS\tSTATUS")

			for _, fn := range functions {
				name := getStringValue(fn, "name")
				runtime := getStringValue(fn, "runtime")
				replicas := getInt32Value(fn, "status", "replicas")
				state := getStringValue(fn, "status", "state")

				fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", name, runtime, replicas, state)
			}

			w.Flush()
			return nil
		},
	}
}

func getStringValue(m map[string]interface{}, keys ...string) string {
	current := m
	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key].(string); ok {
				return val
			}
			return ""
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

func getInt32Value(m map[string]interface{}, keys ...string) int32 {
	current := m
	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key].(float64); ok {
				return int32(val)
			}
			return 0
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return 0
		}
	}
	return 0
}
