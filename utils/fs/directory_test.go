package fs

import (
	"os"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/fs"
)

func TestLastModified_AbsolutePathsForDirectories(t *testing.T) {
	tmpdir := fs.NewDir(t, "test-directory-last-modified-absolute-paths-for-dir",
		fs.WithDir("a"),
		fs.WithDir("b",
			fs.WithDir("c")))
	defer tmpdir.Remove()

	for index, dir := range []string{"a", "b", "b/c"} {
		mtime := time.Now().AddDate(0, 0, index+10)
		assert.NilError(t, touch(tmpdir.Join(dir, "file"), mtime))

		actual, err := LastModified(&LastModifiedSearch{
			Paths: []string{tmpdir.Path()},
			Root:  tmpdir.Path(),
		})
		assert.NilError(t, err)
		assert.Equal(t, actual, mtime)
	}
}

func TestLastModified_RelativePathsForDirectories(t *testing.T) {
	tmpdir := fs.NewDir(t, "test-directory-last-modified-relative-paths-for-dir",
		fs.WithDir("inner",
			fs.WithDir("a"),
			fs.WithDir("b",
				fs.WithDir("c"))))
	defer tmpdir.Remove()

	for index, dir := range []string{"a", "b", "b/c"} {
		mtime := time.Now().AddDate(0, 0, index+10)
		assert.NilError(t, touch(tmpdir.Join("inner", dir, "file"), mtime))

		actual, err := LastModified(&LastModifiedSearch{
			Root:  tmpdir.Join("inner"),
			Paths: []string{"a", "b", "b/c"},
		})
		assert.NilError(t, err)
		assert.Equal(t, actual, mtime)
	}
}

func TestLastModified_RelativePathsForFile(t *testing.T) {
	tmpdir := fs.NewDir(t, "test-directory-last-modified-relative-paths-for-file",
		fs.WithDir("inner",
			fs.WithDir("a")))
	defer tmpdir.Remove()

	mtime := time.Now().AddDate(0, 0, 10)
	assert.NilError(t, touch(tmpdir.Join("inner", "a", "file"), mtime))

	actual, err := LastModified(&LastModifiedSearch{
		Root:  tmpdir.Join("inner"),
		Paths: []string{"a/file"},
	})
	assert.NilError(t, err)
	assert.Equal(t, actual, mtime)
}

func TestLastModifiedExcludesFile(t *testing.T) {
	tmpdir := fs.NewDir(t, "test-directory-last-modified-excludes-file",
		fs.WithDir("a"),
		fs.WithDir("b",
			fs.WithDir("c")))
	defer tmpdir.Remove()

	for index, dir := range []string{"a", "b", "b/c"} {
		mtime := time.Now().AddDate(0, 0, index+10)
		assert.Assert(t, cmp.Nil(touch(tmpdir.Join(dir, "file"), mtime)))

		excludedFile := tmpdir.Join(dir, "excluded-file")
		excludedMtime := time.Now().AddDate(0, 0, index+20)
		assert.Assert(t, cmp.Nil(touch(excludedFile, excludedMtime)))

		actual, err := LastModified(&LastModifiedSearch{
			Excludes: []string{"**/**/excluded-file"},
			Paths:    []string{tmpdir.Path()},
			Root:     tmpdir.Path(),
		})
		assert.NilError(t, err)
		assert.Equal(t, actual, mtime)
	}
}

func TestLastModifiedExcludesFolder(t *testing.T) {
	tmpdir := fs.NewDir(t, "test-directory-last-modified-excludes-folder",
		fs.WithDir("a"),
		fs.WithDir("b",
			fs.WithDir("c")))
	defer tmpdir.Remove()

	mtime := time.Now().AddDate(0, 0, 0)
	assert.Assert(t, cmp.Nil(touch(tmpdir.Join("a", "file"), mtime)))

	for index, dir := range []string{"b", "b/c"} {
		ignoredMtime := time.Now().AddDate(0, 0, index+10)
		assert.Assert(t, cmp.Nil(touch(tmpdir.Join(dir, "file"), ignoredMtime)))
	}

	actual, err := LastModified(&LastModifiedSearch{
		Excludes: []string{"b/"},
		Paths:    []string{tmpdir.Path()},
		Root:     tmpdir.Path(),
	})
	assert.NilError(t, err)
	assert.Equal(t, actual, mtime)
}

func touch(name string, mtime time.Time) error {
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return os.Chtimes(name, mtime, mtime)
}
