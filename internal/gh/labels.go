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

func ReadGithubLabelsFromConfig() ([]GithubLabel, error) {

	configPath, err := config.GetProjectConfigPath()
	if err != nil {
		return nil, err
	}

	raw, err := config.ReadYAMLField(configPath, "github.labels")
	if err != nil {
		return nil, err
	}

	// Re-marshal pour s’appuyer sur yaml.Unmarshal et profiter des tags struct
	bytes, err := yaml.Marshal(normalizeKeys(raw))
	if err != nil {
		return nil, fmt.Errorf("serializing labels: %w", err)
	}

	var labels []GithubLabel
	if err := yaml.Unmarshal(bytes, &labels); err != nil {
		return nil, fmt.Errorf("decoding github.labels: %w", err)
	}
	return labels, nil
}

func normalizeKeys(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalized := make(map[string]any, len(typed))
		for key, val := range typed {
			normalized[strings.ToLower(key)] = normalizeKeys(val)
		}
		return normalized
	case map[any]any:
		normalized := make(map[string]any, len(typed))
		for key, val := range typed {
			strKey, ok := key.(string)
			if !ok {
				continue
			}
			normalized[strings.ToLower(strKey)] = normalizeKeys(val)
		}
		return normalized
	case []any:
		slice := make([]any, len(typed))
		for i, val := range typed {
			slice[i] = normalizeKeys(val)
		}
		return slice
	default:
		return typed
	}
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

	return labels, nil
}
