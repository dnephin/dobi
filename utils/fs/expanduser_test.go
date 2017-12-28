package fs

import (
	"os"
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
)

func TestExpandUserNothingToExpand(t *testing.T) {
	expected := "does/not/expand"
	path, err := ExpandUser(expected)

	assert.Check(t, is.Nil(err))
	assert.Check(t, is.Equal(expected, path))
}

func TestExpandUserJustTilde(t *testing.T) {
	path, err := ExpandUser("~")

	assert.Check(t, is.Nil(err))
	assert.Check(t, is.Equal(os.Getenv("HOME"), path))
}

func TestExpandUserCurrentUser(t *testing.T) {
	path, err := ExpandUser("~/rest/of/path")
	expected := os.Getenv("HOME") + "/rest/of/path"

	assert.Check(t, is.Nil(err))
	assert.Check(t, is.Equal(expected, path))
}

func TestExpandUserOtherUser(t *testing.T) {
	_, err := ExpandUser("~otheruser/rest/of/path")

	assert.Check(t, is.ErrorContains(err, ""))
}
