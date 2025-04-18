//go:build !windows

package cmd

import (
	"os/exec"
	"syscall"
)

// setDetachedProcessAttributes configures the command for detached execution on Unix-like systems.
func setDetachedProcessAttributes(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// Unter Linux/Unix: Prozess in einer neuen Sitzung starten, um ihn vom Terminal zu lösen
	cmd.SysProcAttr.Setsid = true
}
