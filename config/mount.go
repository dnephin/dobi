package config

import (
	"fmt"

	"github.com/dnephin/dobi/execenv"
)

// MountConfig is a data object for a container mount
type MountConfig struct {
	Bind     string `config:"required"`
	Path     string `config:"required"`
	ReadOnly bool
}

// Dependencies returns an empty list, Mount resources have no dependencies
func (c *MountConfig) Dependencies() []string {
	return []string{}
}

// Validate checks that all fields have acceptable values
func (c *MountConfig) Validate(path Path, config *Config) *PathError {
	return nil
}

func (c *MountConfig) String() string {
	return fmt.Sprintf("Create directory '%s' to be mounted at '%s'", c.Bind, c.Path)
}

// Resolve resolves variables in the resource
func (c *MountConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	return c, nil
}

func mountFromConfig(name string, values map[string]interface{}) (Resource, error) {
	mount := &MountConfig{}
	return mount, Transform(name, values, mount)
}

func init() {
	RegisterResource("mount", mountFromConfig)
}
