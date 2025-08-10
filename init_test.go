package main

import (
	"os"
	"testing"
)

func TestInit_ReturnsFlags(t *testing.T) {
	// Save and restore original os.Args
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"kpasscli", "-item", "TestItem", "-kdbpath", "test.kdbx", "-out", "stdout", "-print-config", "-debug", "-verify"}
	flags := Init()

	if flags == nil {
		t.Fatal("Init returned nil flags")
	}
	if flags.Item != "TestItem" {
		t.Errorf("expected Item to be 'TestItem', got '%v'", flags.Item)
	}
	if flags.KdbPath != "test.kdbx" {
		t.Errorf("expected KdbPath to be 'test.kdbx', got '%v'", flags.KdbPath)
	}
	if flags.Out != "stdout" {
		t.Errorf("expected Out to be 'stdout', got '%v'", flags.Out)
	}
}
