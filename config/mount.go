package config

import (
	"fmt"

	"github.com/dnephin/dobi/execenv"
	"github.com/dnephin/dobi/utils/fs"
)

// MountConfig A **mount** resource creates a host bind mount.
// name: mount
// example: A mount named ``source`` that mounts the current host directory as
// ``/app/code`` in the container.
//
// .. code-block:: yaml
//
//     mount=source:
//         bind: .
//         path: /app/code
//
type MountConfig struct {
	// Bind The host path to create and mount
	Bind string `config:"required"`
	// Path The container path of the mount
	Path string `config:"required"`
	// ReadOnly Set the mount to be read-only
	ReadOnly bool
	// File When true create an empty file instead of a directory
	File bool
	// Mode The file mode to set on the host file or directory when it is
	// created.
	// default: ``0777`` *(for directories)*, ``0644`` *(for files)*
	Mode int `config:"validate"`
}

// Dependencies returns an empty list, Mount resources have no dependencies
func (c *MountConfig) Dependencies() []string {
	return []string{}
}

// Validate checks that all fields have acceptable values
func (c *MountConfig) Validate(path Path, config *Config) *PathError {
	return nil
}

// ValidateMode validates Mode and sets a default
func (c *MountConfig) ValidateMode() error {
	if c.Mode != 0 {
		return nil
	}
	switch c.File {
	case true:
		c.Mode = 0644
	default:
		c.Mode = 0755
	}
	return nil
}

func (c *MountConfig) String() string {
	var filetype string
	switch c.File {
	case true:
		filetype = "file"
	default:
		filetype = "directory"
	}
	return fmt.Sprintf("Create %s %q to be mounted at %q", filetype, c.Bind, c.Path)
}

// Resolve resolves variables in the resource
func (c *MountConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	var err error
	c.Path, err = env.Resolve(c.Path)
	if err != nil {
		return c, err
	}
	c.Bind, err = fs.ExpandUser(c.Bind)
	return c, err
}

func mountFromConfig(name string, values map[string]interface{}) (Resource, error) {
	mount := &MountConfig{}
	return mount, Transform(name, values, mount)
}

func init() {
	RegisterResource("mount", mountFromConfig)
}
