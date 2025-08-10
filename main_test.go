package main

import (
	"errors"
	"testing"

	"github.com/tobischo/gokeepasslib/v3"

	"kpasscli/src/cmd"
	"kpasscli/src/config"
	"kpasscli/src/keepass"
	"kpasscli/src/output"
	"kpasscli/src/search"
)

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
}

// Ensure fakeHandler implements output.Handler
var _ output.Handler = (*fakeHandler)(nil)

type fakeResult struct {
	fieldVal string
	fieldErr error
	Path     string
}

func (r fakeResult) GetField(field string) (string, error) {
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

func fakeNewHandler(outputErr error) func(output.Type) output.Handler {
	return func(output.Type) output.Handler {
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
		func(string) string { return "" },
	)
	if err == nil || err.Error() != "multiple items found" {
		t.Errorf("expected multiple items found error, got %v", err)
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
		func(string) string { return "password" },
	)
	if err != nil {
		t.Errorf("expected success, got %v", err)
	}
}
