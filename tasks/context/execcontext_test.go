package context

import (
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthConfigNoAuthConfig(t *testing.T) {
	context := ExecuteContext{}
	auth := context.GetAuthConfig("https://bogus")
	assert.Equal(t, auth, docker.AuthConfiguration{})
}
