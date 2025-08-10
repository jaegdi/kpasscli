package keepass

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/term"

	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/output"
	"kpasscli/src/search"
)

// OpenDatabase opens and decodes a KeePass database file.
// Parameters:
//
//	path: Path to the KeePass database file
//	password: Password to decrypt the database
//
// Returns:
//
//	*gokeepasslib.Database: Decoded database object
//	error: Any error encountered during opening or decoding
func OpenDatabase(path string, password string) (*gokeepasslib.Database, error) {
	file, err := os.Open(path)
	if err != nil {
		debug.Log("Error opening file: %v\n", err)
		return nil, err
	}
	defer file.Close()
	debug.Log("OpenDatabase %s %s", path, strings.Repeat("*", len(password)))

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)
	// debug.Log("OpenDatabase\n%v\n", db)

	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		debug.Log("Error decoding database: %v\n", err)
		return nil, err
	}

	if err := db.UnlockProtectedEntries(); err != nil {
		debug.Log("Error unlocking protected entries: %v\n", err)
		return nil, err
	}

	return db, nil
}

// PasswordPromptFunc defines a function type for prompting the user for a password.
type PasswordPromptFunc func() (string, error)

// ResolvePassword retrieves the database password from a file, executable, or prompt.
// The promptFunc parameter is optional; if nil, getPasswordFromPrompt is used.
// ResolvePassword retrieves the database password from a file or executable.
// It first checks if the password parameter is provided. If not, it falls back to the configuration.
// If the configuration specifies an executable, it runs the executable to get the password.
// If the password parameter is a named pipe, it reads the password from the pipe.
// If the password parameter is a regular file, it reads the password from the file.
// If the password parameter is an executable, it runs the executable to get the password.
// Parameters:
//   - passParam: The password parameter provided via command-line flag.
//   - cfg: The configuration object containing default paths and executables.
//   - kdbpassenv: The environment variable for the password, if set.
//   - promptFunc: Optional function to prompt for password if no file/executable is provided.
//
// Returns:
//   - string: The resolved password.
//   - error: Any error encountered during resolution.
func ResolvePassword(passParam string, cfg *config.Config, kdbpassenv string, promptFunc ...PasswordPromptFunc) (string, error) {

	passfile := ""
	if passParam != "" {
		passfile = passParam
	} else if kdbpassenv != "" {
		passfile = kdbpassenv
	} else if cfg.PasswordFile != "" {
		passfile = cfg.PasswordFile
	} else if cfg.PasswordExecutable != "" {
		passfile = cfg.PasswordExecutable
	} else {
		// Use injected promptFunc if provided, else default
		if len(promptFunc) > 0 && promptFunc[0] != nil {
			return promptFunc[0]()
		}
		return getPasswordFromPrompt()
	}
	// Resolve environment variables in passfile
	// passfile = os.ExpandEnv(passfile)
	debug.Log("PassFile: %v", passfile)

	// Check if passfile is an executable in $PATH
	if execPath, err := exec.LookPath(passfile); err == nil {
		debug.Log("passfile: %v is an executable: %v", passfile, execPath)
		passfile = execPath
	}
	info, err := os.Stat(passfile)
	if err != nil {
		debug.Log("passfile: %v Error: %v", passfile, err.Error())
		return "", fmt.Errorf("password must be provided via file or executable")
	}
	debug.Log("%+v", info)

	if info.Mode()&os.ModeNamedPipe != 0 {
		// Read password from process substitution
		data, err := os.ReadFile(passfile)
		if err != nil {
			debug.Log(err.Error())
			return "", err
		}
		password := strings.TrimSpace(string(data))
		debug.Log("Resolved password from named pipe: %s", strings.Repeat("*", len(password)))
		return password, nil
	}

	if info.Mode()&0111 != 0 {
		// Execute file and read password from stdout
		cmd := exec.Command(passfile)
		output, err := cmd.Output()
		if err != nil {
			debug.Log(err.Error())
			return "", err
		}
		password := strings.TrimSpace(string(output))
		debug.Log("Resolved password from executable: %s", strings.Repeat("*", len(password)))
		return password, nil
	}

	if info.Mode().IsRegular() {
		// Read password from file
		data, err := os.ReadFile(passfile)
		if err != nil {
			debug.Log(err.Error())
			return "", err
		}
		password := strings.TrimSpace(string(data))
		debug.Log("Resolved password from file: %s", strings.Repeat("*", len(password)))
		return password, nil
	}

	return "", fmt.Errorf("password must be provided via file or executable")
}

// getPasswordFromPrompt prompts the user to enter a password securely.
// It reads the password input without echoing it to the terminal, trims any
// leading or trailing whitespace, and returns the password as a string.
// If an error occurs while reading the password, it returns an empty string
// and the encountered error.
// Returns:
//   - string: The entered password, or an empty string if an error occurs.
//   - error: An error if one occurs while reading the password.
func getPasswordFromPrompt() (string, error) {
	// If no valid file or executable is found, prompt the user for the password
	fmt.Print("Enter password: ")
	var password string
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
	}
	password = strings.TrimSpace(string(passwordBytes))
	fmt.Println()
	debug.Log("Resolved password from prompt: %s", strings.Repeat("*", len(password)))
	if password != "" {
		return password, nil
	}
	return "", err
}

// ResolveDatabasePath returns the KeePass database path based on flag, environment, or config.
//
// Parameters:
//   - flagPath: The path provided via command-line flag.
//   - cfg: The configuration object containing a default database path.
//
// Returns:
//   - string: The resolved database path, or empty string if not found.
func ResolveDatabasePath(flagPath string, cfg *config.Config) string {
	if flagPath != "" {
		return flagPath
	}
	if envPath := os.Getenv("KPASSCLI_KDBPATH"); envPath != "" {
		return envPath
	}
	if cfg != nil && cfg.DatabasePath != "" {
		return cfg.DatabasePath
	}
	return ""
}

// GetAllFields finds a specific entry by path and displays all its fields.
//
// Parameters:
//   - db: The KeePass database to search.
//   - config: The configuration object for output formatting.
//   - itemPath: The path of the entry to display.
//
// Returns:
//   - error: Any error encountered during the operation.
//
// GetAllFields finds a specific entry by path and displays all its fields.
//
// Parameters:
//   - db: The KeePass database to search.
//   - config: The configuration object for output formatting.
//   - itemPath: The path of the entry to display.
//
// Returns:
//   - error: Any error encountered during the operation.
func GetAllFields(db *gokeepasslib.Database, config *config.Config, itemPath string) error {
	return GetAllFieldsWithFinder(db, config, itemPath, nil, nil)
}

// GetAllFieldsWithFinder is like GetAllFields but allows dependency injection for testing.
// If finder is nil, uses search.NewFinder(db). If showAllFields is nil, uses output.ShowAllFields.
func GetAllFieldsWithFinder(
	db *gokeepasslib.Database,
	config *config.Config,
	itemPath string,
	finder search.FinderInterface,
	showAllFields func(*gokeepasslib.Entry, config.Config),
) error {
	if finder == nil {
		finder = search.NewFinder(db)
	}
	results, err := finder.Find(itemPath)
	if err != nil {
		return fmt.Errorf("error finding entry '%s': %w", itemPath, err)
	}
	if len(results) == 0 {
		return fmt.Errorf("entry not found: %s", itemPath)
	}
	if len(results) > 1 {
		var foundPaths []string
		for _, res := range results {
			foundPaths = append(foundPaths, res.Path)
		}
		return fmt.Errorf("multiple entries found for '%s', please specify a unique path: %s", itemPath, strings.Join(foundPaths, ", "))
	}
	singleEntry := results[0].Entry
	if singleEntry == nil {
		return fmt.Errorf("found result for '%s', but entry data is unexpectedly nil", itemPath)
	}
	if showAllFields == nil {
		showAllFields = output.ShowAllFields
	}
	showAllFields(singleEntry, *config)
	return nil
}
