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

type CommandConfigSuite struct {
	suite.Suite
	command *CommandConfig
	conf    *Config
}

func TestCommandConfigSuite(t *testing.T) {
	suite.Run(t, new(CommandConfigSuite))
}

func (s *CommandConfigSuite) SetupTest() {
	s.command = &CommandConfig{}
	s.conf = NewConfig()
}

func (s *CommandConfigSuite) TestString() {
	s.command.Use = "builder"
	s.command.Command = "run"
	s.command.Artifact = "foo"
	s.Equal(s.command.String(), "Run 'run' using the 'builder' image to create 'foo'")
}

func (s *CommandConfigSuite) TestValidateMissingUse() {
	s.command.Use = "example"
	err := s.command.Validate(s.conf)
	s.Error(err)
	s.Contains(err.Error(), "example is not an image resource")
}

func (s *CommandConfigSuite) TestValidateMissingVolume() {
	s.command.Use = "example"
	s.command.Volumes = []string{"one", "two"}

	s.conf.Resources["example"] = &ImageConfig{}

	err := s.command.Validate(s.conf)
	s.Error(err)
	s.Contains(err.Error(), "one is not a volume resource")
}
