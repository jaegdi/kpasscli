package output

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tobischo/gokeepasslib/v3"
	wrappers "github.com/tobischo/gokeepasslib/v3/wrappers"
	"golang.design/x/clipboard"

	"kpasscli/src/config"
)

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
	h := &stdHandler{outputType: OutputType("unknown")}
	err := h.Output("test")
	if err == nil {
		t.Error("expected error for unknown output type")
	}
}

// The following tests are for coverage only; they do not test clipboard or stdout side effects.
func TestStdHandler_toClipboard(t *testing.T) {
	h := &stdHandler{outputType: ClipboardType, clipboard: &RealClipboard{}}
	// We expect an error or nil, but do not check clipboard contents.
	_ = h.toClipboard("test")
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
	h := NewHandler(StdoutType, nil)
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
	h := NewHandler(ClipboardType, nil)
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

func TestShowAllFields_Text(t *testing.T) {
	entry := &gokeepasslib.Entry{
		Values: []gokeepasslib.ValueData{
			{Key: "Title", Value: gokeepasslib.V{Content: "TestTitle"}},
			{Key: "UserName", Value: gokeepasslib.V{Content: "TestUser"}},
			{Key: "URL", Value: gokeepasslib.V{Content: "http://example.com"}},
			{Key: "Notes", Value: gokeepasslib.V{Content: "Some notes"}},
			{Key: "Custom", Value: gokeepasslib.V{Content: "CustomValue"}},
		},
		Times: gokeepasslib.TimeData{
			CreationTime:         &wrappers.TimeWrapper{Time: time.Now()},
			LastModificationTime: &wrappers.TimeWrapper{Time: time.Now()},
			LastAccessTime:       &wrappers.TimeWrapper{Time: time.Now()},
		},
	}
	cfg := config.Config{OutputFormat: "text"}
	ShowAllFields(entry, cfg)
}

func TestShowAllFields_JSON(t *testing.T) {
	entry := &gokeepasslib.Entry{
		Values: []gokeepasslib.ValueData{
			{Key: "Title", Value: gokeepasslib.V{Content: "TestTitle"}},
			{Key: "UserName", Value: gokeepasslib.V{Content: "TestUser"}},
			{Key: "URL", Value: gokeepasslib.V{Content: "http://example.com"}},
			{Key: "Notes", Value: gokeepasslib.V{Content: "Some notes"}},
			{Key: "Custom", Value: gokeepasslib.V{Content: "CustomValue"}},
		},
		Times: gokeepasslib.TimeData{
			CreationTime:         &wrappers.TimeWrapper{Time: time.Now()},
			LastModificationTime: &wrappers.TimeWrapper{Time: time.Now()},
			LastAccessTime:       &wrappers.TimeWrapper{Time: time.Now()},
		},
	}
	cfg := config.Config{OutputFormat: "json"}
	ShowAllFields(entry, cfg)
}

func Test_getValue(t *testing.T) {
	entry := &gokeepasslib.Entry{
		Values: []gokeepasslib.ValueData{{Key: "foo", Value: gokeepasslib.V{Content: "bar"}}},
	}
	if got := getValue(entry, "foo"); got != "bar" {
		t.Errorf("getValue = %q, want 'bar'", got)
	}
	if got := getValue(entry, "baz"); got != "" {
		t.Errorf("getValue for missing key = %q, want ''", got)
	}
}

func Test_printNonEmptyValue(t *testing.T) {
	printNonEmptyValue("Key", "Value") // Should print
	printNonEmptyValue("Key", "")      // Should not print
}

func Test_formatTime(t *testing.T) {
	tm := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	if got := formatTime(tm); got == "" {
		t.Error("formatTime returned empty string")
	}
}

func Test_isAdditionalField(t *testing.T) {
	if isAdditionalField("Title") {
		t.Error("Title should not be additional field")
	}
	if !isAdditionalField("Custom") {
		t.Error("Custom should be additional field")
	}
}

func Test_showAllFieldsJson(t *testing.T) {
	entry := &gokeepasslib.Entry{
		Values: []gokeepasslib.ValueData{
			{Key: "Title", Value: gokeepasslib.V{Content: "TestTitle"}},
			{Key: "Custom", Value: gokeepasslib.V{Content: "CustomValue"}},
		},
		Times: gokeepasslib.TimeData{
			CreationTime:         &wrappers.TimeWrapper{Time: time.Now()},
			LastModificationTime: &wrappers.TimeWrapper{Time: time.Now()},
			LastAccessTime:       &wrappers.TimeWrapper{Time: time.Now()},
		},
	}
	showAllFieldsJson(entry)
}
