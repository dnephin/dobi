package config

import (
	"fmt"

	"github.com/dnephin/configtf"
	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/utils/fs"
)

// MountConfig A **mount** resource creates a host bind mount or named volume
// mount.
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
//     mount=named:
//         name: app-data
//         path: /data
//
type MountConfig struct {
	// Bind The host path to create and mount. This field supports expansion of
	// `~` to the current users home directory.
	Bind string
	// Path The container path of the mount
	Path string `config:"required"`
	// Name The name of a named volume
	Name string
	// ReadOnly Set the mount to be read-only
	ReadOnly bool
	// File When true create an empty file instead of a directory
	File bool
	// Mode The file mode to set on the host file or directory when it is
	// created.
	// default: ``0755`` *(for directories)*, ``0644`` *(for files)*
	Mode int `config:"validate"`
	Annotations
}

// Dependencies returns an empty list, Mount resources have no dependencies
func (c *MountConfig) Dependencies() []string {
	return []string{}
}

// Validate checks that all fields have acceptable values
func (c *MountConfig) Validate(path pth.Path, config *Config) *pth.Error {
	switch {
	case c.Bind != "" && c.Name != "":
		return pth.Errorf(path, "\"name\" and \"bind\" can not be used together")
	case c.Bind == "" && c.Name == "":
		return pth.Errorf(path, "One of \"name\" or \"bind\" must be set")
	case c.Name != "" && c.Mode != 0:
		return pth.Errorf(path, "\"mode\" can not be used with named volumes")
	case c.Name != "" && c.File:
		return pth.Errorf(path, "\"file\" can not be used with named volumes")
	}
	return nil
}

// ValidateMode validates Mode and sets a default
func (c *MountConfig) ValidateMode() error {
	if c.Mode != 0 || c.Name != "" {
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
	var mount string
	switch {
	case c.File:
		mount = fmt.Sprintf("file %q", c.Bind)
	case c.Name != "":
		mount = "named volume"
	default:
		mount = fmt.Sprintf("directory %q", c.Bind)
	}
	return fmt.Sprintf("Create %s to be mounted at %q", mount, c.Path)
}

// IsBind returns true if the mount is a bind mount to a host directory
func (c *MountConfig) IsBind() bool {
	return c.Bind != ""
}

// Resolve resolves variables in the resource
func (c *MountConfig) Resolve(resolver Resolver) (Resource, error) {
	conf := *c
	var err error
	conf.Path, err = resolver.Resolve(c.Path)
	if err != nil {
		return &conf, err
	}
	conf.Name, err = resolver.Resolve(c.Name)
	if err != nil {
		return &conf, err
	}
	bind, err := resolver.Resolve(c.Bind)
	if err != nil {
		return &conf, err
	}
	conf.Bind, err = fs.ExpandUser(bind)
	return &conf, err
}

func mountFromConfig(name string, values map[string]interface{}) (Resource, error) {
	mount := &MountConfig{}
	return mount, configtf.Transform(name, values, mount)
}

func init() {
	RegisterResource("mount", mountFromConfig)
}
