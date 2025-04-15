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
// -create-config | -cc: Create an example config file
// -print-config | -pc: print the current detected config to stdout

type Flags struct {
	KdbPath       string
	KdbPassword   string
	Item          string
	FieldName     string
	Out           string
	ConfigPath    string
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

func ParseFlags() *Flags {
	flags := &Flags{}

	// Define flags with both long and short versions
	flag.StringVar(&flags.KdbPath, "kdbpath", "", "Path to KeePass database file")
	flag.StringVar(&flags.KdbPath, "p", "", "Path to KeePass database file (shorthand)")

	flag.StringVar(&flags.KdbPassword, "kdbpassword", "", "Password file or executable to get password")
	flag.StringVar(&flags.KdbPassword, "w", "", "Password file or executable to get password (shorthand)")

	flag.StringVar(&flags.ConfigPath, "config", "~/.config/kpasscli/config.yaml", "Path to configuration file")
	flag.StringVar(&flags.ConfigPath, "c", "~/.config/kpasscli/config.yaml", "Path to configuration file (shorthand)")

	flag.StringVar(&flags.Item, "item", "", "Item to search for")
	flag.StringVar(&flags.Item, "i", "", "Item to search for (shorthand)")

	flag.StringVar(&flags.FieldName, "fieldname", "Password", "Field name to retrieve")
	flag.StringVar(&flags.FieldName, "f", "Password", "Field name to retrieve (shorthand)")

	flag.StringVar(&flags.Out, "out", "", "Output type (clipboard/stdout)")
	flag.StringVar(&flags.Out, "o", "", "Output type (clipboard/stdout) (shorthand)")

	flag.BoolVar(&flags.ShowAll, "show-all", false, "Show all entries of an item.")
	flag.BoolVar(&flags.ShowAll, "a", false, "Show all entries of an item. (shorthand)")

	flag.BoolVar(&flags.CaseSensitive, "case-sensitive", false, "Enable case-sensitive search")
	flag.BoolVar(&flags.CaseSensitive, "cs", false, "Enable case-sensitive search (shorthand)")

	flag.BoolVar(&flags.ExactMatch, "exact-match", false, "Enable exact match search")
	flag.BoolVar(&flags.ExactMatch, "e", false, "Enable exact match search (shorthand)")

	flag.BoolVar(&flags.ShowMan, "man", false, "Show manual page")
	flag.BoolVar(&flags.ShowMan, "m", false, "Show manual page (shorthand)")

	flag.BoolVar(&flags.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&flags.ShowHelp, "h", false, "Show help message (shorthand)")

	flag.BoolVar(&flags.VerifyFlag, "verify", false, "Enable verify message")
	flag.BoolVar(&flags.VerifyFlag, "v", false, "Enable verify message (shorthand)")

	flag.BoolVar(&flags.DebugFlag, "debug", false, "Enable debug logging")
	flag.BoolVar(&flags.DebugFlag, "d", false, "Enable debug logging (shorthand)")

	flag.BoolVar(&flags.CreateConfig, "create-config", false, "Create an example config file")
	flag.BoolVar(&flags.CreateConfig, "cc", false, "Create an example config file (shorthand)")

	flag.BoolVar(&flags.PrintConfig, "print-config", false, "Print current configuration")
	flag.BoolVar(&flags.PrintConfig, "pc", false, "Print current configuration (shorthand)")

	flag.Usage = doc.ShowHelp
	flag.Parse()

	return flags
}
