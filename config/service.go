package config

import (
	"github.com/dnephin/configtf"

	"fmt"
	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/execenv"
)

// ServiceConfig is a resource
type ServiceConfig struct {
	Replicas int
	JobConfig
	dependent
	describable
}

// NewServiceConfig creates a new ServiceConfig with default values
func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}
func serviceFromConfig(name string, values map[string]interface{}) (Resource, error) {
	cmd := &ServiceConfig{}
	return cmd, configtf.Transform(name, values, cmd)
}

func init() {
	RegisterResource("service", serviceFromConfig)
}

func (c *ServiceConfig) String() string {
	return fmt.Sprintf("Service user '%s' image", c.Use)
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *ServiceConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Depends, c.Mounts...)...)
}

// Validate checks that all fields have acceptable values
func (c *ServiceConfig) Validate(path pth.Path, config *Config) *pth.Error {
	return nil
}

// Resolve resolves variables in the resource
func (c *ServiceConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	conf := *c
	return &conf, nil
}
