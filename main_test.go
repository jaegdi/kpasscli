package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tobischo/gokeepasslib/v3"
	"golang.design/x/clipboard"

	"kpasscli/src/cmd"
	"kpasscli/src/config"
	"kpasscli/src/keepass"
	"kpasscli/src/output"
	"kpasscli/src/search"
)

// MockClipboard implements ClipboardService for testing
type MockClipboard struct {
	InitFunc  func() error
	WriteFunc func(t clipboard.Format, content []byte) <-chan struct{}
	Written   []string
}

func (m *MockClipboard) Init() error {
	if m.InitFunc != nil {
		return m.InitFunc()
	}
	return nil
}

func (m *MockClipboard) Write(t clipboard.Format, content []byte) <-chan struct{} {
	m.Written = append(m.Written, string(content))
	if m.WriteFunc != nil {
		return m.WriteFunc(t, content)
	}
	return nil
}

// noResultFinder implements search.FinderInterface and always returns no results
type noResultFinder struct{}

func (f *noResultFinder) Find(item string) ([]search.Result, error) {
	return nil, nil
}

// searchErrorFinder implements search.FinderInterface and always returns a search error
type searchErrorFinder struct{}

func (f *searchErrorFinder) Find(item string) ([]search.Result, error) {
	return nil, errors.New("searchfail")
}

// --- Fakes for dependency injection ---

type fakeHandler struct {
	outputErr error
	captured  string
}

// Ensure fakeHandler implements output.Handler
var _ output.Handler = (*fakeHandler)(nil)

type fakeResult struct {
	fieldVal string
	fieldErr error
	Path     string
	fields   map[string]string
}

func (r fakeResult) GetField(field string) (string, error) {
	if r.fields != nil {
		if val, ok := r.fields[field]; ok {
			return val, nil
		}
		// If fields map is provided but key not found, return error or empty?
		// For this test, let's assume if map is present, we look there.
		// But to keep backward compatibility with other tests that don't set fields,
		// we should fall back to fieldVal/fieldErr if fields is nil.
		return "", fmt.Errorf("field '%s' not found", field)
	}
	if r.fieldErr != nil {
		return "", r.fieldErr
	}
	return r.fieldVal, nil
}

type fakeFinder struct {
	results []fakeResult
	findErr error
	Options interface{}
}

func (f *fakeFinder) Find(item string) ([]fakeResult, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.results, nil
}

// Ensure fakeHandler implements output.Handler
var _ output.Handler = (*fakeHandler)(nil)

func (h *fakeHandler) Output(val string) error {
	h.captured = val
	return h.outputErr
}

// --- Test helpers for dependency injection ---

func fakeLoadConfig(err error) func(string) (*config.Config, error) {
	return func(string) (*config.Config, error) {
		if err != nil {
			return nil, err
		}
		return &config.Config{DefaultOutput: "stdout"}, nil
	}
}

func fakeResolveDBPath(path string) func(string, *config.Config) string {
	return func(string, *config.Config) string { return path }
}

func fakeResolvePassword(pw string, err error) func(string, *config.Config, string, ...keepass.PasswordPromptFunc) (string, error) {
	return func(string, *config.Config, string, ...keepass.PasswordPromptFunc) (string, error) {
		return pw, err
	}
}

func fakeOpenDatabase(db *gokeepasslib.Database, err error) func(string, string) (*gokeepasslib.Database, error) {
	return func(string, string) (*gokeepasslib.Database, error) {
		return db, err
	}
}

// FakeFinder implements the Find method for testing
type FakeFinder struct {
	results []search.Result
	findErr error
	Options search.SearchOptions
}

func (f *FakeFinder) Find(item string) ([]search.Result, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.results, nil
}

func fakeNewFinder(results []search.Result, findErr error) func(*gokeepasslib.Database) *FakeFinder {
	return func(*gokeepasslib.Database) *FakeFinder {
		return &FakeFinder{results: results, findErr: findErr}
	}
}

func fakeNewHandler(outputErr error) func(output.OutputType, output.ClipboardService) output.Handler {
	return func(output.OutputType, output.ClipboardService) output.Handler {
		return &fakeHandler{outputErr: outputErr}
	}
}

// --- Tests for RunApp ---

func TestRunApp_ItemRequired(t *testing.T) {
	flags := &cmd.Flags{}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return search.NewFinder(db) },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "item parameter is required" {
		t.Errorf("expected item parameter error, got %v", err)
	}
}

func TestRunApp_ConfigLoadError(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	err := RunApp(
		flags,
		fakeLoadConfig(errors.New("fail")),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return search.NewFinder(db) },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err != nil {
		t.Errorf("expected warning only, got %v", err)
	}
}

func TestRunApp_NoDBPath(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath(""),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return search.NewFinder(db) },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "no KeePass database path provided" {
		t.Errorf("expected db path error, got %v", err)
	}
}

func TestRunApp_PasswordError(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("", errors.New("pwfail")),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return search.NewFinder(db) },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "Error getting password: pwfail" {
		t.Errorf("expected password error, got %v", err)
	}
}

func TestRunApp_DatabaseOpenError(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, errors.New("dbfail")),
		func(db *gokeepasslib.Database) search.FinderInterface { return search.NewFinder(db) },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "Error opening database: dbfail" {
		t.Errorf("expected db open error, got %v", err)
	}
}

func TestRunApp_SearchError(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return &searchErrorFinder{} },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "Error searching for item: searchfail" {
		t.Errorf("expected search error, got %v", err)
	}
}

func TestRunApp_NoResults(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return &noResultFinder{} },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "no items found" {
		t.Errorf("expected no items found error, got %v", err)
	}
}

func TestRunApp_MultipleResults(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	fakeResults := []search.Result{
		{Path: "entry1", Entry: nil},
		{Path: "entry2", Entry: nil},
	}
	fakeFinder := &FakeFinder{results: fakeResults}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "multiple items found" {
		t.Errorf("expected multiple items found error, got %v", err)
	}
}

func TestRunApp_ClearClipboard(t *testing.T) {
	flags := &cmd.Flags{ClearClipboard: true, ClearAfter: 0} // 0 to avoid sleep in test
	mockClipboard := &MockClipboard{}

	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return &noResultFinder{} },
		fakeNewHandler(nil),
		mockClipboard,
		func(string) string { return "" },
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}

	// Verify that Write was called with empty string
	if len(mockClipboard.Written) == 0 {
		t.Error("expected clipboard write, got none")
	} else {
		lastWrite := mockClipboard.Written[len(mockClipboard.Written)-1]
		if lastWrite != "" {
			t.Errorf("expected clipboard cleared (empty string), got %q", lastWrite)
		}
	}
}

func TestRunApp_GetFieldError(t *testing.T) {
	flags := &cmd.Flags{Item: "foo"}
	// Simulate a Finder that returns a result with a GetField error
	// Use a real Entry with no values to trigger GetField error
	fakeResults := []search.Result{{Path: "entry1", Entry: &gokeepasslib.Entry{}}}
	fakeFinder := &FakeFinder{results: fakeResults}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "nonexistent" },
	)
	if err == nil || err.Error() == "" || err.Error()[:20] != "Error getting field:" {
		t.Errorf("expected get field error, got %v", err)
	}
}

func TestRunApp_OutputError(t *testing.T) {
	flags := &cmd.Flags{Item: "foo", FieldName: "password"}
	fakeResults := []search.Result{
		{Path: "entry1", Entry: &gokeepasslib.Entry{
			Values: []gokeepasslib.ValueData{{Key: "password", Value: gokeepasslib.V{Content: "secret"}}},
		}},
	}
	fakeFinder := &FakeFinder{results: fakeResults}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		fakeNewHandler(errors.New("outfail")),
		&MockClipboard{},
		func(string) string { return "password" },
	)
	if err == nil || err.Error() != "Error outputting value: outfail" {
		t.Errorf("expected output error, got %v", err)
	}
}

func TestRunApp_Success(t *testing.T) {
	flags := &cmd.Flags{Item: "foo", FieldName: "password"}
	fakeResults := []search.Result{
		{Path: "entry1", Entry: &gokeepasslib.Entry{
			Values: []gokeepasslib.ValueData{{Key: "password", Value: gokeepasslib.V{Content: "secret"}}},
		}},
	}
	fakeFinder := &FakeFinder{results: fakeResults}
	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		fakeNewHandler(nil),
		&MockClipboard{},
		func(string) string { return "password" },
	)
	if err != nil {
		t.Errorf("expected success, got %v", err)
	}
}

func TestRunApp_WithTOTP(t *testing.T) {
	flags := &cmd.Flags{Item: "foo", FieldName: "password", PasswordTotp: true}
	// Mock a valid base32 secret. "JBSWY3DPEHPK3PXP" is base32 for "Hello!"
	fakeResults := []search.Result{
		{
			Path: "entry1",
			Entry: &gokeepasslib.Entry{
				Values: []gokeepasslib.ValueData{
					{Key: "password", Value: gokeepasslib.V{Content: "secret"}},
					{Key: "TimeOtp-Secret-Base32", Value: gokeepasslib.V{Content: "JBSWY3DPEHPK3PXP"}},
				},
			},
		},
	}
	// We need to use the fakeResult struct which implements the GetField method we modified
	// But wait, the Finder returns search.Result, which uses the real GetField method on *gokeepasslib.Entry.
	// In my main.go modification, I call results[0].GetField.
	// In the test, I am using a FakeFinder that returns []search.Result (or []fakeResult?).
	// The interface FinderInterface returns []Result.
	// search.Result struct has GetField method.
	// I cannot easily mock the method of a struct in Go unless it's an interface.
	// search.Result is a struct.
	// However, search.Result.GetField logic is:
	// func (r *Result) GetField(fieldName string) (string, error) { ... }
	// It iterates over r.Entry.Values.
	// So if I populate Entry.Values correctly, the real GetField method will work!
	// So I don't need to modify fakeResult struct for this test if I use the real search.Result and Entry.
	// But wait, the tests above use FakeFinder which returns []search.Result.
	// And search.Result has a real GetField method.
	// So I just need to populate the Entry in the search.Result.

	fakeFinder := &FakeFinder{results: fakeResults}

	// We need to capture the output to verify the token is appended.
	// The handler is mocked.
	mockHandler := &fakeHandler{}

	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		func(output.OutputType, output.ClipboardService) output.Handler { return mockHandler },
		&MockClipboard{},
		func(string) string { return "" },
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}

	// Verify that the output contains the password "secret"
	// We can't verify the exact TOTP token because it changes with time,
	// but we can verify it starts with "secret" and is longer than "secret".
	if len(mockHandler.captured) <= 6 {
		t.Errorf("expected output to be longer than password length (6), got %s", mockHandler.captured)
	}
	if mockHandler.captured[:6] != "secret" {
		t.Errorf("expected output to start with 'secret', got %s", mockHandler.captured)
	}
}

func TestRunApp_WithTOTP_MissingSecret(t *testing.T) {
	flags := &cmd.Flags{Item: "foo", FieldName: "password", PasswordTotp: true}
	// Mock an entry without the TOTP secret
	fakeResults := []search.Result{
		{
			Path: "entry1",
			Entry: &gokeepasslib.Entry{
				Values: []gokeepasslib.ValueData{
					{Key: "password", Value: gokeepasslib.V{Content: "secret"}},
				},
			},
		},
	}

	fakeFinder := &FakeFinder{results: fakeResults}
	mockHandler := &fakeHandler{}

	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		func(output.OutputType, output.ClipboardService) output.Handler { return mockHandler },
		&MockClipboard{},
		func(string) string { return "" },
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}

	// Verify that the output is just the password
	if mockHandler.captured != "secret" {
		t.Errorf("expected output to be 'secret', got %s", mockHandler.captured)
	}
}

func TestRunApp_WithTOTP_InvalidSecret(t *testing.T) {
	flags := &cmd.Flags{Item: "foo", FieldName: "password", PasswordTotp: true}
	// Mock an entry with an invalid base32 secret
	fakeResults := []search.Result{
		{
			Path: "entry1",
			Entry: &gokeepasslib.Entry{
				Values: []gokeepasslib.ValueData{
					{Key: "password", Value: gokeepasslib.V{Content: "secret"}},
					{Key: "TimeOtp-Secret-Base32", Value: gokeepasslib.V{Content: "INVALID_BASE32!"}},
				},
			},
		},
	}

	fakeFinder := &FakeFinder{results: fakeResults}
	mockHandler := &fakeHandler{}

	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		func(output.OutputType, output.ClipboardService) output.Handler { return mockHandler },
		&MockClipboard{},
		func(string) string { return "" },
	)

	if err == nil || err.Error()[:28] != "Error generating TOTP token:" {
		t.Errorf("expected TOTP generation error, got %v", err)
	}
}

func TestRunApp_TotpFlag(t *testing.T) {
	flags := &cmd.Flags{Item: "foo", TotpFlag: true}
	// Mock a valid base32 secret
	fakeResults := []search.Result{
		{
			Path: "entry1",
			Entry: &gokeepasslib.Entry{
				Values: []gokeepasslib.ValueData{
					{Key: "password", Value: gokeepasslib.V{Content: "secret"}},
					{Key: "TimeOtp-Secret-Base32", Value: gokeepasslib.V{Content: "JBSWY3DPEHPK3PXP"}},
				},
			},
		},
	}

	fakeFinder := &FakeFinder{results: fakeResults}
	mockHandler := &fakeHandler{}

	err := RunApp(
		flags,
		fakeLoadConfig(nil),
		fakeResolveDBPath("db"),
		fakeResolvePassword("pw", nil),
		fakeOpenDatabase(nil, nil),
		func(db *gokeepasslib.Database) search.FinderInterface { return fakeFinder },
		func(output.OutputType, output.ClipboardService) output.Handler { return mockHandler },
		&MockClipboard{},
		func(string) string { return "" },
	)

	if err != nil {
		t.Errorf("expected success, got %v", err)
	}

	// Verify that the output is a 6-digit token
	if len(mockHandler.captured) != 6 {
		t.Errorf("expected output to be 6 digits, got %s", mockHandler.captured)
	}
	// Verify it's numeric (simple check)
	for _, c := range mockHandler.captured {
		if c < '0' || c > '9' {
			t.Errorf("expected output to be numeric, got %s", mockHandler.captured)
		}
	}
}
