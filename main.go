// main.go
package main

import (
	"fmt"
	"os"

	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/keepass"
	"kpasscli/src/output"
	"kpasscli/src/search"
)

// main is the entry point of the kpasscli application. It initializes logging,
func main() {
	// Initialize cli
	flags := Init()
	//
	// here starts the real work
	//
	debug.Log("Starting kpasscli with item: %s", flags.Item) // Debug-Log hinzugefügt

	if flags.Item == "" {
		fmt.Println("Error: Item parameter is required")
		os.Exit(1)
	}

	cfg, err := config.Load(flags.ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load config file: %v\n", err)
	} else {
		debug.Log("Loaded config file: %v", cfg)
	}

	// Resolve database path
	dbPath := keepass.ResolveDatabasePath(flags.KdbPath, cfg)
	if dbPath == "" {
		fmt.Fprintf(os.Stderr, "Error: No KeePass database path provided")
		os.Exit(1)
	} else {
		debug.Log("Resolved database path: %s", dbPath) // Debug-Log hinzugefügt
	}

	// Get database password
	kdbpasswordenv := ""
	if kpclipassparam := os.Getenv("KPASSCLI_kdbpassword"); kpclipassparam != "" {
		kdbpasswordenv = kpclipassparam
	}
	password, err := keepass.ResolvePassword(flags.KdbPassword, cfg, kdbpasswordenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting password: %v\n", err)
		os.Exit(1)
	}

	// Open database
	db, err := keepass.OpenDatabase(dbPath, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}

	if flags.ShowAll {

		err = keepass.GetAllFields(db, flags.Item)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting all fields: %v\n", err)
		}
		return

	} else {

		// Get output handler
		outputType := output.ResolveOutputType(flags.Out, cfg)
		handler := output.NewHandler(outputType)

		// Create finder with search options
		finder := search.NewFinder(db)
		finder.Options = search.SearchOptions{
			CaseSensitive: flags.CaseSensitive,
			ExactMatch:    flags.ExactMatch,
		}
		results, err := finder.Find(flags.Item)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching for item: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			fmt.Fprintf(os.Stderr, "No items found")
			os.Exit(1)
		}

		if len(results) > 1 {
			fmt.Fprintf(os.Stderr, "Multiple items found:")
			for _, result := range results {
				fmt.Fprintf(os.Stderr, "- %s\n", result.Path)
			}
			os.Exit(1)
		}

		// Get output handler
		outputType := resolveOutputType(flags.Out, cfg)
		handler := output.NewHandler(outputType)

		// Get and output field value
		value, err := results[0].GetField(flags.FieldName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting field: %v\n", err)
			os.Exit(1)
		}

		if err := handler.Output(value); err != nil {
			fmt.Fprintf(os.Stderr, "Error outputting value: %v\n", err)
			os.Exit(1)
		}
		return
	}
}
