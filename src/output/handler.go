package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	// Hinzugefügt für Debug-Logs

	"github.com/tobischo/gokeepasslib/v3"
	"golang.design/x/clipboard"

	"kpasscli/src/config"
	"kpasscli/src/debug"
)

type OutputChannel string

const (
	Clipboard  OutputChannel = "clipboard"
	Stdout     OutputChannel = "stdout"
	lineBreak                = "----------------------------------------"
	timeFormat               = "2006-01-02 15:04:05"
)

// Handler is an interface for outputting values.
type Handler interface {
	Output(string) error
}

// OutputType defines the type of output (clipboard or stdout)
type OutputType string

const (
	ClipboardType OutputType = "clipboard"
	StdoutType    OutputType = "stdout"
)

// stdHandler is the default implementation of Handler.
type stdHandler struct {
	outputType OutputType
}

// NewHandler creates a new Handler instance with the specified output type.
//
// Parameters:
//   - outputType: The OutputType specifying how output should be handled (clipboard or stdout).
//
// Returns:
//   - Handler: A new Handler instance for the specified output type.
func NewHandler(outputType OutputType) Handler {
	return &stdHandler{outputType: outputType}
}

// Output outputs the given value based on the handler's output type.
// Parameters:
//   - value: The value to be output.
//
// Returns:
//   - error: Any error encountered during output.
func (h *stdHandler) Output(value string) error {
	debug.Log("Outputting value: %s", value)
	switch h.outputType {
	case ClipboardType:
		return h.toClipboard(value)
	case StdoutType:
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
func (h *stdHandler) toClipboard(value string) error {
	debug.Log("Copying to clipboard: %s", value)
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
func (h *stdHandler) toStdout(value string) error {
	debug.Log("Printing to stdout: %s", value)
	fmt.Println(value)
	return nil
}

// IsValidType checks if the provided output type is valid.
//
// Parameters:
//   - outputType: The output type to check (string).
//
// Returns:
//   - bool: True if the output type is valid, false otherwise.
func IsValidType(outputType string) bool {
	switch OutputType(outputType) {
	case ClipboardType, StdoutType:
		return true
	default:
		return false
	}
}

// ShowAllFields displays all fields of a KeePass entry
// in a human-readable format.
// Parameters:
//   - entry: The KeePass entry to display.
//   - config: The configuration object containing output format settings.
//
// Returns:
//   - error: Any error encountered during the operation.
func ShowAllFields(entry *gokeepasslib.Entry, config config.Config) {
	if entry == nil {
		return
	}

	if config.OutputFormat == "json" {
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
	printNonEmptyValue("Created", formatTime(entry.Times.CreationTime.Time))
	printNonEmptyValue("Modified", formatTime(entry.Times.LastModificationTime.Time))
	printNonEmptyValue("Accessed", formatTime(entry.Times.LastAccessTime.Time))
	fmt.Println(lineBreak)
}

// Helper functions
// getValue retrieves the value for a given key from the entry's values.
// Parameters:
//   - entry: The KeePass entry to search.
//   - key: The key for which to retrieve the value.
//
// Returns:
//   - string: The value associated with the key, or an empty string if not found
func getValue(entry *gokeepasslib.Entry, key string) string {
	for _, v := range entry.Values {
		if v.Key == key {
			return v.Value.Content
		}
	}
	return ""
}

// printNonEmptyValue prints a key-value pair if the value is not empty.
// Parameters:
//   - key: The key to print.
//   - value: The value to print.
func printNonEmptyValue(key, value string) {
	if value != "" {
		fmt.Printf("%s: %s\n", key, value)
	}
}

// formatTime formats a time.Time object into a string using the predefined time format.
// Parameters:
//   - t: The time.Time object to format.
//
// Returns:
//   - string: The formatted time string.
func formatTime(t time.Time) string {
	return t.Format(timeFormat)
}

// isAdditionalField checks if a key is considered an additional field.
// Parameters:
//   - key: The key to check.
//
// Returns:
//   - bool: True if the key is an additional field, false otherwise.
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

// showAllFieldsJson outputs all fields of a KeePass entry in JSON format.
// Parameters:
//   - entry: The KeePass entry to display.
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
	data.Metadata.Created = formatTime(entry.Times.CreationTime.Time)
	data.Metadata.Modified = formatTime(entry.Times.LastModificationTime.Time)
	data.Metadata.Accessed = formatTime(entry.Times.LastAccessTime.Time)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON output: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}

// ResolveOutputType determines the output type based on the provided flag, environment variable, or configuration.
//
// Order of precedence:
//  1. If the flagOut parameter is not empty, it returns the corresponding output type.
//  2. If the environment variable "KPASSCLI_OUT" is set and valid, it returns the corresponding output type.
//  3. If the cfg parameter is not nil and cfg.DefaultOutput is not empty, it returns the corresponding output type.
//  4. If none of the above conditions are met, it defaults to output.StdoutType.
//
// Parameters:
//   - flagOut: A string representing the output type specified by a flag.
//   - cfg: A pointer to a config.Config struct that may contain a default output type.
//
// Returns:
//   - OutputType: The resolved output type based on the provided inputs.
func ResolveOutputType(flagOut string, cfg *config.Config) OutputType {
	if flagOut != "" {
		return OutputType(flagOut)
	}
	if kpcliout := os.Getenv("KPASSCLI_OUT"); kpcliout != "" {
		if IsValidType(kpcliout) {
			return OutputType(kpcliout)
		}
	}
	if cfg != nil && cfg.DefaultOutput != "" {
		return OutputType(cfg.DefaultOutput)
	}
	return StdoutType
}
