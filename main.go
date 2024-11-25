// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"kpasscli/config"
	"kpasscli/debug"
	"kpasscli/doc"
	"kpasscli/keepass"
	"kpasscli/output"
	"kpasscli/search"
)

func main() {
	// Initialize logging
	log.SetFlags(log.LstdFlags) // Entferne log.Lshortfile, um Dateinamen und Zeilennummern zu unterdrücken

	var (
		kdbPath       = flag.String("kdbpath", "", "Path to KeePass database file")
		kdbPass       = flag.String("kdbpass", "", "Password file or executable to get password")
		item          = flag.String("item", "", "Item to search for")
		fieldName     = flag.String("fieldname", "Password", "Field name to retrieve")
		out           = flag.String("out", "", "Output type (clipboard/stdout)")
		caseSensitive = flag.Bool("case-sensitive", false, "Enable case-sensitive search")
		exactMatch    = flag.Bool("exact-match", false, "Enable exact match search")
		showMan       = flag.Bool("man", false, "Show manual page")
		showHelp      = flag.Bool("help", false, "Show help message")
		debugFlag     = flag.Bool("debug", false, "Enable debug logging")
	)

	flag.Usage = doc.ShowHelp
	flag.Parse()

	if *debugFlag {
		debug.Enable()
	}

	if *showMan {
		doc.ShowMan()
		return
	}

	if *showHelp {
		doc.ShowHelp()
		return
	}

	debug.Log("Starting kpasscli with item: %s", *item) // Debug-Log hinzugefügt

	if *item == "" {
		fmt.Println("Error: Item parameter is required")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: Could not load config file: %v\n", err)
	}

	// Resolve database path
	dbPath := keepass.ResolveDatabasePath(*kdbPath, cfg)
	debug.Log("Resolved database path: %s", dbPath) // Debug-Log hinzugefügt
	if dbPath == "" {
		fmt.Println("Error: No KeePass database path provided")
		os.Exit(1)
	}

	// Get database password
	kdbpassenv := ""
	if kpclipassparam := os.Getenv("KPASSCLI_KDBPASS"); kpclipassparam != "" {
		kdbpassenv = kpclipassparam
	}
	password, err := keepass.ResolvePassword(*kdbPass, cfg, kdbpassenv)
	if err != nil {
		fmt.Printf("Error getting password: %v\n", err)
		os.Exit(1)
	}

	// Open database
	db, err := keepass.OpenDatabase(dbPath, password)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}

	// Create finder with search options
	finder := search.NewFinder(db)
	finder.Options = search.SearchOptions{
		CaseSensitive: *caseSensitive,
		ExactMatch:    *exactMatch,
	}
	results, err := finder.Find(*item)
	if err != nil {
		fmt.Printf("Error searching for item: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("No items found")
		os.Exit(1)
	}

	if len(results) > 1 {
		fmt.Println("Multiple items found:")
		for _, result := range results {
			fmt.Printf("- %s\n", result.Path)
		}
		os.Exit(1)
	}

	// Get output handler
	outputType := resolveOutputType(*out, cfg)
	handler := output.NewHandler(outputType)

	// Get and output field value
	value, err := results[0].GetField(*fieldName)
	if err != nil {
		fmt.Printf("Error getting field: %v\n", err)
		os.Exit(1)
	}

	if err := handler.Output(value); err != nil {
		fmt.Printf("Error outputting value: %v\n", err)
		os.Exit(1)
	}
}

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
