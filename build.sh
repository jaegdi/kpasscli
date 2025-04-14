#!/bin/bash

echo "Build linux binary of kpasscli"
go build -v -o dist/kpasscli

echo "Build windows binary of kpasscli"
GOOS=windows GOARCH=amd64 go build -v -o dist/kpasscli.exe

exit
