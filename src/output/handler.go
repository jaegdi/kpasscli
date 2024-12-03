package output

import (
	"fmt"
	"kpasscli/src/debug"

	// Hinzugefügt für Debug-Logs
	"golang.design/x/clipboard"
)

type Type string

const (
	Clipboard Type = "clipboard"
	Stdout    Type = "stdout"
)

type Handler struct {
	outputType Type
}

// NewHandler creates a new Handler instance with the specified output type.
// Parameters:
//   - outputType: The type of output (clipboard or stdout).
//
// Returns:
//   - *Handler: A new Handler instance.
func NewHandler(outputType Type) *Handler {
	return &Handler{outputType: outputType}
}

// Output outputs the given value based on the handler's output type.
// Parameters:
//   - value: The value to be output.
//
// Returns:
//   - error: Any error encountered during output.
func (h *Handler) Output(value string) error {
	debug.Log("Outputting value: %s", value) // Debug-Log hinzugefügt
	switch h.outputType {
	case Clipboard:
		return h.toClipboard(value)
	case Stdout:
		return h.toStdout(value)
	default:
		return fmt.Errorf("unknown output type: %s", h.outputType)
	}
}

// toClipboard copies the given value to the system clipboard.
// Parameters:
//   - value: The value to be copied to the clipboard.
//
// Returns:
//   - error: Any error encountered during the clipboard operation.
func (h *Handler) toClipboard(value string) error {
	debug.Log("Copying to clipboard: %s", value) // Debug-Log hinzugefügt
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("failed to initialize clipboard: %v", err)
	}
	clipboard.Write(clipboard.FmtText, []byte(value))
	return nil
}

// toStdout prints the given value to the standard output.
// Parameters:
//   - value: The value to be printed to stdout.
//
// Returns:
//   - error: Any error encountered during the stdout operation.
func (h *Handler) toStdout(value string) error {
	debug.Log("Printing to stdout: %s", value) // Debug-Log hinzugefügt
	fmt.Println(value)
	return nil
}

// IsValidType checks if the provided output type is valid.
// Parameters:
//   - outputType: The output type to check.
//
// Returns:
//   - bool: True if the output type is valid, false otherwise.
func IsValidType(outputType string) bool {
	switch Type(outputType) {
	case Clipboard, Stdout:
		return true
	default:
		return false
	}
}
