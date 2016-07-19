package config

import (
	"fmt"

	"github.com/dnephin/dobi/execenv"
)

// MountConfig A **mount** resource creates a host bind mount.
// name: mount
// example: A mount named ``source`` that mounts the current host directory as
// ``/app/code`` in the container.
//
// .. code-block:: yaml
//
//     mount=source:
//       bind: .
//       path: /app/code
type MountConfig struct {
	// Bind The host path to create and mount
	Bind string `config:"required"`
	// Path The container path of the mount
	Path string `config:"required"`
	// ReadOnly Set the mount to be read-only
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
