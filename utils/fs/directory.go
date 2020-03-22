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
func LastModified(fileOrDir ...string) (time.Time, error) {
	var latest time.Time

	// TODO: does this error contain enough context?
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".dobi" {
			return filepath.SkipDir
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
