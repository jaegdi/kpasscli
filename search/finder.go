package search

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
)

type Result struct {
	Path  string
	Entry *gokeepasslib.Entry
}

// SearchOptions defines the search behavior
type SearchOptions struct {
	CaseSensitive bool
	ExactMatch    bool
}

// Finder handles searching through the KeePass database
type Finder struct {
	db      *gokeepasslib.Database
	Options SearchOptions // Add Options field to Finder struct
}

// NewFinder creates a new Finder instance with default options
func NewFinder(db *gokeepasslib.Database) *Finder {
	return &Finder{
		db:      db,
		Options: DefaultSearchOptions(),
	}
}

// DefaultSearchOptions returns the default search options
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		CaseSensitive: false, // Case-insensitive by default
		ExactMatch:    false, // Partial matching by default
	}
}

// Find searches for entries in the KeePass database.
// Parameters:
//
//	query: Search query, can be absolute path, relative path, or entry name
//
// Returns:
//
//	[]Result: Array of matching entries with their paths
//	error: Any error encountered during search
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

// findByAbsolutePath finds an entry using an absolute path.
// Parameters:
//
//	path: Absolute path starting with "/"
//
// Returns:
//
//	*gokeepasslib.Entry: The found entry or nil
//	error: Any error encountered during search
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

// findBySubpath searches for entries matching a relative path pattern.
// The search is performed recursively through all groups in the database.
// Parameters:
//
//	query: Relative path pattern (e.g., "Banking/Account")
//
// Returns:
//
//	[]Result: Array of matching entries with their full paths
//	error: Any error encountered during search
func (f *Finder) findBySubpath(query string) ([]Result, error) {
	parts := strings.Split(query, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid subpath query: must contain at least one '/'")
	}

	var results []Result
	targetName := parts[len(parts)-1] // Last component is the entry name
	subPath := parts[:len(parts)-1]   // Other components form the path

	// Use default search options
	opts := DefaultSearchOptions()

	// Start recursive search from root group
	err := f.searchGroupForSubpath(&f.db.Content.Root.Groups[0], "", subPath, targetName, &results, opts)
	if err != nil {
		return nil, fmt.Errorf("subpath search failed: %w", err)
	}

	return results, nil
}

// searchGroupForSubpath recursively searches through groups for matching paths.
// Parameters:
//
//	group: Current group being searched
//	currentPath: Full path to current group
//	searchPath: Remaining path components to match
//	targetName: Name of the entry to find
//	results: Slice to collect matching results
//	opts: Search options controlling matching behavior
func (f *Finder) searchGroupForSubpath(
	group *gokeepasslib.Group,
	currentPath string,
	searchPath []string,
	targetName string,
	results *[]Result,
	opts SearchOptions,
) error {
	// Build the full path for the current group
	groupPath := currentPath
	if group.Name != "" {
		if groupPath == "" {
			groupPath = group.Name
		} else {
			groupPath = filepath.Join(groupPath, group.Name)
		}
	}

	// If we're at the target depth (matched all path components)
	if len(searchPath) == 0 {
		// Search for entries with matching name in this group
		for _, entry := range group.Entries {
			var title string
			for _, v := range entry.Values {
				if v.Key == "Title" {
					title = v.Value.Content
					break
				}
			}
			if matchesName(title, targetName, opts) {
				fullPath := filepath.Join(groupPath, title)
				*results = append(*results, Result{
					Path:  "/" + fullPath, // Ensure path starts with /
					Entry: &entry,
				})
			}
		}
	}

	// If there are more path components to match
	if len(searchPath) > 0 {
		// Check if current group matches the next path component
		if matchesName(group.Name, searchPath[0], opts) {
			// Recursively search subgroups with remaining path components
			for i := range group.Groups {
				err := f.searchGroupForSubpath(&group.Groups[i], groupPath, searchPath[1:], targetName, results, opts)
				if err != nil {
					return err
				}
			}
		}
	}

	// Always search all subgroups for potential matches
	// This allows finding matches even if intermediate path components don't match exactly
	for i := range group.Groups {
		err := f.searchGroupForSubpath(&group.Groups[i], groupPath, searchPath, targetName, results, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

// matchesName checks if two strings match according to the search options
func matchesName(value, pattern string, opts SearchOptions) bool {
	if opts.CaseSensitive {
		if opts.ExactMatch {
			return value == pattern
		}
		return strings.Contains(value, pattern)
	}

	// Case-insensitive comparison
	valueLower := strings.ToLower(value)
	patternLower := strings.ToLower(pattern)

	if opts.ExactMatch {
		return valueLower == patternLower
	}
	return strings.Contains(valueLower, patternLower)
}

// GetField returns the value of the specified field from the entry
func (r *Result) GetField(fieldName string) (string, error) {
	// fieldNameLower := strings.ToLower(fieldName)

	for _, v := range r.Entry.Values {
		if strings.EqualFold(v.Key, fieldName) {
			return v.Value.Content, nil
		}
	}

	return "", fmt.Errorf("field '%s' not found", fieldName)
}

// String returns a string representation of the result
func (r *Result) String() string {
	var title string
	for _, v := range r.Entry.Values {
		if v.Key == "Title" {
			title = v.Value.Content
			break
		}
	}
	return fmt.Sprintf("%s [%s]", r.Path, title)
}

func (f *Finder) findByName(query string) ([]Result, error) {
	// TODO: Implement name search
	return []Result{}, nil
}
