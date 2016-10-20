package config

import (
	"testing"

	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/execenv"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
	config *Config
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) SetupTest() {
	s.config = NewConfig()
}

type StubResource struct{}

func (r StubResource) Dependencies() []string {
	return nil
}

func (r StubResource) Validate(path pth.Path, config *Config) *pth.Error {
	return nil
}

func (r StubResource) Resolve(env *execenv.ExecEnv) (Resource, error) {
	return r, nil
}

func (r StubResource) Describe() string {
	return ""
}

func (r StubResource) String() string {
	return ""
}

func (s *ConfigSuite) TestSorted() {
	s.config.Resources = map[string]Resource{
		"beta":  StubResource{},
		"alpha": StubResource{},
		"cabo":  StubResource{},
	}
	sorted := s.config.Sorted()
	s.Equal([]string{"alpha", "beta", "cabo"}, sorted)
}
