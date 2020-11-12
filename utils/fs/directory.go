package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/pkg/fileutils"
)

// LastModifiedSearch provides the means by which to specify your search parameters when
// finding the last modified file.
type LastModifiedSearch struct {
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
	var rootPath string
	var err error

	rootPath = search.Root
	if rootPath == "" {
		if rootPath, err = os.Getwd(); err != nil {
			return time.Time{}, err
		}
	}

	// Make absolute path out of `rootPath` to ensure proper functionality later on
	if !filepath.IsAbs(rootPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return time.Time{}, err
		}
		rootPath = cwd + string(os.PathSeparator) + rootPath
	}

	pm, err := fileutils.NewPatternMatcher(search.Excludes)
	if err != nil {
		return time.Time{}, err
	}

	walker := func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("can't stat '%s'", filePath)
			}
			return err
		}
		if relFilePath, err := filepath.Rel(rootPath, filePath); err != nil {
			return err
		} else if skip, err := filepathMatches(pm, relFilePath); err != nil {
			return err
		} else if skip {
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
		// Append the cwd if `path` is relative. This is needed because
		// source/artifact handling (`filepath.Rel()`) needs absolute paths to work properly.
		if !filepath.IsAbs(path) {
			cwd, err := os.Getwd()
			if err != nil {
				return time.Time{}, err
			}
			path = cwd + string(os.PathSeparator) + path
		}
		info, err := os.Stat(path)
		if err != nil {
			return latest, err
		}
		switch info.IsDir() {
		case false:
			if relPath, err := filepath.Rel(rootPath, path); err != nil {
				return time.Time{}, err
			} else if skip, err := filepathMatches(pm, relPath); err != nil {
				return time.Time{}, err
			} else if skip {
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

func filepathMatches(matcher *fileutils.PatternMatcher, file string) (bool, error) {
	file = filepath.Clean(file)
	if file == "." {
		// Don't let them exclude everything, kind of silly.
		return false, nil
	}
	return matcher.Matches(file)
}
