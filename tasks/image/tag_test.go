package image

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/mock/gomock"
	"github.com/gotestyourself/gotestyourself/assert"
)

func setupMockClient(t *testing.T) (*client.MockDockerClient, func()) {
	mock := gomock.NewController(t)
	mockClient := client.NewMockDockerClient(mock)
	return mockClient, func() { mock.Finish() }
}

func setupCtxAndConfig(
	mockClient *client.MockDockerClient,
) (*context.ExecuteContext, *config.ImageConfig) {
	ctx := &context.ExecuteContext{
		Client:     mockClient,
		WorkingDir: "/dir",
	}
	config := &config.ImageConfig{
		Image: "imagename",
		Tags:  []string{"tag"},
	}
	return ctx, config
}

func TestTagImageNothingToTag(t *testing.T) {
	ctx := &context.ExecuteContext{}
	config := &config.ImageConfig{
		Image: "imagename",
		Tags:  []string{"tag"},
	}
	err := tagImage(ctx, config, "imagename:tag")
	assert.NilError(t, err)
}

func TestTagImageWithTag(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	mockClient.EXPECT().TagImage("imagename:tag", docker.TagImageOptions{
		Repo:  "imagename",
		Tag:   "foo",
		Force: true,
	})

	ctx, config := setupCtxAndConfig(mockClient)
	err := tagImage(ctx, config, "imagename:foo")
	assert.NilError(t, err)
}

func TestTagImageWithFullImageName(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	mockClient.EXPECT().TagImage("imagename:tag", docker.TagImageOptions{
		Repo:  "othername",
		Tag:   "bar",
		Force: true,
	})
	ctx, config := setupCtxAndConfig(mockClient)
	err := tagImage(ctx, config, "othername:bar")
	assert.NilError(t, err)
}

func TestTagImageWithFullImageNameAndHost(t *testing.T) {
	mockClient, teardown := setupMockClient(t)
	defer teardown()
	mockClient.EXPECT().TagImage("imagename:tag", docker.TagImageOptions{
		Repo:  "localhost:3030/othername",
		Tag:   "bar",
		Force: true,
	})
	ctx, config := setupCtxAndConfig(mockClient)
	err := tagImage(ctx, config, "localhost:3030/othername:bar")
	assert.NilError(t, err)
}
