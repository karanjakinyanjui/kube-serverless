package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type FunctionSpec struct {
	Name        string            `yaml:"name" json:"name"`
	Runtime     string            `yaml:"runtime" json:"runtime"`
	Handler     string            `yaml:"handler" json:"handler"`
	Code        string            `yaml:"code" json:"code"`
	CodeFile    string            `yaml:"codeFile,omitempty" json:"-"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	MinReplicas int32             `yaml:"minReplicas,omitempty" json:"minReplicas,omitempty"`
	MaxReplicas int32             `yaml:"maxReplicas,omitempty" json:"maxReplicas,omitempty"`
	Triggers    []Trigger         `yaml:"triggers,omitempty" json:"triggers,omitempty"`
}

type Trigger struct {
	Type   string            `yaml:"type" json:"type"`
	Config map[string]string `yaml:"config" json:"config"`
}

func newDeployCommand() *cobra.Command {
	var (
		functionFile string
		runtime      string
		handler      string
		codeFile     string
		minReplicas  int32
		maxReplicas  int32
	)

	cmd := &cobra.Command{
		Use:   "deploy [function-name]",
		Short: "Deploy a function",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var spec FunctionSpec

			if functionFile != "" {
				// Load from file
				data, err := ioutil.ReadFile(functionFile)
				if err != nil {
					return fmt.Errorf("failed to read function file: %w", err)
				}

				if err := yaml.Unmarshal(data, &spec); err != nil {
					return fmt.Errorf("failed to parse function file: %w", err)
				}
			} else {
				// Build from flags
				if len(args) == 0 {
					return fmt.Errorf("function name is required")
				}

				spec.Name = args[0]
				spec.Runtime = runtime
				spec.Handler = handler
				spec.MinReplicas = minReplicas
				spec.MaxReplicas = maxReplicas

				if codeFile != "" {
					code, err := ioutil.ReadFile(codeFile)
					if err != nil {
						return fmt.Errorf("failed to read code file: %w", err)
					}
					spec.Code = string(code)
				}
			}

			// If codeFile is specified in spec, load it
			if spec.CodeFile != "" {
				code, err := ioutil.ReadFile(spec.CodeFile)
				if err != nil {
					return fmt.Errorf("failed to read code file: %w", err)
				}
				spec.Code = string(code)
			}

			return deployFunction(&spec)
		},
	}

	cmd.Flags().StringVarP(&functionFile, "file", "f", "", "Function specification file (YAML)")
	cmd.Flags().StringVarP(&runtime, "runtime", "r", "nodejs18", "Runtime (nodejs18, python39, go119)")
	cmd.Flags().StringVar(&handler, "handler", "index.handler", "Function handler")
	cmd.Flags().StringVarP(&codeFile, "code", "c", "", "Code file path")
	cmd.Flags().Int32Var(&minReplicas, "min-replicas", 0, "Minimum replicas")
	cmd.Flags().Int32Var(&maxReplicas, "max-replicas", 10, "Maximum replicas")

	return cmd
}

func deployFunction(spec *FunctionSpec) error {
	jsonData, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal function spec: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/functions", apiURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to deploy function: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to deploy function: %s", string(body))
	}

	fmt.Printf("Function '%s' deployed successfully\n", spec.Name)
	return nil
}
