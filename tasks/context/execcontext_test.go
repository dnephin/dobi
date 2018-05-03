package context

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestGetAuthConfigNoAuthConfig(t *testing.T) {
	context := ExecuteContext{}
	auth := context.GetAuthConfig("https://bogus")
	assert.Check(t, is.DeepEqual(auth, docker.AuthConfiguration{}))
}
