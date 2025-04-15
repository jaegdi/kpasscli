package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func runKpasscli(args ...string) (string, error) {
	ostyp := "lin"
	distdir := "./dist/linux-amd64"
	os := runtime.GOOS
	arch := runtime.GOARCH
	if arch == "arm64" {
		distdir = "./dist/" + os + "-" + arch
	} else if runtime.GOARCH == "arm" {
		distdir = "./dist/" + os + "-" + arch
	} else {
		distdir = "./dist/" + os + "-amd64"
	}
	switch os {
	case "windows":
		ostyp = "win"
	case "darwin":
		ostyp = "mac"
	case "linux":
		ostyp = "lin"
	default:
		ostyp = "lin"
	}
	cmd := exec.Command(distdir+"/kpasscli", append([]string{"-c", "config-" + ostyp + ".yaml"}, args...)...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func TestKpasscliPassword(t *testing.T) {
	result, err := runKpasscli("-i", "pw2")
	if err != nil {
		t.Fatalf("Error running kpasscli: %v", err)
	}
	expected := "password2"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestKpasscliUsername(t *testing.T) {
	result, err := runKpasscli("-i", "pw1.1", "-f", "username")
	if err != nil {
		t.Fatalf("Error running kpasscli: %v", err)
	}
	expected := "user1.1"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestKpasscliMultipleEntries(t *testing.T) {
	result, err := runKpasscli("-i", "pw")
	if err != nil {
		t.Logf("Expected error running kpasscli: %v", err)
		return
	}
	fmt.Println("Result:", result)
	expected := `- /Root/testgroup1/testgroup1.1/testpw1.1
- /Root/testgroup2/testpw2`
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
