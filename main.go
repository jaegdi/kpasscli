// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/doc"
	"kpasscli/src/keepass"
	"kpasscli/src/output"
	"kpasscli/src/search"
)

type Flags struct {
	KdbPath       string
	KdbPass       string
	Item          string
	FieldName     string
	Out           string
	CaseSensitive bool
	ExactMatch    bool
	ShowMan       bool
	ShowHelp      bool
	DebugFlag     bool
	VerifyFlag    bool
	CreateConfig  bool
}

func parseFlags() *Flags {
	flags := &Flags{}

	// Define flags with both long and short versions
	flag.StringVar(&flags.KdbPath, "kdbpath", "", "Path to KeePass database file")
	flag.StringVar(&flags.KdbPath, "p", "", "Path to KeePass database file (shorthand)")

	flag.StringVar(&flags.KdbPass, "kdbpass", "", "Password file or executable to get password")
	flag.StringVar(&flags.KdbPass, "w", "", "Password file or executable to get password (shorthand)")

	flag.StringVar(&flags.Item, "item", "", "Item to search for")
	flag.StringVar(&flags.Item, "i", "", "Item to search for (shorthand)")

	flag.StringVar(&flags.FieldName, "fieldname", "Password", "Field name to retrieve")
	flag.StringVar(&flags.FieldName, "f", "Password", "Field name to retrieve (shorthand)")

	flag.StringVar(&flags.Out, "out", "", "Output type (clipboard/stdout)")
	flag.StringVar(&flags.Out, "o", "", "Output type (clipboard/stdout) (shorthand)")

	flag.BoolVar(&flags.CaseSensitive, "case-sensitive", false, "Enable case-sensitive search")
	flag.BoolVar(&flags.CaseSensitive, "c", false, "Enable case-sensitive search (shorthand)")

	flag.BoolVar(&flags.ExactMatch, "exact-match", false, "Enable exact match search")
	flag.BoolVar(&flags.ExactMatch, "e", false, "Enable exact match search (shorthand)")

	flag.BoolVar(&flags.ShowMan, "man", false, "Show manual page")
	flag.BoolVar(&flags.ShowMan, "m", false, "Show manual page (shorthand)")

	flag.BoolVar(&flags.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&flags.ShowHelp, "h", false, "Show help message (shorthand)")

	flag.BoolVar(&flags.VerifyFlag, "verify", false, "Enable verify message")
	flag.BoolVar(&flags.VerifyFlag, "v", false, "Enable verify message (shorthand)")

	flag.BoolVar(&flags.DebugFlag, "debug", false, "Enable debug logging")
	flag.BoolVar(&flags.DebugFlag, "d", false, "Enable debug logging (shorthand)")

	flag.BoolVar(&flags.CreateConfig, "create-config", false, "Create an example config file")
	flag.BoolVar(&flags.CreateConfig, "cc", false, "Create an example config file (shorthand)")

	flag.Usage = doc.ShowHelp
	flag.Parse()

	return flags
}

// main is the entry point of the kpasscli application. It initializes logging,
// parses command-line flags, and performs various actions based on the provided
// flags. The main function supports the following flags:
//
// -kdbpath|-p: Path to KeePass database file
// -kdbpass|-w: Password file or executable to get password
// -item|-i: Item to search for
// -fieldname|-f: Field name to retrieve (default: "Password")
// -out|-o: Output type (clipboard/stdout)
// -case-sensitive|-c: Enable case-sensitive search
// -exact-match|-e: Enable exact match search
// -man|-m: Show manual page
// -help|-h: Show help message
// -debug|-d: Enable debug logging
// -create-config|-cc: Create an example config file
//
// The function performs the following steps:
// 1. Initializes logging based on the debug flag.
// 2. Creates an example config file if the create-config flag is set.
// 3. Displays the manual page if the man flag is set.
// 4. Displays the help message if the help flag is set.
// 5. Loads the configuration file.
// 6. Resolves the KeePass database path.
// 7. Retrieves the database password.
// 8. Opens the KeePass database.
// 9. Searches for the specified item in the database.
// 10. Outputs the value of the specified field using the specified output handler.
func main() {
	// Initialize logging
	log.SetFlags(log.LstdFlags)
	flags := parseFlags()
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
		if err := config.CreateExampleConfig(); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Example config file created successfully.")
		return
	}
	if flags.ShowMan {
		doc.ShowMan()
		return
	}
	if flags.ShowHelp {
		doc.ShowHelp()
		return
	}

	//
	// here starts the real work
	//
	debug.Log("Starting kpasscli with item: %s", flags.Item) // Debug-Log hinzugefügt

	if flags.Item == "" {
		fmt.Println("Error: Item parameter is required")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: Could not load config file: %v\n", err)
	}

	// Resolve database path
	dbPath := keepass.ResolveDatabasePath(flags.KdbPath, cfg)
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
	password, err := keepass.ResolvePassword(flags.KdbPass, cfg, kdbpassenv)
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
		CaseSensitive: flags.CaseSensitive,
		ExactMatch:    flags.ExactMatch,
	}
	results, err := finder.Find(flags.Item)
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
	outputType := resolveOutputType(flags.Out, cfg)
	handler := output.NewHandler(outputType)

	// Get and output field value
	value, err := results[0].GetField(flags.FieldName)
	if err != nil {
		fmt.Printf("Error getting field: %v\n", err)
		os.Exit(1)
	}

	if err := handler.Output(value); err != nil {
		fmt.Printf("Error outputting value: %v\n", err)
		os.Exit(1)
	}
}

// resolveOutputType determines the output type based on the provided flag,
// environment variable, or configuration. It follows this order of precedence:
// 1. If the flagOut parameter is not empty, it returns the corresponding output type.
// 2. If the environment variable "KPASSCLI_OUT" is set and valid, it returns the corresponding output type.
// 3. If the cfg parameter is not nil and cfg.DefaultOutput is not empty, it returns the corresponding output type.
// 4. If none of the above conditions are met, it defaults to output.Stdout.
//
// Parameters:
// - flagOut: A string representing the output type specified by a flag.
// - cfg: A pointer to a config.Config struct that may contain a default output type.
//
// Returns:
// - output.Type: The resolved output type based on the provided inputs.
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
