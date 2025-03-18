#!/bin/bash
set -eo pipefail



# build the linux binary of kpasscli
echo "Build linux binary of kpasscli"
go build -v

# build the windows binary of kpasscli
echo "Build windows binary of kpasscli"
GOOS=windows GOARCH=amd64 go build -v

# check kpasscli and if it works, then upload to artifactory
if ./kpasscli -c config.yaml -i pw1.1 -f username; then
    echo "Push to artifactory"

    artifactory-upload.sh -lf=kpasscli       -tr=scptools-bin-dev-local   -tf=/tools/kpasscli/kpasscli-linux/
    artifactory-upload.sh -lf=kpasscli       -tr=scpas-bin-dev-local      -tf=/istag_and_image_management/kpasscli-linux/

    artifactory-upload.sh -lf=kpasscli.exe   -tr=scptools-bin-dev-local   -tf=/tools/kpassclis/kpasscli-windows/
    artifactory-upload.sh -lf=kpasscli.exe   -tr=scpas-bin-dev-local      -tf=/istag_and_image_management/kpasscli-windows/

    echo "Copy it to share folder PEWI4124://Daten"
    cp kpasscli kpasscli.exe  /gast-drive-d/Daten/
fi

echo
echo
echo "#  B U I L D   I M A G E   T O O L   F O R   U B I 7"
BINARY_NAME="kpasscli"
BINARY_NAME_UBI7="dist/${BINARY_NAME}-ubi7"
IMAGE="${BINARY_NAME}:ubi7"
CONTAINER_NAME="${BINARY_NAME}-container"

# build ubi7 binary in image
/usr/bin/podman build -t "$IMAGE" -f Dockerfile .

echo "##########  copy binary from container to local  ##########"
if podman ps -a | rg "$CONTAINER_NAME" >/dev/null; then
    podman rm "$CONTAINER_NAME"
fi
podman create --name "$CONTAINER_NAME" "localhost/$IMAGE"
podman cp "$CONTAINER_NAME":/app/dist/kpasscli "$BINARY_NAME_UBI7"
scp "$BINARY_NAME_UBI7" cid-scp0-tls-v01-mgmt:
podman rm "$CONTAINER_NAME"

artifactory-upload.sh -lf="$BINARY_NAME_UBI7" -tr=scptools-bin-dev-local  -tf=ocp-stable-4.16/clients/$BINARY_NAME

