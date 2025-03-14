#!/bin/bash
set -e

[ -f dist/kpasscli ] && cp dist/kpasscli "old-version/kpasscli_$(date +%F_%T)"
[ -f dist/kpasscli.exe ] && cp dist/kpasscli.exe "old-version/kpasscli.exe_$(date +%F_%T).exe"

# echo Generate the config-clusters.go
# build/scripts/generate_config.sh

echo "Build linux binary of kpasscli"
go build -v -o dist/kpasscli

echo "Build windows binary of kpasscli"
GOOS=windows GOARCH=amd64 go build -v -o dist/kpasscli.exe

if dist/kpasscli -h > /dev/null; then
    echo "Push to artifactory"

    artifactory-upload.sh -lf=dist/kpasscli       -tr=scptools-bin-dev-local   -tf=/tools/kpasscli/kpasscli-linux/ || echo "Upload to artifactory failed"
    artifactory-upload.sh -lf=dist/kpasscli       -tr=scpas-bin-dev-local      -tf=/kpasscli/kpasscli-linux/ || echo "Upload to artifactory failed"

    artifactory-upload.sh -lf=dist/kpasscli.exe   -tr=scptools-bin-dev-local   -tf=/tools/kpasscli/kpasscli-windows/ || echo "Upload to artifactory failed"
    artifactory-upload.sh -lf=dist/kpasscli.exe   -tr=scpas-bin-dev-local      -tf=/kpasscli/kpasscli-windows/ || echo "Upload to artifactory failed"

    # jf rt u --server-id saas --flat kpasscli  /scptools-bin-develop/tools/kpasscli/kpasscli-linux/kpasscli
    # jf rt u --server-id saas --flat kpasscli  /scpas-bin-develop/kpasscli/kpasscli-linux/kpasscli
    # jf rt u --server-id saas --flat kpasscli.exe  /scptools-bin-develop/tools/kpasscli/kpasscli-windows/kpasscli.exe
    # jf rt u --server-id saas --flat kpasscli.exe  /scpas-bin-develop/kpasscli/kpasscli-windows/kpasscli.exe

    echo "Copy it to share folder PEWI4124://Daten"
    cp dist/kpasscli dist/kpasscli.exe  /gast-drive-d/Daten/
fi
