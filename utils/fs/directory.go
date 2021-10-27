package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/moby/moby/pkg/fileutils"
)

// LastModifiedSearch provides the means by which to specify your search parameters when
// finding the last modified file.
type LastModifiedSearch struct {
	// Root must be set to the absolute path of the directory to traverse. Any
	// relative paths in Paths and Excludes will be considered relative to this
	// root directory.
	Root     string
	Excludes []string
	Paths    []string
}

// LastModified returns the latest modified time for all the files and
// directories. The files in each directory are checked for their last modified
// time.
// TODO: use go routines to speed this up
// nolint: gocyclo
func LastModified(search *LastModifiedSearch) (time.Time, error) {
	var latest time.Time
	var err error

	pm, err := fileutils.NewPatternMatcher(search.Excludes)
	if err != nil {
		return time.Time{}, err
	}

	isExcluded := func(path string) (bool, error) {
		relPath, err := filepath.Rel(search.Root, path)
		if err != nil {
			return false, err
		}
		if relPath == "." {
			// Don't let them exclude everything, kind of silly.
			return false, nil
		}
		return pm.Matches(relPath)
	}

	walker := func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("can't stat '%s'", filePath)
			}
			return err
		}

		skip, err := isExcluded(filePath)
		switch {
		case err != nil:
			return err
		case skip:
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
		return nil
	}

	for _, path := range search.Paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(search.Root, path)
		}

		info, err := os.Stat(path)
		if err != nil {
			return latest, fmt.Errorf("internal error: %w", err)
		}
		switch info.IsDir() {
		case false:
			skip, err := isExcluded(path)
			switch {
			case err != nil:
				return time.Time{}, err
			case skip:
				continue
			}

			if info.ModTime().After(latest) {
				latest = info.ModTime()
				continue
			}
		default:
			if err := filepath.Walk(path, walker); err != nil {
				return latest, err
			}
		}
	}
	return latest, nil
}
