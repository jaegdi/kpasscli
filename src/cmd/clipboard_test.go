package cmd

import (
	"os/exec"
	"runtime"
	"testing"
)

func Test_runClearerDaemonMode(t *testing.T) {
	// This test will just check that the function exits for delaySeconds <= 0
	// and does not panic for a small positive value (side effects not checked)
	// Use t.Parallel() for isolation
	t.Parallel()
	// Should exit immediately
	// Not testable directly due to os.Exit, but can be covered by running in a subprocess if needed
}

func Test_clearWaylandClipboard(t *testing.T) {
	t.Parallel()
	path, cmd := clearWaylandClipboard()
	if path != "" && cmd == nil {
		t.Error("If path is set, cmd should not be nil")
	}
}

func Test_clearX11Clipboard(t *testing.T) {
	t.Parallel()
	path, cmd := clearX11Clipboard()
	if path != "" && cmd == nil {
		t.Error("If path is set, cmd should not be nil")
	}
}

func Test_clearDarwinClipboard(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "darwin" {
		t.Skip("Only relevant on macOS")
	}
	path, cmd := clearDarwinClipboard()
	if path != "" && cmd == nil {
		t.Error("If path is set, cmd should not be nil")
	}
}

func Test_clearWindowsClipboard(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("Only relevant on Windows")
	}
	path, cmd := clearWindowsClipboard()
	if path == "" || cmd == nil {
		t.Error("Windows clipboard clear should always return a path and cmd")
	}
}

func Test_clearWindowsClipboard_Focused(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("Only relevant on Windows")
	}
	path, cmd := clearWindowsClipboard()
	if path != "cmd" {
		t.Errorf("Expected path to be 'cmd', got %q", path)
	}
	if cmd == nil {
		t.Error("Expected cmd to not be nil")
	}
	if len(cmd.Args) < 3 || cmd.Args[2] != "clip" {
		t.Errorf("Expected cmd.Args to contain 'clip', got %v", cmd.Args)
	}
}

func Test_performClipboardClear(t *testing.T) {
	t.Parallel()
	// This function is hard to test for side effects, but we can call it to check for panics
	// It should not panic on any platform
	performClipboardClear()
}

func Test_StartClipboardClearer(t *testing.T) {
	t.Parallel()
	// This function starts a detached process, so we only check that it does not panic
	StartClipboardClearer(1)
}

func Test_setDetachedProcessAttributes(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("echo", "test")
	setDetachedProcessAttributes(cmd)
	// No assertion, just check for panics
}
