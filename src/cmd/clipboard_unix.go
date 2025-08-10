//go:build !windows

package cmd

import (
	"os/exec"
	"syscall"
)

// setDetachedProcessAttributes configures the command for detached execution on Unix-like systems.
//
// It sets the Setsid attribute to true, which starts the process in a new session,
// detaching it from the controlling terminal. This is commonly used for daemonizing
// processes on Unix/Linux.
//
// Parameters:
//   - cmd: The exec.Cmd object to configure for detached execution.
func setDetachedProcessAttributes(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// Unter Linux/Unix: Prozess in einer neuen Sitzung starten, um ihn vom Terminal zu lösen
	cmd.SysProcAttr.Setsid = true
}
