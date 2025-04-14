package config

import (
	"fmt"
	"kpasscli/src/debug"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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
}

// Load reads and parses the configuration file.
// Returns:
//
//	*Config: Parsed configuration
//	error: Any error encountered while reading or parsing
//
// Load loads the configuration from a YAML file. It first checks for the
// presence of "config.yaml" in the current directory. If the file does not
// exist there, it attempts to load the configuration from the user's home
// directory under ".config/kpasscli/config.yaml". The function reads the
// configuration file, unmarshals its content into a Config struct, and
// returns a pointer to the Config struct along with any error encountered
// during the process.
//
// Returns:
//   - *Config: A pointer to the Config struct containing the configuration
//     data.
//   - error: An error object if any error occurred during the loading or
//     unmarshalling of the configuration file, otherwise nil.
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
	return &config, nil
}

// CreateExampleConfig creates an example configuration file in the current directory.
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
