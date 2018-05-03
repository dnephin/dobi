package image

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestParseAuthRepoWithUserRepo(t *testing.T) {
	repo, err := parseAuthRepo("dnephin/foo")
	assert.NilError(t, err)
	assert.Check(t, is.Equal(repo, defaultRepo))
}

func TestParseAuthRepoPrivateRepoAndPort(t *testing.T) {
	repo, err := parseAuthRepo("myrepo.net:3434/dnephin/foo")
	assert.NilError(t, err)
	assert.Check(t, is.Equal(repo, "myrepo.net:3434"))
}

func TestParseAuthRepoPrivateRepo(t *testing.T) {
	repo, err := parseAuthRepo("myrepo.net/dnephin/foo:tag")
	assert.NilError(t, err)
	assert.Check(t, is.Equal(repo, "myrepo.net"))
}

func TestParseAuthRepoPrivateRepoNoUsername(t *testing.T) {
	repo, err := parseAuthRepo("myrepo.net/foo")
	assert.NilError(t, err)
	assert.Check(t, is.Equal(repo, "myrepo.net"))
}
