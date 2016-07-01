package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

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
	s.conf.Resources["example"] = &AliasConfig{}
	s.run.Use = "example"
	err := s.run.Validate(s.conf)
	s.Error(err)
	s.Contains(err.Error(), "example is not an image resource")
}

func (s *RunConfigSuite) TestValidateMissingVolume() {
	s.conf.Resources["one"] = NewImageConfig()
	s.conf.Resources["two"] = NewImageConfig()
	s.conf.Resources["example"] = NewImageConfig()
	s.run.Use = "example"
	s.run.Volumes = []string{"one", "two"}

	err := s.run.Validate(s.conf)
	s.Error(err)
	s.Contains(err.Error(), "one is not a volume resource")
}
