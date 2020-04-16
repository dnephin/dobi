package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dnephin/dobi/logging"
	"github.com/docker/docker/pkg/fileutils"
	git "github.com/gogits/git-module"
)

// LastModifiedSearch provides the means by which to specify your search parameters when
// finding the last modified file.
type LastModifiedSearch struct {
	// Root must be set to the absolute path of the directory to traverse. Any
	// relative paths in Paths and Excludes will be considered relative to this
	// root directory.
	Root     string
	UseGit   bool
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
	var lastMod time.Time
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

		if lastMod = fileLastModified(search.UseGit, filePath, info); lastMod.After(latest) { // nolint: lll
			latest = lastMod
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

			if lastMod = fileLastModified(search.UseGit, path, info); lastMod.After(latest) {
				latest = lastMod
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

func fileLastModified(useGit bool, filePath string, info os.FileInfo) time.Time {
	if useGit {
		cmd := git.NewCommand("log", "--max-count=1", `--format=%ci`, "--", filePath)
		cwd, err := os.Getwd()
		if err != nil {
			logging.Log.Warn("could not determine current working directory — falling back to file mtime") // nolint: lll
			return info.ModTime()
		}
		stdout, err := cmd.RunInDir(cwd)
		if err != nil {
			logging.Log.Warnf("there was an error grabbing commit information for %s — falling back to file mtime", filePath) // nolint: lll
			return info.ModTime()
		}
		latest, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(stdout))
		if err != nil {
			logging.Log.Warnf("could not determine committer date for %s — falling back to file mtime", filePath) // nolint: lll
			return info.ModTime()
		}
		return latest
	}

	return info.ModTime()
}
