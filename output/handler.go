package output

import (
	"fmt"
	"kpasscli/debug"

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

func NewHandler(outputType Type) *Handler {
	return &Handler{outputType: outputType}
}

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

func (h *Handler) toClipboard(value string) error {
	debug.Log("Copying to clipboard: %s", value) // Debug-Log hinzugefügt
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("failed to initialize clipboard: %v", err)
	}
	clipboard.Write(clipboard.FmtText, []byte(value))
	return nil
}

func (h *Handler) toStdout(value string) error {
	debug.Log("Printing to stdout: %s", value) // Debug-Log hinzugefügt
	fmt.Println(value)
	return nil
}
