
# kpasscli

A keepass client to query item values from a keepass db by cli.


## NAME
kpasscli - KeePass database command line interface for automation

## SYNOPSIS
kpasscli [-kdbpath path] [-kdbpass path] -item name [-fieldname field] [-out type] [-man] [-help]

## DESCRIPTION
kpasscli is a command-line tool for querying KeePass database files.
It allows retrieving entries and their fields using various search methods.

If no -kdbpass is given, then kpasscli asks for the password interactively by a passwored prompt.

If the item is found, it takes per default the value of the password field or if the parameter -fieldname is given, the value of this field.

Then it depends of the output config, if this is set to
- stdout: The value is printed to stdout
- clipboard: The value is copied into the clipboard and can be pasted wherever it is needed

## OPTIONS

###    -kdbpath path  or envvar KPASSCLI_KDBPATH  or config file: database_path
Path to the KeePass database file. If not specified, the tool will look for
the path in the KDBPATH environment variable or the config file.

###    -kdbpass path  or envvar KPASSCLI_KDBPASS  or config file: password_executable|password_file
Path to a file containing the database password or to an executable that outputs the password. For security reasons, the password cannot be provided directly on the command line.

###    -item name
The entry to search for. This can be:
- An absolute path starting with "/" (e.g., "/MY_KP_ROOT/Personal/Banking/Account")
- A relative path (e.g., "Banking/Account")
- A simple name (e.g., "Account")

###    -fieldname field
The field to retrieve from the entry. Defaults to "Password".
Common fields: Title, UserName, Password, URL, Notes

###    -out type   or envvar KPASSCLI_OUT  or config file: default_output
How to output the retrieved value. Options:
- stdout: Print to standard output (default)
- clipboard: Copy to system clipboard

###    -createConfig
Create an example configuration file

###    -man
Display this manual page

###    -help
Display brief help message

## SEARCH BEHAVIOR
###    Absolute Path
eg. **-item=/root/subpath/subpath/to/entry**

Searches for an exact match at the specified location in the database.
It returns the value of the item, per default the password or if the -fieldname parameter is given, the value of this field.

###    Relative Path
eg. **-item=subpath/to/entry**

Searches through all groups in the keepass-db for a matching subpath with the entry.
If multiple matches are found, returns with error and lists all matches.

Otherwise it returns the value of the item, per default the password, or if the -fieldname parameter is given, the value of this field.

###    Simple Name (entry):
eg.  **-item=entry**

Searches all matching entries regardless of location.
If multiple matches are found, returns with error and lists all matches.

Otherwise it returns the value of the item, per default the password or, if the -fieldname parameter is given, the value of this field.

## CONFIGURATION

Configuration can be provided via a config.yaml file with the following fields:
- **database_path**:       Default path to the KeePass database
- **default_output**:      Default output type (stdout/clipboard)
- **password_file**:       file which contains the password to open the keepass db
- **password_executable**: the path to the executable, that returns the password to open the keepass database.
This method can be safe, if the executable itself asks for a general password to run it.
## Password retrieval methods
take care, this can be unsecure if you not protect the password file
or the executable properly

To create a example config file, kpasscli can be executed with parameter  -createConfig


## ENVIRONMENT VARIABLES
###    KPASSCLI_KDBPATH
Alternative way to specify the KeePass database path
###    KPASSCLI_OUT
Alternative way to specify the output type (stdout/clipboard)
###    KPASSCLI_KDBPASS
Alternative way to specify the password file or executable

## EXAMPLES

    # Get password for a specific entry:
    kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="/Personal/Banking/Account"

    # Get username instead of password:
    kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="Account" -fieldname=UserName

    # Copy password to clipboard:
    kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="Account" -out=clipboard
    # or with short flags
    kpasscli -p /path/to/db.kdbx -w /path/to/pass.txt -i Acount -o clipboard

    # Provide the password to open the keepass db by a executable
    kpasscli -kdbpath=/path/to/db.kdbx -kdbpass  <(script_to_get_keepass-db_password) -item Account
    # but also works without process substitution, kpasscli checks if the given password-file is a executable
    # and execute it automatically
    kpasscli -kdbpath=/path/to/db.kdbx -kdbpass script_to_get_keepass-db_password -item Account

# SECURITY
- Database passwords must be provided via file or executable
- Clipboard contents are not automatically cleared
- Be cautious when using clipboard output on shared systems

# AUTHOR
Dirk Jäger

# LICENSE
GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007
`

