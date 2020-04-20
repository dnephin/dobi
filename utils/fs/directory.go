package fs

import (
	"os"
	"path/filepath"
	"time"
)

// LastModified returns the latest modified time for all the files and
// directories. The files in each directory are checked for their last modified
// time.
// TODO: use go routines to speed this up
// nolint: gocyclo
func LastModified(excludes []string, fileOrDir ...string) (time.Time, error) {
	var latest time.Time
	var excludedFileOrDir string

	// TODO: does this error contain enough context?
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, excludedFileOrDir = range excludes {
			if excludedFileOrDir == info.Name() {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
		return nil
	}

	for _, file := range fileOrDir {
		info, err := os.Stat(file)
		if err != nil {
			return latest, err
		}
		switch info.IsDir() {
		case false:
			for _, excludedFileOrDir = range excludes {
				if excludedFileOrDir == info.Name() {
					continue
				}
			}
			if info.ModTime().After(latest) {
				latest = info.ModTime()
				continue
			}
		default:
			if err := filepath.Walk(file, walker); err != nil {
				return latest, err
			}
		}
	}
	return latest, nil
}
