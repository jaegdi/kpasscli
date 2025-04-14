package keepass

import (
	"fmt"
	"kpasscli/src/config"
	"kpasscli/src/debug"
	"kpasscli/src/output"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/crypto/ssh/terminal"
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

// ResolvePassword retrieves the database password from a file or executable.
// It first checks if the password parameter is provided. If not, it falls back to the configuration.
// If the configuration specifies an executable, it runs the executable to get the password.
// If the password parameter is a named pipe, it reads the password from the pipe.
// If the password parameter is a regular file, it reads the password from the file.
// If the password parameter is an executable, it runs the executable to get the password.
// Parameters:
//   - passParam: Path to password file or executable
//   - cfg: Configuration object containing password file or executable paths
//
// Returns:
//   - string: The resolved password
//   - error: Any error encountered during password retrieval
func ResolvePassword(passParam string, cfg *config.Config, kdbpassenv string) (string, error) {

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
		return getPasswordFromPrompt()
	}
	// Resolve environment variables in passfile
	passfile = os.ExpandEnv(passfile)
	debug.Log(passfile)

	// Check if passfile is an executable in $PATH
	if execPath, err := exec.LookPath(passfile); err == nil {
		passfile = execPath
	}
	info, err := os.Stat(passfile)
	if err != nil {
		debug.Log("passfile: %v Error: %v", passfile, err.Error())
		return "", fmt.Errorf("password must be provided via file or executable")
	}
	debug.Log("%v", info)

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
func getPasswordFromPrompt() (string, error) {
	// If no valid file or executable is found, prompt the user for the password
	fmt.Print("Enter password: ")
	var password string
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
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

// ResolveDatabasePath determines the path to the KeePass database file.
// It checks the following sources in order of precedence:
// 1. The provided flagPath argument. If it is not empty, it is returned.
// 2. The environment variable "KPASSCLI_KDBPATH". If it is set, its value is returned.
// 3. The DatabasePath field in the provided config.Config object. If it is not nil and its DatabasePath field is not empty, it is returned.
// If none of these sources provide a path, an empty string is returned.
//
// Parameters:
// - flagPath: A string representing the path provided via a command-line flag.
// - cfg: A pointer to a config.Config object that may contain the database path.
//
// Returns:
// A string representing the resolved database path, or an empty string if no path is found.
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

func GetAllFields(db *gokeepasslib.Database, itemPath string) error {
	entry, err := search.FindEntry(db, itemPath)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("entry not found: %s", itemPath)
	}

	output.ShowAllFields(entry)
	return nil
}
