package config

import (
	"testing"

	"github.com/dnephin/dobi/execenv"
	"github.com/stretchr/testify/assert"
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

func (s *ConfigSuite) TestSorted() {
	s.config.Resources = map[string]Resource{
		"beta":  &ImageConfig{},
		"alpha": &ImageConfig{},
		"cabo":  &ImageConfig{},
	}
	sorted := s.config.Sorted()
	s.Equal([]string{"alpha", "beta", "cabo"}, sorted)
}

func TestResourceResolveDoesNotMutate(t *testing.T) {
	env := execenv.NewExecEnv("execid", "project", ".")

	for name, fromConfigFunc := range resourceTypeRegistry {
		value := make(map[string]interface{})
		resource, err := fromConfigFunc(name, value)
		assert.Nil(t, err)
		resolved, err := resource.Resolve(env)
		assert.Nil(t, err)
		assert.True(t, resource != resolved,
			"Expected different pointers for %q: %p, %p",
			name, resource, resolved)
	}
}
