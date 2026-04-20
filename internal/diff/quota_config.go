package diff

import (
	"os"

	"gopkg.in/yaml.v3"
)

type quotaConfigFile struct {
	DefaultMax int            `yaml:"default_max"`
	Rules      map[string]int `yaml:"rules"`
}

// LoadQuotaOptions reads a YAML file and returns a QuotaOptions value.
// If path is empty, a zero-value QuotaOptions is returned.
func LoadQuotaOptions(path string) (QuotaOptions, error) {
	if path == "" {
		return QuotaOptions{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return QuotaOptions{}, err
	}
	var cfg quotaConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return QuotaOptions{}, err
	}
	return QuotaOptions{
		DefaultMax: cfg.DefaultMax,
		Rules:      cfg.Rules,
	}, nil
}
