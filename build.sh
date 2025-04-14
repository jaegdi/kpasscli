#!/bin/bash

# echo Generate the config-clusters.go
# build/scripts/generate_config.sh

echo "Build linux binary of kpasscli"
mkdir -p dist/linux-amd64
GOOS=linux GOARCH=amd64 go build -v -o dist/linux-amd64/kpasscli

echo "Build windows binary of kpasscli"
mkdir -p dist/windows-amd64
GOOS=windows GOARCH=amd64 go build -v -o dist/windows-amd64/kpasscli.exe

echo "Build darwin binary of kpasscli"
mkdir -p dist/darwin-amd64
GOOS=darwin GOARCH=amd64 go build -v -o dist/darwin-amd64/kpasscli

echo "Build darwin arm64 binary of kpasscli"
mkdir -p dist/darwin-arm64
GOOS=darwin GOARCH=arm64 go build -v -o dist/darwin-arm64/kpasscli
