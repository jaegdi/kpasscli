#!/bin/bash
set -e

[ -f kpasscli ] && cp kpasscli "old-version/kpasscli_$(date +%F_%T)"
[ -f kpasscli.exe ] && cp kpasscli.exe "old-version/kpasscli.exe_$(date +%F_%T).exe"

# echo Generate the config-clusters.go
# build/scripts/generate_config.sh

echo "Build linux binary of kpasscli"
go build -v

echo "Build windows binary of kpasscli"
GOOS=windows GOARCH=amd64 go build -v

if ./kpasscli -h > /dev/null; then
    echo "Push to artifactory"

    artifactory-upload.sh -lf=kpasscli       -tr=scptools-bin-dev-local   -tf=/tools/kpasscli/kpasscli-linux/ || echo "Upload to artifactory failed"
    artifactory-upload.sh -lf=kpasscli       -tr=scpas-bin-dev-local      -tf=/kpasscli/kpasscli-linux/ || echo "Upload to artifactory failed"

    artifactory-upload.sh -lf=kpasscli.exe   -tr=scptools-bin-dev-local   -tf=/tools/kpasscli/kpasscli-windows/ || echo "Upload to artifactory failed"
    artifactory-upload.sh -lf=kpasscli.exe   -tr=scpas-bin-dev-local      -tf=/kpasscli/kpasscli-windows/ || echo "Upload to artifactory failed"

    # jf rt u --server-id saas --flat kpasscli  /scptools-bin-develop/tools/kpasscli/kpasscli-linux/kpasscli
    # jf rt u --server-id saas --flat kpasscli  /scpas-bin-develop/kpasscli/kpasscli-linux/kpasscli
    # jf rt u --server-id saas --flat kpasscli.exe  /scptools-bin-develop/tools/kpasscli/kpasscli-windows/kpasscli.exe
    # jf rt u --server-id saas --flat kpasscli.exe  /scpas-bin-develop/kpasscli/kpasscli-windows/kpasscli.exe

    echo "Copy it to share folder PEWI4124://Daten"
    cp kpasscli kpasscli.exe  /gast-drive-d/Daten/
fi
