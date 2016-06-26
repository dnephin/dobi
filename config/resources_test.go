package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ImageConfigSuite struct {
	suite.Suite
	image *ImageConfig
}

func TestImageConfigSuite(t *testing.T) {
	suite.Run(t, new(ImageConfigSuite))
}

func (s *ImageConfigSuite) SetupTest() {
	s.image = NewImageConfig()
	s.image.Image = "example"
}

func (s *ImageConfigSuite) TestString() {
	s.image.Context = "./files"
	s.Equal("Build image 'example' from 'files/Dockerfile'", s.image.String())
}

func (s *ImageConfigSuite) TestValidateMissingDependencies() {
	s.image.Depends = []string{"one", "two"}
	conf := NewConfig()
	err := s.image.Validate(conf)
	s.Error(err)
	s.Contains(err.Error(), "missing dependencies: one, two")
}

type RunConfigSuite struct {
	suite.Suite
	run  *RunConfig
	conf *Config
}

func TestRunConfigSuite(t *testing.T) {
	suite.Run(t, new(RunConfigSuite))
}

func (s *RunConfigSuite) SetupTest() {
	s.run = &RunConfig{}
	s.conf = NewConfig()
}

func (s *RunConfigSuite) TestString() {
	s.run.Use = "builder"
	s.run.Command = "run"
	s.run.Artifact = "foo"
	s.Equal(s.run.String(), "Run 'run' using the 'builder' image to create 'foo'")
}

func (s *RunConfigSuite) TestValidateMissingUse() {
	s.run.Use = "example"
	err := s.run.Validate(s.conf)
	s.Error(err)
	s.Contains(err.Error(), "example is not an image resource")
}

func (s *RunConfigSuite) TestValidateMissingVolume() {
	s.run.Use = "example"
	s.run.Volumes = []string{"one", "two"}

	s.conf.Resources["example"] = &ImageConfig{}

	err := s.run.Validate(s.conf)
	s.Error(err)
	s.Contains(err.Error(), "one is not a volume resource")
}
