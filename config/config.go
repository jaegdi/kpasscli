package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
    DatabasePath  string `yaml:"database_path"`
    DefaultOutput string `yaml:"default_output"`
}

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
