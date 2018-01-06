package context

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
)

func TestGetAuthConfigNoAuthConfig(t *testing.T) {
	context := ExecuteContext{}
	auth := context.GetAuthConfig("https://bogus")
	assert.Check(t, is.Compare(auth, docker.AuthConfiguration{}))
}
