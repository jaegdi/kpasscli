
# kpasscli

A secure command-line interface for KeePass database entries designed for automation, security, and seamless integration with workflows.

kpasscli provides a secure way to query KeePass database entries without exposing passwords in scripts or logs. It's ideal for developers, system administrators, and security-conscious users who need to programmatically access credentials while maintaining strict security standards.

## NAME
kpasscli - KeePass database command line interface for automation

## SYNOPSIS
kpasscli [-kdbpath path] [-kdbpass path] -item name [-fieldname field] [-out type] [-man] [-help]

## DESCRIPTION
kpasscli is a command-line tool for securely retrieving KeePass database entries using various search methods and output configurations. It prioritizes security by:
- Never exposing passwords in command line arguments
- Supporting secure password retrieval via files or executables
- Providing clipboard output without auto-clearing
- Enabling workflow integration without credential leakage

The tool automatically handles password interactions through secure mechanisms (file-based or executable-based) when no interactive password prompt is specified.

## KEY FEATURES
- 🔒 **Security-first design**: Passwords never appear in command line history or process lists
- 🔄 **Flexible search**: Supports absolute paths, relative paths, and simple names
- 🧠 **Smart field selection**: Default to password field or customize with `-fieldname`
- 📦 **Output control**: Print to stdout or copy to clipboard
- ⚙️ **Configurable**: Customizable via environment variables or config files
- 🛡️ **Secure password handling**: Supports password files and secure executables


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

kpasscli uses a layered configuration approach:
1. Environment variables (highest priority)
2. Config file (`~/.config/kpasscli/config.yaml`)
3. Command-line flags

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

### Get password for specific entry:
```bash
kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="/Personal/Banking/Account"
```

### Get username instead of password:
```bash
kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="Account" -fieldname=UserName
```

### Copy password to clipboard:
```bash
kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=/path/to/pass.txt -item="Account" -out=clipboard
```

### Use password executable:
```bash
kpasscli -kdbpath=/path/to/db.kdbx -kdbpass=generate_password.sh -item="Account"
```

## CONFIGURATION EXAMPLE
Create a secure config file:
```yaml
database_path: /home/user/keepsass/db.kdbx
password_file: /home/user/keepsass/pass.txt
default_output: clipboard
```

## ENVIRONMENT VARIABLES
| Variable | Description |
|----------|-------------|
| `KPASSCLI_KDBPATH` | Database path override |
| `KPASSCLI_OUT` | Output type override (stdout/clipboard) |
| `KPASSCLI_KDBPASS` | Password source override |

# SECURITY
- Passwords are **never** exposed in command line arguments
- Database passwords must be provided via only by user readeable file or an executable
- Clipboard contents are automatically cleared after a configurable delay
- Be cautious when using clipboard output on shared systems

# AUTHOR
Dirk Jäger

# LICENSE
GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007
`

