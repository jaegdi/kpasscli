package search

import (
	"fmt"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
)

type Result struct {
	Path  string
	Entry *gokeepasslib.Entry
}

type Finder struct {
	db *gokeepasslib.Database
}

func NewFinder(db *gokeepasslib.Database) *Finder {
	return &Finder{db: db}
}

func (f *Finder) Find(query string) ([]Result, error) {
	var results []Result

	if strings.HasPrefix(query, "/") {
		// Absolute path search
		entry, err := f.findByAbsolutePath(query)
		if err != nil {
			return nil, fmt.Errorf("absolute path search failed: %w", err)
		}
		if entry != nil {
			results = append(results, Result{Path: query, Entry: entry})
		}
	} else if strings.Contains(query, "/") {
		// Subpath search
		var err error
		results, err = f.findBySubpath(query)
		if err != nil {
			return nil, fmt.Errorf("subpath search failed: %w", err)
		}
	} else {
		// Name search
		var err error
		results, err = f.findByName(query)
		if err != nil {
			return nil, fmt.Errorf("name search failed: %w", err)
		}
	}

	return results, nil
}

func (f *Finder) findByAbsolutePath(path string) (*gokeepasslib.Entry, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	if parts[0] == "" {
		return nil, fmt.Errorf("invalid path format")
	}

	currentGroup := &f.db.Content.Root.Groups[0]

	// Navigate through groups
	for i := 0; i < len(parts)-1; i++ {
		found := false
		for _, group := range currentGroup.Groups {
			if group.Name == parts[i] {
				currentGroup = &group
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("group not found: %s", parts[i])
		}
	}

	// Search for entry in final group
	targetName := parts[len(parts)-1]
	for _, entry := range currentGroup.Entries {
		if entry.GetTitle() == targetName {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("entry not found: %s", targetName)
}

func (f *Finder) findBySubpath(query string) ([]Result, error) {
	// TODO: Implement subpath search
	return []Result{}, nil
}

func (f *Finder) findByName(query string) ([]Result, error) {
	// TODO: Implement name search
	return []Result{}, nil
}

// Helper method to get a field value from an entry
func (r *Result) GetField(fieldName string) (string, error) {
	switch strings.ToLower(fieldName) {
	case "title":
		return r.Entry.GetTitle(), nil
	case "username":
		return r.Entry.GetContent("Username"), nil
	case "password":
		return r.Entry.GetPassword(), nil
	case "url":
		return r.Entry.GetContent("URL"), nil
	case "notes":
		return r.Entry.GetContent("Notes"), nil
	default:
		// Search in custom fields
		for _, v := range r.Entry.Values {
			if strings.EqualFold(v.Key, fieldName) {
				return v.Value.Content, nil
			}
		}
		return "", fmt.Errorf("field '%s' not found", fieldName)
	}
}
