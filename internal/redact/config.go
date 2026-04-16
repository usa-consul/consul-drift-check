package redact

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FileConfig represents the YAML structure for redaction rules.
type FileConfig struct {
	SensitiveKeys []string `yaml:"sensitive_keys"`
	Mask          string   `yaml:"mask"`
}

// LoadConfig reads redaction options from a YAML file at path.
// If path is empty, an empty Options struct is returned.
func LoadConfig(path string) (Options, error) {
	if path == "" {
		return Options{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Options{}, fmt.Errorf("redact: read config: %w", err)
	}

	var fc FileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return Options{}, fmt.Errorf("redact: parse config: %w", err)
	}

	return Options{
		SensitiveKeys: fc.SensitiveKeys,
		Mask:          fc.Mask,
	}, nil
}
