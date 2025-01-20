package doc

import (
	"fmt"
)

const manPage = `NAME
    kpasscli - KeePass database command line interface

SYNOPSIS
    kpasscli [-kdbpath|-p path] [-kdbpassword|-w path] -item|-i name [-fieldname|-f field] [-out|-o type] [-verify|-v] [-man|-m] [-help|-h]

DESCRIPTION
    kpasscli is a command-line tool for querying KeePass database files.
    It allows retrieving entries and their fields using various search methods.

    The intention for this tool is to use it in automation scripts, to get secret strings like
    cert keys, key passwords, user or tech. user passwords, tokens, ..., which are stored in a keepass-db.
    And it supports optional to open the keepass-db without an interactive password prompt.

    If no -kdbpassword|-w is given, then kpasscli asks for the password to open the keepass-db interactively by a passwored prompt.

    If the item is found, it takes per default the value of the password field or if
    the parameter -fieldname|-f is given, the value of this field.

    Then it depends of the output config, if this is set to
    - stdout: The value is printed to stdout
    - clipboard: The value is copied into the clipboard and can be pasted wherever it is needed

OPTIONS
    -kdbpath|-p path
        Path to the KeePass database file. If not specified, the tool will look for
        the path in the KDBPATH environment variable or the config file.

    -kdbpassword|-w password-file
        Path to a file containing the database password or to an executable that
        outputs the password. For security reasons, the password cannot be provided
        directly on the command line.

    -item|-i name
        The entry to search for. This can be:
        - An absolute path starting with "/" (e.g., "/Personal/Banking/Account")
        - A relative path (e.g., "Banking/Account")
        - A simple name (e.g., "Account")

    -fieldname|-f field
        The field to retrieve from the entry. Defaults to "Password".
        Common fields: Title, UserName, Password, URL, Notes

    -out|-o type
        How to output the retrieved value. Options:
        - stdout: Print to standard output (default)
        - clipboard: Copy to system clipboard

    -case-sensitive|-c
        Enable case-sensitive search

    -exact-match|-e
        Enable exact match search

    -create-config|-cc
        Create an example configuration file

    -man|-m
        Display this manual page

    -help|-h
        Display brief help message

    -verify|-v
        Show the path of found item

    -debug|-d
        Enable debug logging

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
    - database_path:       Default path to the KeePass database
    - default_output:      Default output type (stdout/clipboard)

    # Password retrieval methods, take care, this can be unsecure if you not protect the password file
    # or the executable properly. See SECURITY
    - password_file:       file which contains the password to open the keepass db
    - password_executable: the path to the executable, that returns the password to open the keepass database.
                           This method can be safe, if the executable itself asks for a general password to run it.

ENVIRONMENT
    KPASSCLI_KDBPATH       Alternative way to specify the KeePass database path
    KPASSCLI_OUT           Alternative way to specify the output type (stdout/clipboard)
    KPASSCLI_kdbpassword   Alternative way to specify the password file or executable

SECURITY
    To enable noninteractive access to open the keepass-db, there are two options:
    - provide a password-file
    - provide a executable that prints the password to STDOUT

    In both cases there are security risks, if this is not well prepared.

    A secure way is to use a wallet that is opened with the user login, like kwallet, if you use KDE Desktop.

EXAMPLES
    Get password for a specific entry:
        kpasscli -kdbpath=/path/to/db.kdbx -kdbpassword=/path/to/pass.txt -item="/Personal/Banking/Account"
        kpasscli -p=/path/to/db.kdbx -w=/path/to/pass.txt -i="/Personal/Banking/Account"

    Get username instead of password:
        kpasscli -kdbpath=/path/to/db.kdbx -kdbpassword=/path/to/pass.txt -item="Account" -fieldname=UserName
        kpasscli -p=/path/to/db.kdbx -w=/path/to/pass.txt -i="Account" -f=UserName

    Copy password to clipboard:
        kpasscli -kdbpath=/path/to/db.kdbx -kdbpassword=/path/to/pass.txt -item="Account" -out=clipboard
        kpasscli -p=/path/to/db.kdbx -w=/path/to/pass.txt -i="Account" -o=clipboard

SECURITY
    - Database passwords must be provided via file or executable
    - Clipboard contents are not automatically cleared
    - Be cautious when using clipboard output on shared systems

AUTHOR
    Dirk Jäger

LICENSE
	  GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007
`

func ShowMan() {
	fmt.Print(manPage)
	fmt.Println()
}
