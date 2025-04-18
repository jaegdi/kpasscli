// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/keepass"
	"kpasscli/src/output"
	"kpasscli/src/search"
)

// Funktion zum Starten des Clear-Daemons (Goroutine)
func startClipboardClearer(delaySeconds int) {
	if delaySeconds <= 0 {
		return // Nichts tun, wenn Verzögerung 0 oder negativ ist
	}

	debug.Log("Starting clipboard clearer goroutine with delay: %d seconds", delaySeconds)

	go func() {
		time.Sleep(time.Duration(delaySeconds) * time.Second)
		debug.Log("Timer finished. Attempting to clear clipboard.")

		var cmd *exec.Cmd
		var clearCmdPath string
		var args []string
		var err error

		switch runtime.GOOS {
		case "linux":
			// Prüfe auf Wayland zuerst
			if os.Getenv("WAYLAND_DISPLAY") != "" {
				clearCmdPath, err = exec.LookPath("wl-copy")
				if err == nil {
					args = []string{"-c"}
					debug.Log("Using wl-copy to clear Wayland clipboard")
				} else {
					debug.Log("wl-copy not found, trying xclip/xsel for Wayland (might not work)")
					// Fallback zu X11-Tools, falls wl-copy nicht da ist, aber unwahrscheinlich, dass es funktioniert
					clearCmdPath, err = exec.LookPath("xclip")
					if err == nil {
						args = []string{"-selection", "clipboard", "-in", "/dev/null"}
						debug.Log("Using xclip as fallback")
					} else {
						clearCmdPath, err = exec.LookPath("xsel")
						if err == nil {
							args = []string{"-b", "-c"}
							debug.Log("Using xsel as fallback")
						}
					}
				}
			} else { // Annahme: X11
				clearCmdPath, err = exec.LookPath("xclip")
				if err == nil {
					args = []string{"-selection", "clipboard", "-in", "/dev/null"}
					debug.Log("Using xclip to clear X11 clipboard")
				} else {
					clearCmdPath, err = exec.LookPath("xsel")
					if err == nil {
						args = []string{"-b", "-c"}
						debug.Log("Using xsel to clear X11 clipboard")
					}
				}
			}
			if err != nil { // Wenn keines der Linux-Tools gefunden wurde
				log.Printf("Warning: Could not find 'wl-copy', 'xclip', or 'xsel' to clear Linux clipboard: %v", err)
				return
			}

		case "darwin": // macOS
			clearCmdPath, err = exec.LookPath("pbcopy")
			if err == nil {
				// pbcopy leert, wenn man ihm nichts übergibt (EOF auf stdin)
				cmd = exec.Command(clearCmdPath)
				stdin, pipeErr := cmd.StdinPipe()
				if pipeErr != nil {
					log.Printf("Warning: Could not get stdin pipe for pbcopy: %v", pipeErr)
					return
				}
				// Sofort schließen, um EOF zu senden
				closeErr := stdin.Close()
				if closeErr != nil {
					log.Printf("Warning: Could not close stdin pipe for pbcopy: %v", closeErr)
					// Trotzdem versuchen auszuführen
				}
				debug.Log("Using pbcopy to clear macOS clipboard")
			} else {
				log.Printf("Warning: Could not find 'pbcopy' to clear macOS clipboard: %v", err)
				return
			}

		case "windows":
			// Der Befehl 'clip' liest von stdin. echo off | clip leert es effektiv.
			clearCmdPath = "cmd" // Führe den Befehl über cmd.exe aus
			args = []string{"/c", "echo off | clip"}
			debug.Log("Using 'cmd /c echo off | clip' to clear Windows clipboard")

		default:
			log.Printf("Warning: Automatic clipboard clearing not supported on this OS: %s", runtime.GOOS)
			return
		}

		// Führe den Befehl aus (außer für macOS pbcopy, wo cmd schon erstellt wurde)
		if cmd == nil && clearCmdPath != "" {
			cmd = exec.Command(clearCmdPath, args...)
		}

		if cmd != nil {
			runErr := cmd.Run()
			if runErr != nil {
				// Gib eine Warnung aus, aber keinen Fehler, da kpasscli selbst erfolgreich war
				log.Printf("Warning: Failed to clear clipboard: %v", runErr)
			} else {
				debug.Log("Clipboard cleared successfully.")
			}
		} else if clearCmdPath == "" && runtime.GOOS != "darwin" { // Nur loggen wenn kein Befehl gefunden wurde (außer macOS)
			log.Printf("Warning: No suitable command found to clear clipboard on %s.", runtime.GOOS)
		}

	}() // Starte die Goroutine
}

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

		err = keepass.GetAllFields(db, cfg, flags.Item)
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

		if len(results) == 1 { // Nur bei genau einem Ergebnis fortfahren
			value, err := results[0].GetField(flags.FieldName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting field: %v\n", err)
				os.Exit(1)
			}

			// Output the value using the handler
			if err := handler.Output(value); err != nil {
				fmt.Fprintf(os.Stderr, "Error outputting value: %v\n", err)
				os.Exit(1)
			}

			// --- NEUER TEIL: Clipboard Clearer starten ---
			if flags.ClearAfter > 0 && outputType == output.Clipboard {
				// Informiere den Benutzer (optional, aber gut)
				fmt.Fprintf(os.Stderr, "Value copied to clipboard. Will clear in %d seconds.\n", flags.ClearAfter)
				startClipboardClearer(flags.ClearAfter)
			}
			// --- ENDE NEUER TEIL ---

			os.Exit(0) // Erfolgreich beendet nach Ausgabe
		} else if len(results) == 0 {
			fmt.Fprintf(os.Stderr, "No items found\n") // Konsistente Fehlermeldung
			os.Exit(1)
		} else { // Mehr als 1 Ergebnis
			fmt.Fprintf(os.Stderr, "Multiple items found:\n") // Konsistente Fehlermeldung
			for _, result := range results {
				fmt.Fprintf(os.Stderr, "- %s\n", result.Path)
			}
			os.Exit(1)
		}
	}
}
