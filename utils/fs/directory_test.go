package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DirectorySuite struct {
	suite.Suite
	path string
}

func TestDirectorySuite(t *testing.T) {
	suite.Run(t, new(DirectorySuite))
}

func (s *DirectorySuite) SetupTest() {
	var err error
	s.path, err = ioutil.TempDir("", "directory-test")
	s.Require().Nil(err)
}

func (s *DirectorySuite) TearDownTest() {
	s.Nil(os.RemoveAll(s.path))
}

func (s *DirectorySuite) TestLastModified() {
	dirs := []string{
		filepath.Join(s.path, "a"),
		filepath.Join(s.path, "b"),
		filepath.Join(s.path, "b", "c"),
	}
	for _, dir := range dirs {
		s.Require().Nil(os.MkdirAll(dir, 0777))
	}

	assertModTime := func(days int, dir string) {
		mtime := time.Now().AddDate(0, 0, days)
		file := filepath.Join(dir, "file")
		s.Require().Nil(touch(file, mtime))
		actual, err := LastModified(dirs...)
		s.Equal(actual, mtime)
		s.Nil(err)
	}

	for index, dir := range dirs {
		assertModTime(10+index, dir)
	}
}

func touch(name string, mtime time.Time) error {
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	w.Close()

	return os.Chtimes(name, mtime, mtime)
}
