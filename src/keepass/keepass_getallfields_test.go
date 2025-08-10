package keepass

import (
	"errors"
	"testing"

	"kpasscli/src/config"
	"kpasscli/src/search"

	"github.com/tobischo/gokeepasslib/v3"
)

type mockFinder struct {
	results []search.Result
	err     error
}

func (m *mockFinder) Find(query string) ([]search.Result, error) {
	return m.results, m.err
}

func TestGetAllFields_NoResults(t *testing.T) {
	cfg := &config.Config{}
	db := &gokeepasslib.Database{}
	itemPath := "notfound"

	err := GetAllFieldsWithFinder(db, cfg, itemPath, &mockFinder{results: nil, err: nil}, nil)
	if err == nil || err.Error() != "entry not found: notfound" {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestGetAllFields_MultipleResults(t *testing.T) {
	cfg := &config.Config{}
	db := &gokeepasslib.Database{}
	itemPath := "ambiguous"

	results := []search.Result{
		{Path: "/a/b", Entry: &gokeepasslib.Entry{}},
		{Path: "/a/c", Entry: &gokeepasslib.Entry{}},
	}
	err := GetAllFieldsWithFinder(db, cfg, itemPath, &mockFinder{results: results, err: nil}, nil)
	if err == nil || !errors.Is(err, err) {
		t.Errorf("expected multiple entries error, got: %v", err)
	}
}

func TestGetAllFields_NilEntry(t *testing.T) {
	cfg := &config.Config{}
	db := &gokeepasslib.Database{}
	itemPath := "nilentry"

	results := []search.Result{{Path: "/a/b", Entry: nil}}
	err := GetAllFieldsWithFinder(db, cfg, itemPath, &mockFinder{results: results, err: nil}, nil)
	if err == nil || err.Error() != "found result for 'nilentry', but entry data is unexpectedly nil" {
		t.Errorf("expected nil entry error, got: %v", err)
	}
}

func TestGetAllFields_Success(t *testing.T) {
	cfg := &config.Config{}
	db := &gokeepasslib.Database{}
	itemPath := "ok"

	entry := &gokeepasslib.Entry{}
	results := []search.Result{{Path: "/a/b", Entry: entry}}
	called := false
	showAllFields := func(e *gokeepasslib.Entry, c config.Config) {
		called = true
		if e != entry {
			t.Errorf("ShowAllFields called with wrong entry")
		}
	}
	err := GetAllFieldsWithFinder(db, cfg, itemPath, &mockFinder{results: results, err: nil}, showAllFields)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if !called {
		t.Error("ShowAllFields was not called")
	}
}
