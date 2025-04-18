package cmd

import (
	"kpasscli/src/debug"
	"log"
	"os"
	"os/exec"
	"runtime" // Kann für performClipboardClear nützlich sein
	"strconv"

	// "syscall" // Wird hier nicht mehr direkt benötigt
	"time"
)

// HINWEIS: Die Funktion setDetachedProcessAttributes wurde entfernt.
// Ihre Definition befindet sich jetzt ausschließlich in
// clipboard_unix.go und clipboard_windows.go

// --- runClearerDaemonMode bleibt unverändert ---
func runClearerDaemonMode(delaySeconds int) {
	if delaySeconds <= 0 {
		os.Exit(1) // Sollte nicht passieren, aber sicher ist sicher
	}

	// Warte die angegebene Zeit
	time.Sleep(time.Duration(delaySeconds) * time.Second)

	// Führe die eigentliche Löschaktion aus
	performClipboardClear()

	os.Exit(0) // Erfolgreich beendet nach dem Löschen (oder Versuch)
}

// --- performClipboardClear bleibt unverändert ---
func performClipboardClear() {
	debug.Log("Attempting to clear clipboard.")

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
		if err != nil {
			log.Printf("Warning: Could not find 'wl-copy', 'xclip', or 'xsel' to clear Linux clipboard: %v", err)
			return
		}

	case "darwin": // macOS
		clearCmdPath, err = exec.LookPath("pbcopy")
		if err == nil {
			cmd = exec.Command(clearCmdPath)
			stdin, pipeErr := cmd.StdinPipe()
			if pipeErr != nil {
				log.Printf("Warning: Could not get stdin pipe for pbcopy: %v", pipeErr)
				return
			}
			closeErr := stdin.Close() // Send EOF immediately
			if closeErr != nil {
				log.Printf("Warning: Could not close stdin pipe for pbcopy: %v", closeErr)
			}
			debug.Log("Using pbcopy to clear macOS clipboard")
		} else {
			log.Printf("Warning: Could not find 'pbcopy' to clear macOS clipboard: %v", err)
			return
		}

	case "windows":
		clearCmdPath = "cmd"
		args = []string{"/c", "echo off | clip"}
		debug.Log("Using 'cmd /c echo off | clip' to clear Windows clipboard")

	default:
		log.Printf("Warning: Automatic clipboard clearing not supported on this OS: %s", runtime.GOOS)
		return
	}

	if cmd == nil && clearCmdPath != "" {
		cmd = exec.Command(clearCmdPath, args...)
	}

	if cmd != nil {
		runErr := cmd.Run() // Run and wait here, as this *is* the clearer process
		if runErr != nil {
			log.Printf("Warning: Failed to clear clipboard: %v", runErr)
		} else {
			debug.Log("Clipboard cleared successfully.")
		}
	} else if clearCmdPath == "" && runtime.GOOS != "darwin" {
		log.Printf("Warning: No suitable command found to clear clipboard on %s.", runtime.GOOS)
	}
}

// Funktion zum Starten des *neuen Prozesses* zum Löschen des Clipboards
func StartClipboardClearer(delaySeconds int) {
	if delaySeconds <= 0 {
		return
	}

	debug.Log("Attempting to start detached clipboard clearer process with delay: %d seconds", delaySeconds)

	// 1. Finde den Pfad zur aktuell laufenden kpasscli-Executable
	executablePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Could not find executable path to start clearer process: %v", err)
		return
	}

	// 2. Bereite die Argumente für den neuen Prozess vor
	args := []string{
		"--internal-clear-clipboard",                         // Das interne Flag
		"--internal-clear-delay", strconv.Itoa(delaySeconds), // Das Delay übergeben
	}
	// Optional: Debug-Flag weitergeben
	if debug.Enabled() {
		args = append(args, "--debug") // oder "-d"
	}

	// 3. Erstelle das exec.Cmd Objekt
	cmd := exec.Command(executablePath, args...)

	// 4. WICHTIG: Prozess losgelöst starten (detached)
	//    Rufe die plattformspezifische Funktion auf, um SysProcAttr zu setzen.
	setDetachedProcessAttributes(cmd) // <--- Der Aufruf bleibt bestehen

	err = cmd.Start() // Startet den Prozess und kehrt sofort zurück
	if err != nil {
		log.Printf("Warning: Failed to start detached clearer process: %v", err)
		return
	}

	debug.Log("Detached clearer process started successfully (PID: %d). Main process will now exit.", cmd.Process.Pid)

	// WICHTIG: Wir rufen cmd.Process.Release() NICHT auf, da der Prozess weiterlaufen soll.
	// Wir rufen cmd.Wait() NICHT auf.
}
