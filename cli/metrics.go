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

type FunctionMetrics struct {
	Invocations  int64   `json:"invocations"`
	ColdStarts   int64   `json:"coldStarts"`
	AvgDuration  float64 `json:"avgDuration"`
	ErrorRate    float64 `json:"errorRate"`
	CostEstimate float64 `json:"costEstimate"`
}

func newMetricsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "metrics [function-name]",
		Short: "Get function metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			url := fmt.Sprintf("%s/api/v1/functions/%s/metrics", apiURL, name)

			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to get metrics: %w", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to get metrics: %s", string(body))
			}

			var metrics FunctionMetrics
			if err := json.Unmarshal(body, &metrics); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "METRIC\tVALUE")
			fmt.Fprintf(w, "Invocations\t%d\n", metrics.Invocations)
			fmt.Fprintf(w, "Cold Starts\t%d\n", metrics.ColdStarts)
			fmt.Fprintf(w, "Avg Duration\t%.3fs\n", metrics.AvgDuration)
			fmt.Fprintf(w, "Error Rate\t%.2f%%\n", metrics.ErrorRate*100)
			fmt.Fprintf(w, "Cost Estimate\t$%.4f\n", metrics.CostEstimate)
			w.Flush()

			return nil
		},
	}
}
