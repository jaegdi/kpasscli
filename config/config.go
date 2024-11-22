package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration.
type Config struct {
	// DatabasePath is the default path to the KeePass database file
	DatabasePath string `yaml:"database_path"`
	// DefaultOutput specifies the default output type (stdout/clipboard)
	DefaultOutput string `yaml:"default_output"`
}

// Load reads and parses the configuration file.
// Returns:
//
//	*Config: Parsed configuration
//	error: Any error encountered while reading or parsing
func Load() (*Config, error) {
	configPath := filepath.Join("config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
