package output

import (
	"fmt"
	"kpasscli/src/config"
	"kpasscli/src/debug"
	"os"

	// Hinzugefügt für Debug-Logs
	"golang.design/x/clipboard"
)

type OutputChannel string

const (
	Clipboard OutputChannel = "clipboard"
	Stdout    OutputChannel = "stdout"
)

type Handler struct {
	outputChannel OutputChannel
}

// NewHandler creates a new Handler instance with the specified output type.
// Parameters:
//   - outputChannel: The type of output (clipboard or stdout).
//
// Returns:
//   - *Handler: A new Handler instance.
func NewHandler(outputChannel OutputChannel) *Handler {
	return &Handler{outputChannel: outputChannel}
}

// Output outputs the given value based on the handler's output type.
// Parameters:
//   - value: The value to be output.
//
// Returns:
//   - error: Any error encountered during output.
func (h *Handler) Output(value string) error {
	debug.Log("Outputting value: %s", value) // Debug-Log hinzugefügt
	switch h.outputChannel {
	case Clipboard:
		return h.toClipboard(value)
	case Stdout:
		return h.toStdout(value)
	default:
		return fmt.Errorf("unknown output type: %s", h.outputChannel)
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
//   - outputChannel: The output type to check.
//
// Returns:
//   - bool: True if the output type is valid, false otherwise.
func IsValidType(outputChannel string) bool {
	switch OutputChannel(outputChannel) {
	case Clipboard, Stdout:
		return true
	default:
		return false
	}
}

// resolveOutputType determines the output type based on the provided flag,
// environment variable, or configuration. It follows this order of precedence:
// 1. If the flagOut parameter is not empty, it returns the corresponding output type.
// 2. If the environment variable "KPASSCLI_OUT" is set and valid, it returns the corresponding output type.
// 3. If the cfg parameter is not nil and cfg.DefaultOutput is not empty, it returns the corresponding output type.
// 4. If none of the above conditions are met, it defaults to output.Stdout.
//
// Parameters:
// - flagOut: A string representing the output type specified by a flag.
// - cfg: A pointer to a config.Config struct that may contain a default output type.
//
// Returns:
// - output.OutputChannel: The resolved output type based on the provided inputs.
func ResolveOutputType(flagOut string, cfg *config.Config) OutputChannel {
	if flagOut != "" {
		return OutputChannel(flagOut)
	}
	if kpcliout := os.Getenv("KPASSCLI_OUT"); kpcliout != "" {
		if IsValidType(kpcliout) {
			return OutputChannel(kpcliout)
		}
	}
	if cfg != nil && cfg.DefaultOutput != "" {
		return OutputChannel(cfg.DefaultOutput)
	}
	return Stdout
}
