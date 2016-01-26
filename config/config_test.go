package config

import (
	"testing"

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

func (r StubResource) Validate(config *Config) error {
	return nil
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
