package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"kpasscli/src/debug"
)

// Config represents the application configuration.
type Config struct {
	// DatabasePath is the default path to the KeePass database file
	DatabasePath string `yaml:"database_path"`
	// DefaultOutput specifies the default output type (stdout/clipboard)
	DefaultOutput      string `yaml:"default_output"`
	PasswordFile       string `yaml:"password_file"`
	PasswordExecutable string `yaml:"password_executable"`
	ConfigfilePath     string `yaml:"configfile_path"`
	OutputFormat       string `yaml:"output_format"`
}

// Load reads and parses the configuration file from the given path.
//
// Parameters:
//   - configPath: The path to the configuration file to load.
//
// Returns:
//   - *Config: The loaded configuration struct.
//   - error: Any error encountered during loading or parsing.
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = filepath.Join(".", "config.yaml")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(os.Getenv("HOME"), ".config", "kpasscli", "config.yaml")
	}
	debug.Log("Loading config from: %s\n", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	debug.Log("Loaded data: %v\n", string(data))
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	// config.ConfigfilePath = configPath
	y, _ := yaml.Marshal(config)
	debug.Log("Loaded config: %v\n", string(y))
	config.OutputFormat = "text"
	return &config, nil
}

// CreateExampleConfig creates an example configuration file at the specified path.
//
// Parameters:
//   - configPath: The path where the example config file should be created.
//
// Returns:
//   - error: Any error encountered during the creation of the config file.
func CreateExampleConfig(configPath string) error {
	exampleConfig := Config{
		DatabasePath:       "/path/to/your/database.kdbx",
		DefaultOutput:      "stdout",
		PasswordFile:       "/path/to/your/password.txt",
		PasswordExecutable: "[/path/to/your/]password_executable.sh",
	}
	data, err := yaml.Marshal(&exampleConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func (c *Config) Print() {
	fmt.Fprintf(os.Stderr, "Current used Configuration: %s\n", c.ConfigfilePath)
	fmt.Fprintf(os.Stderr, "------------------------------------------\n")
	fmt.Fprintf(os.Stderr, "Database Path: %s\n", c.DatabasePath)
	fmt.Fprintf(os.Stderr, "Default Output: %s\n", c.DefaultOutput)
	fmt.Fprintf(os.Stderr, "Password File: %s\n", c.PasswordFile)
	fmt.Fprintf(os.Stderr, "Password Executable: %s\n", c.PasswordExecutable)
	fmt.Fprintf(os.Stderr, "------------------------------------------\n")
}
