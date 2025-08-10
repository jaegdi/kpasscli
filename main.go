package main

import (
	"fmt"
	"os"

	"github.com/tobischo/gokeepasslib/v3"

	"kpasscli/src/cmd"
	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/keepass"
	"kpasscli/src/output"
	"kpasscli/src/search"
)

// RunApp contains the main application logic and is testable.
func RunApp(
	flags *cmd.Flags,
	loadConfig func(string) (*config.Config, error),
	resolveDBPath func(string, *config.Config) string,
	resolvePassword func(string, *config.Config, string, ...keepass.PasswordPromptFunc) (string, error),
	openDatabase func(string, string) (*gokeepasslib.Database, error),
	newFinder func(*gokeepasslib.Database) search.FinderInterface,
	newHandler func(output.Type) output.Handler,
	getEnv func(string) string,
) error {
	debug.Log("Starting kpasscli with item: %s", flags.Item)

	if flags.Item == "" {
		return fmt.Errorf("item parameter is required")
	}

	config, err := loadConfig(flags.ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load config file: %v\n", err)
		if config == nil {
			return nil
		}
	}
	c
	dbPath := resolveDBPath(flags.KdbPath, config)
	debug.Log("Resolved database path: %s", dbPath)
	if dbPath == "" {
		return fmt.Errorf("no KeePass database path provided")
	}

	kdbpasswordenv := getEnv("KPASSCLI_kdbpassword")
	password, err := resolvePassword(flags.KdbPassword, config, kdbpasswordenv)
	if err != nil {
		return fmt.Errorf("Error getting password: %w", err)
	}

	db, err := openDatabase(dbPath, password)
	if err != nil {
		return fmt.Errorf("Error opening database: %w", err)
	}

	finder := newFinder(db)
	// If the Finder supports Options, set them (for real Finder)
	if f, ok := finder.(*search.Finder); ok {
		f.Options = search.SearchOptions{
			CaseSensitive: flags.CaseSensitive,
			ExactMatch:    flags.ExactMatch,
		}
	}
	results, err := finder.Find(flags.Item)
	if err != nil {
		return fmt.Errorf("Error searching for item: %w", err)
	}

	if len(results) == 0 {
		return fmt.Errorf("no items found")
	}

	if len(results) > 1 {
		for _, result := range results {
			fmt.Fprintf(os.Stderr, "- %s\n", result.Path)
			debug.Log("Found item: %s", result.Path)
		}
		return fmt.Errorf("multiple items found")
	}

	outputType := resolveOutputType(flags.Out, config)
	handler := newHandler(outputType)

	value, err := results[0].GetField(flags.FieldName)
	if err != nil {
		return fmt.Errorf("Error getting field: %w", err)
	}

	if err := handler.Output(value); err != nil {
		return fmt.Errorf("Error outputting value: %w", err)
	}

	return nil
}

func main() {
	flags := Init()
	err := RunApp(
		flags,
		config.Load,
		keepass.ResolveDatabasePath,
		func(passParam string, cfg *config.Config, kdbpassenv string, promptFunc ...keepass.PasswordPromptFunc) (string, error) {
			return keepass.ResolvePassword(passParam, cfg, kdbpassenv, promptFunc...)
		},
		keepass.OpenDatabase,
		func(db *gokeepasslib.Database) search.FinderInterface { return search.NewFinder(db) },
		output.NewHandler,
		os.Getenv,
	)
	if err != nil {
		debug.ErrMsg(err, "kpasscli")
		os.Exit(1)
	}
}

// resolveOutputType remains unchanged
func resolveOutputType(flagOut string, cfg *config.Config) output.Type {
	if flagOut != "" {
		return output.Type(flagOut)
	}
	if kpcliout := os.Getenv("KPASSCLI_OUT"); kpcliout != "" {
		if output.IsValidType(kpcliout) {
			return output.Type(kpcliout)
		}
	}
	if cfg != nil && cfg.DefaultOutput != "" {
		return output.Type(cfg.DefaultOutput)
	}
	return output.Stdout
}
