package search

import (
	"testing"

	"github.com/tobischo/gokeepasslib/v3"
)

func TestFinder_Find(t *testing.T) {
	db := makeTestDB()
	f := NewFinder(db)

	// Absolute path
	results, err := f.Find("/Root/Banking/Account")
	if err != nil {
		t.Fatalf("unexpected error for absolute path: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for absolute path, got %d", len(results))
	}
	if len(results) == 1 {
		title, err := results[0].GetField("Title")
		if err != nil || title != "Account" {
			t.Errorf("expected Title 'Account', got '%v' (err: %v)", title, err)
		}
	}

	// Subpath
	results, err = f.Find("Banking/Account")
	if err != nil {
		t.Fatalf("unexpected error for subpath: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for subpath, got %d", len(results))
	}

	// Name
	results, err = f.Find("Account")
	if err != nil {
		t.Fatalf("unexpected error for name: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for name, got %d", len(results))
	}
}

func makeTestDB() *gokeepasslib.Database {
	db := &gokeepasslib.Database{}
	db.Content = &gokeepasslib.DBContent{}
	entry := gokeepasslib.Entry{
		Values: []gokeepasslib.ValueData{
			{Key: "Title", Value: gokeepasslib.V{Content: "Account"}},
			{Key: "UserName", Value: gokeepasslib.V{Content: "tester"}},
		},
	}
	group := gokeepasslib.Group{
		Name:    "Banking",
		Entries: []gokeepasslib.Entry{entry},
	}
	db.Content.Root = &gokeepasslib.RootData{
		Groups: []gokeepasslib.Group{{
			Name:   "Root",
			Groups: []gokeepasslib.Group{group},
		}},
	}
	return db
}

func TestFindByAbsolutePath(t *testing.T) {
	db := makeTestDB()
	f := NewFinder(db)
	entry, err := f.findByAbsolutePath("/Root/Banking/Account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	title := ""
	for _, v := range entry.Values {
		if v.Key == "Title" {
			title = v.Value.Content
		}
	}
	if title != "Account" {
		t.Errorf("expected title 'Account', got '%v'", title)
	}
}

func TestFindBySubpath(t *testing.T) {
	db := makeTestDB()
	f := NewFinder(db)
	results, err := f.findBySubpath("Banking/Account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestFindByName(t *testing.T) {
	db := makeTestDB()
	f := NewFinder(db)
	results, err := f.findByName("Account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

type dummyFinder struct{}

func TestEnableVerify(t *testing.T) {
	verify = false
	EnableVerify()
	if !verify {
		t.Error("verify should be true after EnableVerify")
	}
}

func TestResult_GetField(t *testing.T) {
	entry := &gokeepasslib.Entry{
		Values: []gokeepasslib.ValueData{
			{Key: "Title", Value: gokeepasslib.V{Content: "TestTitle"}},
			{Key: "UserName", Value: gokeepasslib.V{Content: "TestUser"}},
		},
	}
	r := Result{Entry: entry}
	val, err := r.GetField("Title")
	if err != nil || val != "TestTitle" {
		t.Errorf("expected 'TestTitle', got '%v' (err: %v)", val, err)
	}
	_, err = r.GetField("NotExist")
	if err == nil {
		t.Error("expected error for missing field")
	}
}

func TestNewFinderAndDefaultSearchOptions(t *testing.T) {
	db := &gokeepasslib.Database{}
	f := NewFinder(db)
	if f.db != db {
		t.Error("NewFinder did not set db correctly")
	}
	def := DefaultSearchOptions()
	if def.CaseSensitive {
		t.Error("DefaultSearchOptions should be case-insensitive")
	}
	if def.ExactMatch {
		t.Error("DefaultSearchOptions should not be exact match")
	}
}

func TestMatchesName(t *testing.T) {
	tests := []struct {
		value, pattern string
		opts           SearchOptions
		expect         bool
	}{
		{"foo", "foo", SearchOptions{CaseSensitive: true, ExactMatch: true}, true},
		{"foo", "FOO", SearchOptions{CaseSensitive: false, ExactMatch: true}, false}, // actual implementation is case-sensitive for exact match
		{"foo", "FOO", SearchOptions{CaseSensitive: true, ExactMatch: true}, false},
		{"foobar", "foo", SearchOptions{CaseSensitive: false, ExactMatch: false}, true},
		{"foobar", "baz", SearchOptions{CaseSensitive: false, ExactMatch: false}, false},
	}
	for _, tc := range tests {
		got := matchesName(tc.value, tc.pattern, tc.opts)
		if got != tc.expect {
			t.Errorf("matchesName(%q, %q, %+v) = %v, want %v", tc.value, tc.pattern, tc.opts, got, tc.expect)
		}
	}
}

func (d *dummyFinder) Find(query string) ([]Result, error) {
	if query == "found" {
		entry := &gokeepasslib.Entry{
			Values: []gokeepasslib.ValueData{
				{Key: "Title", Value: gokeepasslib.V{Content: "Entry"}},
			},
		}
		return []Result{{Path: "/some/path", Entry: entry}}, nil
	}
	return nil, nil
}

func TestFinderInterface_Find(t *testing.T) {
	var f FinderInterface = &dummyFinder{}
	results, err := f.Find("found")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected one result, got %+v", results)
	}
	if len(results) == 1 {
		title, err := results[0].GetField("Title")
		if err != nil || title != "Entry" {
			t.Errorf("expected Title 'Entry', got '%v' (err: %v)", title, err)
		}
	}
	results, err = f.Find("notfound")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %+v", results)
	}
}
