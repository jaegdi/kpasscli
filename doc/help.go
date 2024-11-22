package doc

import (
	"fmt"
)

const manPage = `NAME
    kpasscli - KeePass database command line interface

SYNOPSIS
    kpasscli [-kdbpath path] [-kdbpass path] -item name [-fieldname field] [-out type] [-man] [-help]

DESCRIPTION
    kpasscli is a command-line tool for querying KeePass database files.
    It allows retrieving entries and their fields using various search methods.

OPTIONS
    -kdbpath path
        Path to the KeePass database file. If not specified, the tool will look for
        the path in the KDBPATH environment variable or the config file.

    -kdbpass path
        Path to a file containing the database password or to an executable that
        outputs the password. For security reasons, the password cannot be provided
        directly on the command line.

    -item name
        The entry to search for. This can be:
        - An absolute path starting with "/" (e.g., "/Personal/Banking/Account")
        - A relative path (e.g., "Banking/Account")
        - A simple name (e.g., "Account")

    -fieldname field
        The field to retrieve from the entry. Defaults to "Password".
        Common fields: Title, UserName, Password, URL, Notes

    -out type
        How to output the retrieved value. Options:
        - stdout: Print to standard output (default)
        - clipboard: Copy to system clipboard

    -man
        Display this manual page

    -help
        Display brief help message

SEARCH BEHAVIOR
    Absolute Path (/path/to/entry):
        Searches for an exact match at the specified location in the database.

    Relative Path (path/to/entry):
        Searches through all groups for a matching path.
        If multiple matches are found, lists all matches.

    Simple Name (entry):
        Searches all entries regardless of location.
        If multiple matches are found, lists all matches.

CONFIGURATION
    Configuration can be provided via a config.yaml file with the following fields:
    - database_path: Default path to the KeePass database
    - default_output: Default output type (stdout/clipboard)

ENVIRONMENT
    KDBPATH
        Alternative way to specify the KeePass database path

EXAMPLES
    Get password for a specific entry:
        kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="/Personal/Banking/Account"

    Get username instead of password:
        kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="Account" -fieldname=UserName

    Copy password to clipboard:
        kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="Account" -out=clipboard

SECURITY
    - Database passwords must be provided via file or executable
    - Clipboard contents are not automatically cleared
    - Be cautious when using clipboard output on shared systems

AUTHOR
    [Your Name or Organization]

VERSION
    1.0.0
`

func ShowMan() {
    fmt.Println(manPage)
}

func ShowHelp() {
    help := `Usage: kpasscli [OPTIONS]

Options:
    -kdbpath path    Path to KeePass database file
    -kdbpass path    Path to password file or executable
    -item name       Entry to search for
    -fieldname field Field to retrieve (default: Password)
    -out type        Output type (stdout/clipboard)
    -man            Show full manual
    -help           Show this help

Example:
    kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="/Personal/Banking/Account"

For more information, use -man`

    fmt.Println(help)
}
