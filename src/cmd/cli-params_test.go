package cmd

import (
	"flag"
	"testing"
)

func TestParseFlags_NoArgs(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	flags := ParseFlags(fs, []string{})
	if flags == nil {
		t.Fatal("ParseFlags returned nil")
	}
	if flags.Item != "" {
		t.Errorf("expected Item to be empty, got '%v'", flags.Item)
	}
	if flags.FieldName != "Password" {
		t.Errorf("expected FieldName to be 'Password', got '%v'", flags.FieldName)
	}
}

func TestParseFlags_LongFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{"-kdbpath", "db.kdbx", "-kdbpassword", "pw.txt", "-item", "Entry", "-fieldname", "UserName", "-out", "stdout", "-config", "cfg.yaml", "-case-sensitive", "-exact-match", "-man", "-help", "-debug", "-verify", "-create-config", "-print-config"}
	flags := ParseFlags(fs, args)
	if flags.KdbPath != "db.kdbx" {
		t.Errorf("expected KdbPath 'db.kdbx', got '%v'", flags.KdbPath)
	}
	if flags.KdbPassword != "pw.txt" {
		t.Errorf("expected KdbPassword 'pw.txt', got '%v'", flags.KdbPassword)
	}
	if flags.Item != "Entry" {
		t.Errorf("expected Item 'Entry', got '%v'", flags.Item)
	}
	if flags.FieldName != "UserName" {
		t.Errorf("expected FieldName 'UserName', got '%v'", flags.FieldName)
	}
	if flags.Out != "stdout" {
		t.Errorf("expected Out 'stdout', got '%v'", flags.Out)
	}
	if flags.ConfigPath != "cfg.yaml" {
		t.Errorf("expected ConfigPath 'cfg.yaml', got '%v'", flags.ConfigPath)
	}
	if !flags.CaseSensitive || !flags.ExactMatch || !flags.ShowMan || !flags.ShowHelp || !flags.DebugFlag || !flags.VerifyFlag || !flags.CreateConfig || !flags.PrintConfig {
		t.Error("expected all bool flags to be true")
	}
}

func TestParseFlags_ShortFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{"-p", "db.kdbx", "-w", "pw.txt", "-i", "Entry", "-f", "UserName", "-o", "stdout", "-c", "cfg.yaml", "-cs", "-e", "-m", "-h", "-d", "-v", "-cc", "-pc"}
	flags := ParseFlags(fs, args)
	if flags.KdbPath != "db.kdbx" {
		t.Errorf("expected KdbPath 'db.kdbx', got '%v'", flags.KdbPath)
	}
	if flags.KdbPassword != "pw.txt" {
		t.Errorf("expected KdbPassword 'pw.txt', got '%v'", flags.KdbPassword)
	}
	if flags.Item != "Entry" {
		t.Errorf("expected Item 'Entry', got '%v'", flags.Item)
	}
	if flags.FieldName != "UserName" {
		t.Errorf("expected FieldName 'UserName', got '%v'", flags.FieldName)
	}
	if flags.Out != "stdout" {
		t.Errorf("expected Out 'stdout', got '%v'", flags.Out)
	}
	if flags.ConfigPath != "cfg.yaml" {
		t.Errorf("expected ConfigPath 'cfg.yaml', got '%v'", flags.ConfigPath)
	}
	if !flags.CaseSensitive || !flags.ExactMatch || !flags.ShowMan || !flags.ShowHelp || !flags.DebugFlag || !flags.VerifyFlag || !flags.CreateConfig || !flags.PrintConfig {
		t.Error("expected all bool flags to be true")
	}
}
