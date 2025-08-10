package config

import (
	"os"
	"testing"
)

func TestLoad_FileNotFound(t *testing.T) {
	tmpdir := os.TempDir()
	path := tmpdir + "/definitely_missing_config.yaml"
	os.Remove(path)
	// t.Logf("Testing with path: %s", path)
	cfg, err := Load(path)
	if err == nil && cfg == nil {
		t.Errorf("expected error or fallback config, got nil config and nil error")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	f, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("database_path: test.kdbx\n")
	f.Close()
	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DatabasePath != "test.kdbx" {
		t.Errorf("expected database_path to be test.kdbx, got %v", cfg.DatabasePath)
	}
}

func TestCreateExampleConfig(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "exampleconfig.yaml")
	if err != nil {
		t.Fatal(err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	defer os.Remove(path)

	err = CreateExampleConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read created config: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty config file")
	}
}

func TestConfig_Print(t *testing.T) {
	cfg := &Config{
		DatabasePath:   "db.kdbx",
		DefaultOutput:  "stdout",
		PasswordFile:   "pw.txt",
		ConfigfilePath: "config.yaml",
	}
	// Just ensure Print runs without panic
	cfg.Print()
}
