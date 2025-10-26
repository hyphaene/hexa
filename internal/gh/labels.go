package gh

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hyphaene/hexa/internal/config"
	"gopkg.in/yaml.v3"
)

type GithubLabel struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Color       string `json:"color" yaml:"color"`
}

// validateLabels checks that all labels have required fields
// description is optional per GitHub API (has omitempty tag)
func validateLabels(labels []GithubLabel) error {
	for i, label := range labels {
		if label.Name == "" {
			return fmt.Errorf("label[%d]: name is required", i)
		}
		if label.Color == "" {
			return fmt.Errorf("label[%d] (%s): color is required", i, label.Name)
		}
	}
	return nil
}

func ReadGithubLabelsFromConfig() ([]GithubLabel, error) {
	configPath, err := config.GetProjectConfigPath()
	if err != nil {
		return nil, err
	}

	raw, err := config.ReadYAMLField(configPath, "github.labels")
	if err != nil {
		// File or key doesn't exist - legitimate case for fresh repo
		return []GithubLabel{}, nil
	}

	// Marshal/unmarshal leverages struct tags for case-insensitive matching
	bytes, err := yaml.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("serializing labels: %w", err)
	}

	var labels []GithubLabel
	if err := yaml.Unmarshal(bytes, &labels); err != nil {
		return nil, fmt.Errorf("decoding github.labels: %w", err)
	}

	if err := validateLabels(labels); err != nil {
		return nil, err
	}

	return labels, nil
}

func FetchLabels() ([]GithubLabel, error) {
	cmd := exec.Command("gh", "label", "list", "--json", "name,description,color", "--jq", ".")

	response, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gh label list failed: %s (%w)", strings.TrimSpace(string(response)), err)
	}

	var labels []GithubLabel
	if err := json.Unmarshal(response, &labels); err != nil {
		return nil, fmt.Errorf("decoding gh response: %w", err)
	}

	if err := validateLabels(labels); err != nil {
		return nil, err
	}

	return labels, nil
}
