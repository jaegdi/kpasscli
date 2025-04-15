#Requires -Version 5.1 # Optional: Specify minimum PowerShell version if needed
<#
.SYNOPSIS
  Builds the kpasscli Go project for multiple target platforms (Linux, Windows, macOS).

.DESCRIPTION
  This script cross-compiles the kpasscli application using the Go toolchain.
  It sets the GOOS and GOARCH environment variables for each target build
  and places the resulting binaries in corresponding subdirectories under 'dist'.

.NOTES
  Author: Based on original psbuild.sh by jaegdi/Dirk
  Requires: Go toolchain (go command) installed and in PATH.
#>

# echo Generate the config-clusters.go
# build/scripts/generate_config.sh

# --- Build Linux ---
Write-Host ""
Write-Host "Build linux binary of kpasscli (amd64)"
# Equivalent of mkdir -p
New-Item -ItemType Directory -Path 'dist/linux-amd64' -Force | Out-Null
# Set environment variables for the 'go build' command
$env:GOOS = "linux"
$env:GOARCH = "amd64"
# Execute the build
go build -v -o dist/linux-amd64/kpasscli

# --- Build Windows ---
Write-Host ""
Write-Host "Build windows binary of kpasscli (amd64)"
New-Item -ItemType Directory -Path 'dist/windows-amd64' -Force | Out-Null
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -v -o dist/windows-amd64/kpasscli.exe # Note the .exe extension

# --- Build Darwin (macOS Intel) ---
Write-Host ""
Write-Host "Build darwin binary of kpasscli (amd64)"
New-Item -ItemType Directory -Path 'dist/darwin-amd64' -Force | Out-Null
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -v -o dist/darwin-amd64/kpasscli

# --- Build Darwin (macOS Apple Silicon) ---
Write-Host ""
Write-Host "Build darwin arm64 binary of kpasscli (arm64)"
New-Item -ItemType Directory -Path 'dist/darwin-arm64' -Force | Out-Null
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -v -o dist/darwin-arm64/kpasscli

# Optional: Clear the environment variables if the script continues and needs the original values
$env:GOOS = $null
$env:GOARCH = $null

Write-Host ""
Write-Host "Build process finished."
Write-Host ""
