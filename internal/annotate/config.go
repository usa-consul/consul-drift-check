package annotate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type fileSchema struct {
	Rules []Rule `yaml:"rules"`
}

// LoadConfig reads annotation rules from a YAML file at path.
// If path is empty, an empty rule set is returned.
func LoadConfig(path string) ([]Rule, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("annotate: read config: %w", err)
	}
	var schema fileSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("annotate: parse config: %w", err)
	}
	return schema.Rules, nil
}
