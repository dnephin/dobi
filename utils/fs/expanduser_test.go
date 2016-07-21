package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandUserNothingToExpand(t *testing.T) {
	expected := "does/not/expand"
	path, err := ExpandUser(expected)

	assert.Nil(t, err)
	assert.Equal(t, expected, path)
}

func TestExpandUserJustTilde(t *testing.T) {
	path, err := ExpandUser("~")

	assert.Nil(t, err)
	assert.Equal(t, os.Getenv("HOME"), path)
}

func TestExpandUserCurrentUser(t *testing.T) {
	path, err := ExpandUser("~/rest/of/path")
	expected := os.Getenv("HOME") + "/rest/of/path"

	assert.Nil(t, err)
	assert.Equal(t, expected, path)
}

func TestExpandUserOtherUser(t *testing.T) {
	_, err := ExpandUser("~otheruser/rest/of/path")

	assert.Error(t, err)
}
