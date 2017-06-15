package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAuthRepoWithUserRepo(t *testing.T) {
	repo, err := parseAuthRepo("dnephin/foo")
	assert.Nil(t, err)
	assert.Equal(t, repo, defaultRepo)
}

func TestParseAuthRepoPrivateRepoAndPort(t *testing.T) {
	repo, err := parseAuthRepo("myrepo.net:3434/dnephin/foo")
	assert.Nil(t, err)
	assert.Equal(t, repo, "myrepo.net:3434")
}

func TestParseAuthRepoPrivateRepo(t *testing.T) {
	repo, err := parseAuthRepo("myrepo.net/dnephin/foo:tag")
	assert.Nil(t, err)
	assert.Equal(t, repo, "myrepo.net")
}

func TestParseAuthRepoPrivateRepoNoUsername(t *testing.T) {
	repo, err := parseAuthRepo("myrepo.net/foo")
	assert.Nil(t, err)
	assert.Equal(t, repo, "myrepo.net")
}
