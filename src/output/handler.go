package output

import (
	"encoding/json"
	"fmt"
	"kpasscli/src/config"
	"kpasscli/src/debug"

	// Hinzugefügt für Debug-Logs
	"github.com/tobischo/gokeepasslib/v3"
	"golang.design/x/clipboard"
)

type Type string

const (
	Clipboard  Type = "clipboard"
	Stdout     Type = "stdout"
	lineBreak       = "----------------------------------------"
	timeFormat      = "2006-01-02 15:04:05"
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

// ShowAllFields displays all fields of a KeePass entry
func ShowAllFields(entry *gokeepasslib.Entry) {
	if entry == nil {
		return
	}

	if config.Config.OutputFormat == "json" {
		showAllFieldsJson(entry)
		return
	}

	fmt.Println(lineBreak)
	fmt.Printf("Entry Details:\n")
	fmt.Println(lineBreak)

	// Standard fields
	printNonEmptyValue("Title", getValue(entry, "Title"))
	printNonEmptyValue("Username", getValue(entry, "UserName"))
	printNonEmptyValue("URL", getValue(entry, "URL"))
	printNonEmptyValue("Notes", getValue(entry, "Notes"))

	// Additional fields
	hasAdditionalFields := false
	for _, v := range entry.Values {
		if isAdditionalField(v.Key) && v.Value.Content != "" {
			if !hasAdditionalFields {
				fmt.Println(lineBreak)
				fmt.Println("Additional Fields:")
				hasAdditionalFields = true
			}
			printNonEmptyValue(v.Key, v.Value.Content)
		}
	}

	// Metadata
	fmt.Println(lineBreak)
	fmt.Println("Metadata:")
	printNonEmptyValue("Created", entry.Times.CreationTime.Format(timeFormat))
	printNonEmptyValue("Modified", entry.Times.LastModificationTime.Format(timeFormat))
	printNonEmptyValue("Accessed", entry.Times.LastAccessTime.Format(timeFormat))
	fmt.Println(lineBreak)
}

// Helper functions
func getValue(entry *gokeepasslib.Entry, key string) string {
	for _, v := range entry.Values {
		if v.Key == key {
			return v.Value.Content
		}
	}
	return ""
}

func printNonEmptyValue(key, value string) {
	if value != "" {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func isAdditionalField(key string) bool {
	standardFields := map[string]bool{
		"Title":    true,
		"UserName": true,
		"URL":      true,
		"Notes":    true,
		"Password": true,
	}
	return !standardFields[key]
}

func showAllFieldsJson(entry *gokeepasslib.Entry) {
	type entryData struct {
		Title            string            `json:"title"`
		Username         string            `json:"username"`
		URL              string            `json:"url"`
		Notes            string            `json:"notes"`
		AdditionalFields map[string]string `json:"additional_fields,omitempty"`
		Metadata         struct {
			Created  string `json:"created"`
			Modified string `json:"modified"`
			Accessed string `json:"accessed"`
		} `json:"metadata"`
	}

	data := entryData{
		Title:            getValue(entry, "Title"),
		Username:         getValue(entry, "UserName"),
		URL:              getValue(entry, "URL"),
		Notes:            getValue(entry, "Notes"),
		AdditionalFields: make(map[string]string),
	}

	// Fill additional fields
	for _, v := range entry.Values {
		if isAdditionalField(v.Key) && v.Value.Content != "" {
			data.AdditionalFields[v.Key] = v.Value.Content
		}
	}

	// Fill metadata
	data.Metadata.Created = entry.Times.CreationTime.Format(timeFormat)
	data.Metadata.Modified = entry.Times.LastModificationTime.Format(timeFormat)
	data.Metadata.Accessed = entry.Times.LastAccessTime.Format(timeFormat)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Error.Printf("Error creating JSON output: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
