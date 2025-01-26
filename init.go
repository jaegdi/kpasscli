package main

import (
	"flag"
	"fmt"
	"kpasscli/src/cmd"
	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/doc"
	"kpasscli/src/search"
	"log"
	"os"
	"path/filepath"
)

// var Flags *cmd.Flags

func Init() *cmd.Flags {
	log.SetFlags(log.LstdFlags)
	flags := cmd.ParseFlags()
	flag.Usage = doc.ShowHelp
	flag.Parse()

	// switch toggles
	if flags.DebugFlag {
		debug.Enable()
	}
	if flags.VerifyFlag {
		search.EnableVerify()
	}

	// Handle special flags and help messages
	if flags.CreateConfig {
		filename := "config.yaml"
		configPath := filepath.Join(".", filename)
		if err := config.CreateExampleConfig(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Example config file '%s' created successfully.", configPath)
		os.Exit(0)
	}
	if flags.PrintConfig {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		cfg.Print()
	}
	if flags.ShowMan {
		doc.ShowMan()
		os.Exit(0)
	}
	if flags.ShowHelp {
		doc.ShowHelp()
		os.Exit(0)
	}
	return flags
}
