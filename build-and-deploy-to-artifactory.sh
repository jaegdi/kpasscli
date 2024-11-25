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

    artifactory-upload.sh -lf=kpasscli       -tr=scptools-bin-dev-local   -tf=/tools/kpasscli/kpasscli-linux/
    artifactory-upload.sh -lf=kpasscli       -tr=scpas-bin-dev-local      -tf=/istag_and_image_management/kpasscli-linux/

    artifactory-upload.sh -lf=kpasscli.exe   -tr=scptools-bin-dev-local   -tf=/tools/kpasscli/kpasscli-windows/
    artifactory-upload.sh -lf=kpasscli.exe   -tr=scpas-bin-dev-local      -tf=/istag_and_image_management/kpasscli-windows/

    # jf rt u --server-id default --flat kpasscli  /scptools-bin-develop/tools/kpasscli/kpasscli-linux/kpasscli
    # jf rt u --server-id default --flat kpasscli  /scpas-bin-develop/istag_and_image_management/kpasscli-linux/kpasscli
    # jf rt u --server-id default --flat kpasscli.exe  /scptools-bin-develop/tools/kpasscli/kpasscli-windows/kpasscli.exe
    # jf rt u --server-id default --flat kpasscli.exe  /scpas-bin-develop/istag_and_image_management/kpasscli-windows/kpasscli.exe

    echo "Copy it to share folder PEWI4124://Daten"
    cp kpasscli kpasscli.exe  /gast-drive-d/Daten/
fi
