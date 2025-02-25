package doc

import (
	"fmt"
)

func ShowHelp() {
	help := `Usage: kpasscli [OPTIONS]

Options:
    -kdbpath | -p path      Path to KeePass database file
    -kdbpassword | -w path  Path to password file or executable, if not given asks for password interactively
    -item | -i name         Entry to search for
    -fieldname | -f field   Field to retrieve (default: Password)
    -out | -o type          Output type (stdout/clipboard)
    -case-sensitive | -c    Enable case-sensitive search
    -exact-match | -e       Enable exact match search
    -create-config | -cc    Create an example config file
    -verify | -v            Show the path of found item
    -debug | -d             Enable debug logging
    -man | -m               Show full manual
    -help | -h              Show this help

Example:
    kpasscli -kdbpath=/path/to/db.kdbx -kdbpassword=/path/to/pass.txt -item="/Personal/Banking/Account"
    kpasscli -p=/path/to/db.kdbx -w=/path/to/pass.txt -i="/Personal/Banking/Account"

    if keepass-db file and password-file|password-exec and output type is set in the config file
    then it's enough to specify the item and my be the fieldname.

    # for password
    kpasscli -i=/Personal/Banking/Account

    # or if Account is uniq in the keepass-db
    kasscli -i=Account

    # To verify, if the right item was found, you can use the -verify flag
    kpasscli -i=Account -v

    # for username
    kasscli -i=/Personal/Banking/Account -f=UserName

For more information, use -man|-m

AUTHOR
	Dirk Jäger

LICENSE
	GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007`

	fmt.Println(help)
}
