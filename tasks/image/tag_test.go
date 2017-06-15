package image

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TagImageSuite struct {
	suite.Suite
	mock       *gomock.Controller
	mockClient *client.MockDockerClient
	ctx        *context.ExecuteContext
	config     *config.ImageConfig
}

func TestTagImageSuite(t *testing.T) {
	suite.Run(t, new(TagImageSuite))
}

func (s *TagImageSuite) SetupTest() {
	s.mock = gomock.NewController(s.T())
	s.mockClient = client.NewMockDockerClient(s.mock)
	s.ctx = &context.ExecuteContext{
		Client:     s.mockClient,
		WorkingDir: "/dir",
	}
	s.config = &config.ImageConfig{
		Image: "imagename",
		Tags:  []string{"tag"},
	}
}

func (s *TagImageSuite) TearDownTest() {
	s.mock.Finish()
}

func (s *TagImageSuite) TestTagImageNothingToTag() {
	err := tagImage(s.ctx, s.config, "imagename:tag")
	s.Nil(err)
}

func (s *TagImageSuite) TestTagImageWithTag() {
	s.mockClient.EXPECT().TagImage("imagename:tag", docker.TagImageOptions{
		Repo:  "imagename",
		Tag:   "foo",
		Force: true,
	})
	err := tagImage(s.ctx, s.config, "imagename:foo")
	s.Nil(err)
}

func (s *TagImageSuite) TestTagImageWithFullImageName() {
	s.mockClient.EXPECT().TagImage("imagename:tag", docker.TagImageOptions{
		Repo:  "othername",
		Tag:   "bar",
		Force: true,
	})
	err := tagImage(s.ctx, s.config, "othername:bar")
	s.Nil(err)
}

func (s *TagImageSuite) TestTagImageWithFullImageNameAndHost() {
	s.mockClient.EXPECT().TagImage("imagename:tag", docker.TagImageOptions{
		Repo:  "localhost:3030/othername",
		Tag:   "bar",
		Force: true,
	})
	err := tagImage(s.ctx, s.config, "localhost:3030/othername:bar")
	s.Nil(err)
}
