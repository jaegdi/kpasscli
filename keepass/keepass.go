package keepass

import (
	"fmt"
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
		return nil, err
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	decoder := gokeepasslib.NewDecoder(file)
	if err := decoder.Decode(db); err != nil {
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
func ResolvePassword(passParam string) (string, error) {
	if passParam == "" {
		return "", fmt.Errorf("password parameter is required")
	}

	info, err := os.Stat(passParam)
	if err != nil {
		return "", fmt.Errorf("password must be provided via file or executable")
	}

	if info.Mode().IsRegular() {
		// Read password from file
		data, err := os.ReadFile(passParam)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}

	if info.Mode()&0111 != 0 {
		// Execute file and read password from stdout
		cmd := exec.Command(passParam)
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}

	return "", fmt.Errorf("password must be provided via file or executable")
}
