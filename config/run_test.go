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
	s.run.Command = ShlexSlice{original: "run"}
	s.run.Artifact = "foo"
	s.Equal(s.run.String(), "Run 'run' using the 'builder' image to create 'foo'")
}

func (s *RunConfigSuite) TestValidateMissingUse() {
	s.conf.Resources["example"] = &AliasConfig{}
	s.run.Use = "example"
	err := s.run.Validate(NewPath(""), s.conf)
	s.Error(err)
	s.Contains(err.Error(), "example is not an image resource")
}

func (s *RunConfigSuite) TestValidateMissingMount() {
	s.conf.Resources["one"] = NewImageConfig()
	s.conf.Resources["two"] = NewImageConfig()
	s.conf.Resources["example"] = NewImageConfig()
	s.run.Use = "example"
	s.run.Mounts = []string{"one", "two"}

	err := s.run.Validate(NewPath(""), s.conf)
	s.Error(err)
	s.Contains(err.Error(), "one is not a mount resource")
}

func (s *RunConfigSuite) TestRunFromConfig() {
	values := map[string]interface{}{
		"use":        "image-res",
		"command":    "echo foo",
		"entrypoint": "bash -c",
	}
	res, err := runFromConfig("foo", values)
	run, ok := res.(*RunConfig)

	s.Equal(ok, true)
	s.Nil(err)
	s.Equal(run.Use, "image-res")
	s.Equal(run.Command.Value(), []string{"echo", "foo"})
	s.Equal(run.Entrypoint.Value(), []string{"bash", "-c"})
}
