package validate

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds a list of validation rules loaded from a YAML file.
type Config struct {
	Rules []Rule `yaml:"rules"`
}

// LoadConfig reads a YAML config file from path and returns the parsed Config.
// An empty path returns an empty Config without error.
func LoadConfig(path string) (Config, error) {
	if path == "" {
		return Config{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
