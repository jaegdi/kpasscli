package doc

import (
	"bytes"
	"os"
	"testing"
)

func TestShowHelp(t *testing.T) {
	// Capture stdout
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	ShowHelp()
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	// Optionally: check buf.String() for expected content
}
