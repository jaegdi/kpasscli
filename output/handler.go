package output

import (
	"fmt"

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
    if err := clipboard.Init(); err != nil {
        return fmt.Errorf("failed to initialize clipboard: %v", err)
    }
    clipboard.Write(clipboard.FmtText, []byte(value))
    return nil
}

func (h *Handler) toStdout(value string) error {
    fmt.Println(value)
    return nil
}
