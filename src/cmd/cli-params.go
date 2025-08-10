package cmd

import (
	"flag"

	"kpasscli/src/doc"
)

// parses command-line flags, and supports the following flags:
//
// -kdbpath | -p: Path to KeePass database file
// -kdbpassword | -w: Password file or executable to get password
// -item | -i: Item to search for
// -fieldname | -f: Field name to retrieve (default: "Password")
// -out | -o: Output type (clipboard/stdout)
// -case-sensitive | -c: Enable case-sensitive search
// -exact-match | -e: Enable exact match search
// -man | -m: Show manual page
// -help | -h: Show help message
// -debug | -d: Enable debug logging
// -config | -c path: Path to configuration file (default: ~/.config/kpasscli/config.yaml)
// -verify | -v: Enable verify messages
// -create-config | -cc: Create an example config file
// -print-config | -pc: print the current detected config to stdout

type Flags struct {
	KdbPath       string
	KdbPassword   string
	Item          string
	FieldName     string
	Out           string
	ConfigPath    string
	ClearAfter    int
	CaseSensitive bool
	ExactMatch    bool
	ShowMan       bool
	ShowHelp      bool
	DebugFlag     bool
	VerifyFlag    bool
	CreateConfig  bool
	PrintConfig   bool
	ShowAll       bool
}

// ParseFlags parses flags from the provided FlagSet and args. For production, use ParseFlagsDefault().
func ParseFlags(fs *flag.FlagSet, args []string) *Flags {
	flags := &Flags{}
	fs.StringVar(&flags.KdbPath, "kdbpath", "", "Path to KeePass database file")
	fs.StringVar(&flags.KdbPath, "p", "", "Path to KeePass database file (shorthand)")

	fs.StringVar(&flags.KdbPassword, "kdbpassword", "", "Password file or executable to get password")
	fs.StringVar(&flags.KdbPassword, "w", "", "Password file or executable to get password (shorthand)")

	fs.StringVar(&flags.Item, "item", "", "Item to search for")
	fs.StringVar(&flags.Item, "i", "", "Item to search for (shorthand)")

	fs.StringVar(&flags.FieldName, "fieldname", "Password", "Field name to retrieve")
	fs.StringVar(&flags.FieldName, "f", "Password", "Field name to retrieve (shorthand)")

	fs.StringVar(&flags.Out, "out", "", "Output type (clipboard/stdout)")
	fs.StringVar(&flags.Out, "o", "", "Output type (clipboard/stdout) (shorthand)")

	fs.BoolVar(&flags.CaseSensitive, "case-sensitive", false, "Enable case-sensitive search")
	fs.BoolVar(&flags.CaseSensitive, "cs", false, "Enable case-sensitive search (shorthand)")

	fs.BoolVar(&flags.ExactMatch, "exact-match", false, "Enable exact match search")
	fs.BoolVar(&flags.ExactMatch, "e", false, "Enable exact match search (shorthand)")

	fs.BoolVar(&flags.ShowMan, "man", false, "Show manual page")
	fs.BoolVar(&flags.ShowMan, "m", false, "Show manual page (shorthand)")

	fs.BoolVar(&flags.ShowHelp, "help", false, "Show help message")
	fs.BoolVar(&flags.ShowHelp, "h", false, "Show help message (shorthand)")

	fs.BoolVar(&flags.VerifyFlag, "verify", false, "Enable verify message")
	fs.BoolVar(&flags.VerifyFlag, "v", false, "Enable verify message (shorthand)")

	fs.BoolVar(&flags.DebugFlag, "debug", false, "Enable debug logging")
	fs.BoolVar(&flags.DebugFlag, "d", false, "Enable debug logging (shorthand)")

	fs.BoolVar(&flags.CreateConfig, "create-config", false, "Create an example config file")
	fs.BoolVar(&flags.CreateConfig, "cc", false, "Create an example config file (shorthand)")

	fs.BoolVar(&flags.PrintConfig, "print-config", false, "Print current configuration")
	fs.BoolVar(&flags.PrintConfig, "pc", false, "Print current configuration (shorthand)")

	fs.StringVar(&flags.ConfigPath, "config", "~/.config/kpasscli/config.yaml", "Path to configuration file")
	fs.StringVar(&flags.ConfigPath, "c", "~/.config/kpasscli/config.yaml", "Path to configuration file (shorthand)")

	fs.Usage = doc.ShowHelp
	fs.Parse(args) // Parse the flags from the provided args. This is implemented to test the ParseFlags function.
	return flags
}

// ParseFlagsDefault parses flags from the global flag.CommandLine and os.Args[1:].
// This is the default function to use in production.
// It sets up the flags and returns the parsed Flags struct.
// It also sets the usage function to show help documentation.
// It is typically called in the main function of the application.
// It initializes the logging, handles special flags, and exits the program if necessary.
func ParseFlagsDefault() *Flags {
	return ParseFlags(flag.CommandLine, nil)
}
