package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"time"

	"github.com/tobischo/gokeepasslib/v3"
	"golang.design/x/clipboard"

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
	newHandler func(output.OutputType, output.ClipboardService) output.Handler,
	clipboardService output.ClipboardService,
	getEnv func(string) string,
) error {
	// Handle clipboard clearing background process
	if flags.ClearClipboard {
		if flags.ClearAfter > 0 {
			debug.Log("Waiting %d seconds before clearing clipboard...", flags.ClearAfter)
			time.Sleep(time.Duration(flags.ClearAfter) * time.Second)
		}
		if err := clipboardService.Init(); err != nil {
			return fmt.Errorf("failed to initialize clipboard: %v", err)
		}
		clipboardService.Write(clipboard.FmtText, []byte(""))
		debug.Log("Clipboard cleared.")
		return nil
	}

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
	// removed stray 'c'
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

	outputType := output.ResolveOutputType(flags.Out, config)
	handler := newHandler(outputType, clipboardService)

	var value string
	// var token string
	if flags.TotpFlag {
		// totpSecret, err := results[0].GetField("otp")
		// println("totpSecret:", totpSecret)
		// if err != nil {
		// 	return fmt.Errorf("TOTP secret not found: %w", err)
		// }
		// token, err := totp.GenerateCode(totpSecret, time.Now())
		// if err != nil {
		// 	return fmt.Errorf("Error generating TOTP token: %w", err)
		// }
		value, err = results[0].GetTotpToken("otp")
		if err != nil {
			return fmt.Errorf("Error getting field: %w", err)
		}
	} else {
		value, err = results[0].GetField(flags.FieldName)
		if err != nil {
			return fmt.Errorf("Error getting field: %w", err)
		}

		if flags.PasswordTotp {
			token, err := results[0].GetTotpToken("otp")
			if err != nil {
				return fmt.Errorf("Error generating TOTP token: %w", err)
			}
			value = value + token
		}
	}

	// Output the value of the requested item field
	if err := handler.Output(value); err != nil {
		return fmt.Errorf("Error outputting value: %w", err)
	}

	// If output is clipboard and ClearAfter is set, spawn background process
	if outputType == output.ClipboardType && flags.ClearAfter > 0 {
		exe, err := os.Executable()
		if err != nil {
			debug.Log("Failed to get executable path: %v", err)
		} else {
			debug.Log("Spawning background process to clear clipboard after %d seconds", flags.ClearAfter)
			cmd := exec.Command(exe, "--clear-clipboard", "-ca", strconv.Itoa(flags.ClearAfter))
			if err := cmd.Start(); err != nil {
				debug.Log("Failed to start background process: %v", err)
			} else {
				// Detach process
				cmd.Process.Release()
			}
		}
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
		&output.RealClipboard{},
		os.Getenv,
	)
	if err != nil {
		debug.ErrMsg(err, "kpasscli")
		os.Exit(1)
	}
}
