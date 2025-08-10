//go:build windows

package cmd

import (
	"os/exec"
	"syscall"
)

// setDetachedProcessAttributes configures the command for detached execution on Windows.
//
// It sets the CreationFlags attribute to 0x00000008 (DETACHED_PROCESS), which starts the process
// in a new console session, detaching it from the parent. This is commonly used for daemonizing
// processes on Windows. The hex value is used directly for cross-compilation compatibility.
//
// Parameters:
//   - cmd: The exec.Cmd object to configure for detached execution.
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
