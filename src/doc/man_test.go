package doc

import (
	"bytes"
	"os"
	"testing"
)

func TestShowMan(t *testing.T) {
	// Capture stdout
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	ShowMan()
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	// Optionally: check buf.String() for expected content
}
