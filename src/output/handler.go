package output

import (
	"encoding/json"
	"fmt"
	// Hinzugefügt für Debug-Logs

	"golang.design/x/clipboard"

	"kpasscli/src/debug"
)

type OutputChannel string

const (
	Clipboard  OutputChannel = "clipboard"
	Stdout     OutputChannel = "stdout"
	lineBreak                = "----------------------------------------"
	timeFormat               = "2006-01-02 15:04:05"
)

<<<<<<< HEAD

// Handler is an interface for outputting values.
type Handler interface {
	Output(string) error
}

// stdHandler is the default implementation of Handler.
type stdHandler struct {
	outputType Type
=======
type Handler struct {
	outputChannel OutputChannel
>>>>>>> a1d926327a8839ef184d0f5f955fa542c6b2e174
}

// NewHandler creates a new Handler instance with the specified output type.
// Parameters:
//   - outputChannel: The type of output (clipboard or stdout).
//
// Returns:
//   - *Handler: A new Handler instance.
<<<<<<< HEAD
// NewHandler creates a new Handler instance with the specified output type.
func NewHandler(outputType Type) Handler {
	return &stdHandler{outputType: outputType}
=======
func NewHandler(outputChannel OutputChannel) *Handler {
	return &Handler{outputChannel: outputChannel}
>>>>>>> a1d926327a8839ef184d0f5f955fa542c6b2e174
}

// Output outputs the given value based on the handler's output type.
// Parameters:
//   - value: The value to be output.
//
// Returns:
//   - error: Any error encountered during output.
<<<<<<< HEAD
func (h *stdHandler) Output(value string) error {
	debug.Log("Outputting value: %s", value)
	switch h.outputType {
=======
func (h *Handler) Output(value string) error {
	debug.Log("Outputting value: %s", value) // Debug-Log hinzugefügt
	switch h.outputChannel {
>>>>>>> a1d926327a8839ef184d0f5f955fa542c6b2e174
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

// ShowAllFields displays all fields of a KeePass entry
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

func formatTime(t time.Time) string {
	return t.Format(timeFormat)
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
