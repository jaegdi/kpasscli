//go:build windows

package cmd

import (
	"os/exec"
	"syscall"
)

// setDetachedProcessAttributes configures the command for detached execution on Windows.
func setDetachedProcessAttributes(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// Unter Windows: Prozess in einer neuen Konsolensitzung starten (DETACHED_PROCESS)
	// Verwende den Hex-Wert, da die Konstante syscall.DETACHED_PROCESS
	// beim Cross-Compilieren möglicherweise nicht definiert ist.
	// DETACHED_PROCESS = 0x00000008
	cmd.SysProcAttr.CreationFlags = 0x00000008
}
