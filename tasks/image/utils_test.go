package image

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestParseAuthRepoWithUserRepo(t *testing.T) {
	repo := parseAuthRepo("dnephin/foo")
	assert.Check(t, is.Equal(repo, defaultRepo))
}

func TestParseAuthRepoPrivateRepoAndPort(t *testing.T) {
	repo := parseAuthRepo("myrepo.net:3434/dnephin/foo")
	assert.Check(t, is.Equal(repo, "myrepo.net:3434"))
}

func TestParseAuthRepoPrivateRepo(t *testing.T) {
	repo := parseAuthRepo("myrepo.net/dnephin/foo:tag")
	assert.Check(t, is.Equal(repo, "myrepo.net"))
}

func TestParseAuthRepoPrivateRepoNoUsername(t *testing.T) {
	repo := parseAuthRepo("myrepo.net/foo")
	assert.Check(t, is.Equal(repo, "myrepo.net"))
}
