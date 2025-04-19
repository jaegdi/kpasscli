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

// runClearerDaemonMode ist die Funktion, die im separaten, losgelösten Prozess ausgeführt wird.
// Sie wartet die angegebene Anzahl von Sekunden und ruft dann performClipboardClear auf,
// um die Zwischenablage zu löschen. Danach beendet sie den Prozess.
// Parameter:
//   - delaySeconds: Die Anzahl der Sekunden, die gewartet werden soll, bevor die Zwischenablage gelöscht wird.
func runClearerDaemonMode(delaySeconds int) {
	if delaySeconds <= 0 {
		os.Exit(1) // Sollte nicht passieren, aber sicher ist sicher
	}

	// Warte die angegebene Zeit
	debug.Log("Clearer daemon started, waiting for %d seconds.", delaySeconds)
	time.Sleep(time.Duration(delaySeconds) * time.Second)

	// Führe die eigentliche Löschaktion aus
	debug.Log("Clearer daemon delay finished, attempting to clear clipboard.")
	performClipboardClear()

	debug.Log("Clearer daemon finished.")
	os.Exit(0) // Erfolgreich beendet nach dem Löschen (oder Versuch)
}

// clearWaylandClipboard versucht, den Befehl zum Löschen der Zwischenablage unter Wayland zu finden und vorzubereiten.
// Es priorisiert 'wl-copy' und fällt auf 'xclip' oder 'xsel' zurück, falls 'wl-copy' nicht verfügbar ist.
// Rückgabewerte:
//   - clearCmdPath: Der Pfad zum gefundenen Befehl (oder leer, wenn keiner gefunden wurde).
//   - cmd: Ein vorbereitetes *exec.Cmd Objekt (kann nil sein, wenn pbcopy spezielle Behandlung braucht oder kein Befehl gefunden wurde).
func clearWaylandClipboard() (clearCmdPath string, cmd *exec.Cmd) {
	var err error
	var args []string
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
			} else {
				// Kein geeignetes Kommando gefunden
				clearCmdPath = "" // Sicherstellen, dass der Pfad leer ist
			}
		}
	}
	// Erstelle das Cmd-Objekt nur, wenn ein Pfad gefunden wurde
	if clearCmdPath != "" {
		cmd = exec.Command(clearCmdPath, args...)
	}
	return clearCmdPath, cmd
}

// clearX11Clipboard versucht, den Befehl zum Löschen der Zwischenablage unter X11 zu finden und vorzubereiten.
// Es priorisiert 'xclip' und fällt auf 'xsel' zurück, falls 'xclip' nicht verfügbar ist.
// Rückgabewerte:
//   - clearCmdPath: Der Pfad zum gefundenen Befehl (oder leer, wenn keiner gefunden wurde).
//   - cmd: Ein vorbereitetes *exec.Cmd Objekt (oder nil, wenn kein Befehl gefunden wurde).
func clearX11Clipboard() (clearCmdPath string, cmd *exec.Cmd) {
	var err error
	var args []string
	clearCmdPath, err = exec.LookPath("xclip")
	if err == nil {
		args = []string{"-selection", "clipboard", "-in", "/dev/null"}
		debug.Log("Using xclip to clear X11 clipboard")
	} else {
		clearCmdPath, err = exec.LookPath("xsel")
		if err == nil {
			args = []string{"-b", "-c"}
			debug.Log("Using xsel to clear X11 clipboard")
		} else {
			// Kein geeignetes Kommando gefunden
			clearCmdPath = "" // Sicherstellen, dass der Pfad leer ist
			log.Printf("Warning: Could not find 'xclip' or 'xsel' to clear Linux X11 clipboard: %v", err)
		}
	}
	// Erstelle das Cmd-Objekt nur, wenn ein Pfad gefunden wurde
	if clearCmdPath != "" {
		cmd = exec.Command(clearCmdPath, args...)
	}
	return clearCmdPath, cmd
}

// clearDarwinClipboard versucht, den Befehl 'pbcopy' zum Löschen der Zwischenablage unter macOS (Darwin) vorzubereiten.
// Das Löschen erfolgt, indem ein leerer Input an 'pbcopy' gesendet wird (EOF auf stdin).
// Rückgabewerte:
//   - clearCmdPath: Der Pfad zu 'pbcopy' (oder leer, wenn nicht gefunden).
//   - cmd: Ein vorbereitetes *exec.Cmd Objekt mit konfigurierter StdinPipe (oder nil, wenn 'pbcopy' nicht gefunden wurde).
func clearDarwinClipboard() (clearCmdPath string, cmd *exec.Cmd) {
	var err error
	clearCmdPath, err = exec.LookPath("pbcopy")
	if err == nil {
		cmd = exec.Command(clearCmdPath) // Keine Argumente nötig
		stdin, pipeErr := cmd.StdinPipe()
		if pipeErr != nil {
			log.Printf("Warning: Could not get stdin pipe for pbcopy: %v", pipeErr)
			cmd = nil // Fehler beim Vorbereiten, kein Cmd zurückgeben
			return
		}
		// Sende sofort EOF, um die Zwischenablage zu leeren
		closeErr := stdin.Close()
		if closeErr != nil {
			log.Printf("Warning: Could not close stdin pipe for pbcopy: %v", closeErr)
			// Fahren trotzdem fort, vielleicht funktioniert es dennoch
		}
		debug.Log("Using pbcopy to clear macOS clipboard")
	} else {
		log.Printf("Warning: Could not find 'pbcopy' to clear macOS clipboard: %v", err)
		clearCmdPath = "" // Sicherstellen, dass der Pfad leer ist
		cmd = nil
	}
	return clearCmdPath, cmd
}

// clearWindowsClipboard bereitet den Befehl zum Löschen der Zwischenablage unter Windows vor.
// Verwendet 'cmd /c echo off | clip'.
// Rückgabewerte:
//   - clearCmdPath: Der Pfad zu 'cmd'.
//   - cmd: Ein vorbereitetes *exec.Cmd Objekt.
func clearWindowsClipboard() (clearCmdPath string, cmd *exec.Cmd) {
	clearCmdPath = "cmd"
	args := []string{"/c", "echo off | clip"}
	debug.Log("Using 'cmd /c echo off | clip' to clear Windows clipboard")
	cmd = exec.Command(clearCmdPath, args...) // Erstelle das Cmd-Objekt hier
	return clearCmdPath, cmd
}

// performClipboardClear führt die eigentliche Aktion zum Löschen der Zwischenablage aus.
// Es erkennt das Betriebssystem, ruft die entsprechende Hilfsfunktion (clear*Clipboard) auf,
// um den Befehl zu erhalten, und führt diesen dann aus.
// Diese Funktion wird vom Hintergrundprozess (gestartet durch StartClipboardClearer) aufgerufen.
func performClipboardClear() {
	debug.Log("Attempting to clear clipboard.")

	var cmd *exec.Cmd
	var clearCmdPath string // Nur zur Information, ob ein Befehl gefunden wurde

	switch runtime.GOOS {
	case "linux":
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			debug.Log("Detected Wayland display.")
			clearCmdPath, cmd = clearWaylandClipboard()
		} else { // Annahme: X11
			debug.Log("Assuming X11 display.")
			clearCmdPath, cmd = clearX11Clipboard()
		}
	case "darwin": // macOS
		clearCmdPath, cmd = clearDarwinClipboard()
	case "windows":
		clearCmdPath, cmd = clearWindowsClipboard()
	default:
		log.Printf("Warning: Automatic clipboard clearing not supported on this OS: %s", runtime.GOOS)
		return
	}

	// cmd sollte jetzt entweder ein gültiges *exec.Cmd sein oder nil, wenn kein Befehl gefunden/vorbereitet werden konnte.
	if cmd != nil {
		debug.Log("Executing clipboard clear command: %s %v", cmd.Path, cmd.Args)
		runErr := cmd.Run() // Führe den Befehl aus und warte auf das Ergebnis
		if runErr != nil {
			// Bei pbcopy kann ein Fehler auftreten, wenn stdin geschlossen wird, was aber erwartet wird.
			// Prüfe, ob es sich um pbcopy handelt und der Fehler spezifisch ist (optional, falls nötig).
			// Fürs Erste loggen wir den Fehler immer.
			log.Printf("Warning: Failed to clear clipboard: %v", runErr)
		} else {
			debug.Log("Clipboard cleared successfully.")
		}
	} else if clearCmdPath == "" { // Explizite Prüfung, ob überhaupt ein Pfad gefunden wurde
		// Diese Meldung wird bereits in den clear* funktionen geloggt, aber hier nochmal zur Sicherheit.
		log.Printf("Warning: No suitable command found or prepared to clear clipboard on %s.", runtime.GOOS)
	} else {
		// Fall, der eigentlich nicht eintreten sollte: Pfad gefunden, aber cmd ist nil (z.B. Fehler bei pbcopy pipe)
		log.Printf("Warning: Command path found (%s), but command execution object could not be prepared.", clearCmdPath)
	}
}

// StartClipboardClearer startet einen *neuen*, vom Hauptprozess losgelösten (detached) Prozess,
// der nach einer bestimmten Verzögerung die Zwischenablage löscht.
// Der neue Prozess ist eine Instanz der aktuell laufenden kpasscli-Executable,
// die mit speziellen internen Flags aufgerufen wird.
// Parameter:
//   - delaySeconds: Die Anzahl der Sekunden, die der neue Prozess warten soll, bevor er die Zwischenablage löscht.
func StartClipboardClearer(delaySeconds int) {
	if delaySeconds <= 0 {
		debug.Log("Clipboard clearing disabled (delay <= 0).")
		return
	}

	debug.Log("Attempting to start detached clipboard clearer process with delay: %d seconds", delaySeconds)

	// 1. Finde den Pfad zur aktuell laufenden kpasscli-Executable
	executablePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Could not find executable path to start clearer process: %v", err)
		return
	}
	debug.Log("Executable path for clearer: %s", executablePath)

	// 2. Bereite die Argumente für den neuen Prozess vor
	args := []string{
		"--internal-clear-clipboard",                         // Das interne Flag, um runClearerDaemonMode auszulösen
		"--internal-clear-delay", strconv.Itoa(delaySeconds), // Das Delay übergeben
	}
	// Optional: Debug-Flag weitergeben, damit der Hintergrundprozess auch loggt
	if debug.Enabled() {
		args = append(args, "--debug") // oder "-d"
		debug.Log("Passing debug flag to clearer process.")
	}
	debug.Log("Arguments for clearer process: %v", args)

	// 3. Erstelle das exec.Cmd Objekt
	cmd := exec.Command(executablePath, args...)

	// 4. WICHTIG: Prozess losgelöst starten (detached)
	//    Rufe die plattformspezifische Funktion auf, um SysProcAttr zu setzen.
	//    Dies stellt sicher, dass der neue Prozess weiterläuft, auch wenn der Hauptprozess endet.
	setDetachedProcessAttributes(cmd)
	debug.Log("Set detached process attributes.")

	// 5. Starte den Prozess
	err = cmd.Start() // Startet den Prozess und kehrt sofort zurück
	if err != nil {
		log.Printf("Warning: Failed to start detached clearer process: %v", err)
		return
	}

	// Der Prozess wurde erfolgreich gestartet. Der Hauptprozess kann nun normal beendet werden.
	// Der neue Prozess läuft im Hintergrund weiter.
	debug.Log("Detached clearer process started successfully (PID: %d). Main process will now exit.", cmd.Process.Pid)

	// WICHTIG: Wir rufen cmd.Process.Release() NICHT auf, da der Prozess weiterlaufen soll.
	// Wir rufen cmd.Wait() NICHT auf, da wir nicht auf das Ende des Hintergrundprozesses warten wollen.
}
