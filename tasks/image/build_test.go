package image

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ImageRecordSuite struct {
	suite.Suite
	mock   *gomock.Controller
	ctx    *context.ExecuteContext
	config *config.ImageConfig
}

func TestImageRecordSuite(t *testing.T) {
	suite.Run(t, new(ImageRecordSuite))
}

func (s *ImageRecordSuite) SetupTest() {
	s.mock = gomock.NewController(s.T())
	s.ctx = &context.ExecuteContext{
		Client:     client.NewMockDockerClient(s.mock),
		WorkingDir: "/dir",
	}
	s.config = &config.ImageConfig{
		Image: "imagename",
		Tags:  []string{"tag"},
	}
}

func (s *ImageRecordSuite) TearDownTest() {
	s.mock.Finish()
}

func (s *ImageRecordSuite) TestRecordPath() {
	path := recordPath(s.ctx, s.config)
	s.Equal("/dir/.dobi/images/imagename:tag", path)
}

func (s *ImageRecordSuite) TestRecordPathEscapesSlash() {
	s.config.Image = "repo/name"
	path := recordPath(s.ctx, s.config)
	s.Equal("/dir/.dobi/images/repo name:tag", path)
}
