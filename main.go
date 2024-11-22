package main

import (
	"flag"
	"fmt"
	"os"

	"kpasscli/config"
	"kpasscli/keepass"
	"kpasscli/output"
	"kpasscli/search"
)

func main() {
    var (
        kdbPath   = flag.String("kdbpath", "", "Path to KeePass database file")
        kdbPass   = flag.String("kdbpass", "", "Password file or executable to get password")
        item      = flag.String("item", "", "Item to search for")
        fieldName = flag.String("fieldname", "Password", "Field name to retrieve")
        out       = flag.String("out", "", "Output type (clipboard/stdout)")
    )
    flag.Parse()

    if *item == "" {
        fmt.Println("Error: Item parameter is required")
        os.Exit(1)
    }

    cfg, err := config.Load()
    if err != nil {
        fmt.Printf("Warning: Could not load config file: %v\n", err)
    }

    // Resolve database path
    dbPath := resolveDatabasePath(*kdbPath, cfg)
    if dbPath == "" {
        fmt.Println("Error: No KeePass database path provided")
        os.Exit(1)
    }

    // Get database password
    password, err := keepass.ResolvePassword(*kdbPass)
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

    // Search for items
    finder := search.NewFinder(db)
    results, err := finder.Find(*item)
    if err != nil {
        fmt.Printf("Error searching for item: %v\n", err)
        os.Exit(1)
    }

    if len(results) == 0 {
        fmt.Println("No items found")
        os.Exit(0)
    }

    if len(results) > 1 {
        fmt.Println("Multiple items found:")
        for _, result := range results {
            fmt.Printf("- %s\n", result.Path)
        }
        os.Exit(0)
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

func resolveDatabasePath(flagPath string, cfg *config.Config) string {
    if flagPath != "" {
        return flagPath
    }
    if envPath := os.Getenv("KDBPATH"); envPath != "" {
        return envPath
    }
    if cfg != nil && cfg.DatabasePath != "" {
        return cfg.DatabasePath
    }
    return ""
}

func resolveOutputType(flagOut string, cfg *config.Config) output.Type {
    if flagOut != "" {
        return output.Type(flagOut)
    }
    if cfg != nil && cfg.DefaultOutput != "" {
        return output.Type(cfg.DefaultOutput)
    }
    return output.Stdout
}
