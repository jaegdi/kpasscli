package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"kpasscli/src/cmd"
	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/doc"
	"kpasscli/src/search"
)

// var Flags *cmd.Flags

func Init() *cmd.Flags {
	log.SetFlags(log.LstdFlags)
	flags := cmd.ParseFlagsDefault()
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
			debug.ErrMsg(err, "Error creating config file")
			os.Exit(1)
		}
		debug.ErrMsg(nil, fmt.Sprintf("Example config file '%s' created successfully.", configPath))
		os.Exit(0)
	}
	if flags.PrintConfig {
		cfg, err := config.Load(flags.ConfigPath)
		if err != nil {
			debug.ErrMsg(err, "Error loading config")
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
