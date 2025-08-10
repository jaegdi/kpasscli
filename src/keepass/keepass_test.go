package keepass

import (
	"io/ioutil"
	"kpasscli/src/config"
	"os"
	"testing"
)

func TestOpenDatabase_FileNotFound(t *testing.T) {
	_, err := OpenDatabase("/nonexistent/file.kdbx", "pw")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestOpenDatabase_InvalidFile(t *testing.T) {
	f, err := ioutil.TempFile("", "invalid.kdbx")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("notakeepassfile")
	f.Close()
	_, err = OpenDatabase(f.Name(), "pw")
	if err == nil {
		t.Error("expected error for invalid file")
	}
}

func TestResolveDatabasePath(t *testing.T) {
	os.Setenv("KPASSCLI_KDBPATH", "envpath.kdbx")
	cfg := &config.Config{DatabasePath: "cfgpath.kdbx"}
	if got := ResolveDatabasePath("flag.kdbx", cfg); got != "flag.kdbx" {
		t.Errorf("flag path not used: got %v", got)
	}
	if got := ResolveDatabasePath("", cfg); got != "envpath.kdbx" {
		t.Errorf("env path not used: got %v", got)
	}
	os.Unsetenv("KPASSCLI_KDBPATH")
	if got := ResolveDatabasePath("", cfg); got != "cfgpath.kdbx" {
		t.Errorf("config path not used: got %v", got)
	}
	if got := ResolveDatabasePath("", &config.Config{}); got != "" {
		t.Errorf("expected empty string, got %v", got)
	}
}

func TestResolvePassword_PromptFallback(t *testing.T) {
	want := "promptedpass"
	mockPrompt := func() (string, error) {
		return want, nil
	}
	got, err := ResolvePassword("", &config.Config{}, "", mockPrompt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("expected '%v', got '%v'", want, got)
	}
}

func TestResolvePassword_FromFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "pwfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("secretpass\n")
	tmpfile.Close()

	pass, err := ResolvePassword(tmpfile.Name(), &config.Config{}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass != "secretpass" {
		t.Errorf("expected 'secretpass', got '%v'", pass)
	}
}

func TestResolvePassword_FromEnv(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "pwenv.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("envpass\n")
	tmpfile.Close()

	pass, err := ResolvePassword("", &config.Config{}, tmpfile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass != "envpass" {
		t.Errorf("expected 'envpass', got '%v'", pass)
	}
}

func TestResolvePassword_FromConfig(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "pwcfg.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("cfgpass\n")
	tmpfile.Close()

	cfg := &config.Config{PasswordFile: tmpfile.Name()}
	pass, err := ResolvePassword("", cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass != "cfgpass" {
		t.Errorf("expected 'cfgpass', got '%v'", pass)
	}
}

func TestResolvePassword_FileNotExist(t *testing.T) {
	_, err := ResolvePassword("/nonexistent/file.txt", &config.Config{}, "")
	if err == nil {
		t.Error("expected error for nonexistent password file")
	}
}

func TestResolvePassword_FromExecutable(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "pwexec.sh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	script := "#!/bin/sh\necho execpass\n"
	if _, err := tmpfile.WriteString(script); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	if err := os.Chmod(tmpfile.Name(), 0755); err != nil {
		t.Fatal(err)
	}

	pass, err := ResolvePassword(tmpfile.Name(), &config.Config{}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass != "execpass" {
		t.Errorf("expected 'execpass', got '%v'", pass)
	}
}

func Test_getPasswordFromPrompt_success(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	input := []byte("secret\n")
	go func() {
		w.Write(input)
		w.Close()
	}()

	pw, err := getPasswordFromPromptWithReader(r, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw != "secret" {
		t.Errorf("expected 'secret', got '%v'", pw)
	}
}

func Test_getPasswordFromPrompt_empty(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	input := []byte("\n")
	go func() {
		w.Write(input)
		w.Close()
	}()

	pw, _ := getPasswordFromPromptWithReader(r, -1)
	if pw != "" {
		t.Errorf("expected empty password, got '%v'", pw)
	}
}
