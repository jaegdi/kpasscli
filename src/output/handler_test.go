package output

import (
	"os"
	"strings"
	"testing"

	"golang.design/x/clipboard"
)

// --- Additional tests for output package ---
func TestNewHandler(t *testing.T) {
	h := NewHandler(Stdout)
	if h == nil {
		t.Error("NewHandler returned nil")
	}
	if _, ok := h.(*stdHandler); !ok {
		t.Error("NewHandler did not return *stdHandler")
	}
}

func TestIsValidType(t *testing.T) {
	if !IsValidType("stdout") {
		t.Error("IsValidType should return true for 'stdout'")
	}
	if !IsValidType("clipboard") {
		t.Error("IsValidType should return true for 'clipboard'")
	}
	if IsValidType("invalid") {
		t.Error("IsValidType should return false for 'invalid'")
	}
}

func TestStdHandler_Output_UnknownType(t *testing.T) {
	h := &stdHandler{outputType: Type("unknown")}
	err := h.Output("test")
	if err == nil {
		t.Error("expected error for unknown output type")
	}
}

// The following tests are for coverage only; they do not test clipboard or stdout side effects.
func TestStdHandler_toClipboard(t *testing.T) {
	h := &stdHandler{outputType: Clipboard}
	// We expect an error or nil, but do not check clipboard contents.
	_ = h.toClipboard("test")
}

func TestStdHandler_toStdout(t *testing.T) {
	h := &stdHandler{outputType: Stdout}
	// Should not panic or error
	err := h.toStdout("stdhandler to stdout test")
	if err != nil {
		t.Errorf("toStdout returned error: %v", err)
	}
}

type testHandler struct {
	called bool
	msg    string
}

func (h *testHandler) Output(msg string) error {
	h.called = true
	h.msg = msg
	return nil
}

func TestHandler_Output(t *testing.T) {
	h := &testHandler{}
	var handler Handler = h
	err := handler.Output("hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !h.called || h.msg != "hello" {
		t.Errorf("expected Output to be called with 'hello', got called=%v, msg=%v", h.called, h.msg)
	}
}

func TestOutput_Output_Stdout(t *testing.T) {
	h := NewHandler(Stdout)
	expected := "stdout test inkl. Debug verification"
	// Capture stdout
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := h.Output(expected)
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Errorf("Output(Stdout) returned error: %v", err)
	}
	// Read from pipe
	var buf [128]byte
	n, _ := r.Read(buf[:])
	got := string(buf[:n])
	// Remove trailing newline for comparison
	got = strings.TrimSuffix(got, "\n")
	if got != expected {
		t.Errorf("stdout = %q, want %q", got, expected)
	}
}

func TestOutput_Output_Clipboard(t *testing.T) {
	h := NewHandler(Clipboard)
	expected := "clipboard test"
	err := h.Output(expected)
	if err != nil {
		t.Errorf("Output(Clipboard) returned error: %v", err)
	}
	// Check clipboard content
	if err := clipboard.Init(); err == nil {
		got := string(clipboard.Read(clipboard.FmtText))
		if got != expected {
			t.Errorf("clipboard content = %q, want %q", got, expected)
		}
	}
}
