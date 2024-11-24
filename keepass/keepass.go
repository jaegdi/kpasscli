package keepass

import (
	"fmt"
	"kpasscli/config"
	"kpasscli/debug"
	"os"
	"os/exec"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
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
// Parameters:
//
//	passParam: Path to password file or executable
//
// Returns:
//
//	string: The resolved password
//	error: Any error encountered during password retrieval
func ResolvePassword(passParam string, cfg *config.Config) (string, error) {
	if passParam == "" && cfg.PasswordFile == "" && cfg.PasswordExecutable == "" {
		return "", fmt.Errorf("password parameter -kdbpass is required")
	}

	debug.Log("ResolvePassword %v", cfg)
	if cfg.PasswordExecutable != "" && passParam == "" {
		cmd := exec.Command(cfg.PasswordExecutable)
		debug.Log("cmd %v", cmd)
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		password := strings.TrimSpace(string(output))
		debug.Log("Resolved password from executable: %s", strings.Repeat("*", len(password)))
		return password, nil
	}

	info, err := os.Stat(passParam)
	if err != nil {
		return "", fmt.Errorf("password must be provided via file or executable")
	}

	debug.Log("%v", info)
	if info.Mode()&os.ModeNamedPipe != 0 {
		// Read password from process substitution
		data, err := os.ReadFile(passParam)
		if err != nil {
			return "", err
		}
		password := strings.TrimSpace(string(data))
		debug.Log("Resolved password from named pipe: %s", strings.Repeat("*", len(password)))
		return password, nil
	}
	if info.Mode().IsRegular() {
		// Read password from file
		data, err := os.ReadFile(passParam)
		if err != nil {
			return "", err
		}
		password := strings.TrimSpace(string(data))
		debug.Log("Resolved password from file: %s", strings.Repeat("*", len(password)))
		return password, nil
	}

	if info.Mode()&0111 != 0 {
		// Execute file and read password from stdout
		cmd := exec.Command(passParam)
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		password := strings.TrimSpace(string(output))
		debug.Log("Resolved password from executable: %s", strings.Repeat("*", len(password)))
		return password, nil
	}

	return "", fmt.Errorf("password must be provided via file or executable")
}
