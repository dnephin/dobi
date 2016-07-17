package config

import (
	"fmt"
	"reflect"

	"github.com/dnephin/dobi/execenv"
	shlex "github.com/kballard/go-shellquote"
)

// RunConfig is a data object for a command resource
type RunConfig struct {
	Use           string `config:"required"`
	Artifact      string
	Command       ShlexSlice
	Mounts        []string
	Privileged    bool
	Interactive   bool
	Depends       []string
	Env           []string
	Entrypoint    ShlexSlice
	ProvideDocker bool
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *RunConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Mounts, c.Depends...)...)
}

// Validate checks that all fields have acceptable values
func (c *RunConfig) Validate(path Path, config *Config) *PathError {
	if err := c.validateUse(config); err != nil {
		return PathErrorf(path.add("use"), err.Error())
	}
	if err := c.validateMounts(config); err != nil {
		return PathErrorf(path.add("mounts"), err.Error())
	}
	return nil
}

func (c *RunConfig) validateUse(config *Config) error {
	err := fmt.Errorf("%s is not an image resource", c.Use)

	res, ok := config.Resources[c.Use]
	if !ok {
		return err
	}

	switch res.(type) {
	case *ImageConfig:
	default:
		return err
	}

	return nil
}

func (c *RunConfig) validateMounts(config *Config) error {
	for _, mount := range c.Mounts {
		err := fmt.Errorf("%s is not a mount resource", mount)

		res, ok := config.Resources[mount]
		if !ok {
			return err
		}

		switch res.(type) {
		case *MountConfig:
		default:
			return err
		}
	}
	return nil
}

func (c *RunConfig) String() string {
	artifact, command := "", ""
	if c.Artifact != "" {
		artifact = fmt.Sprintf(" to create '%s'", c.Artifact)
	}
	// TODO: look for entrypoint as well as command
	if !c.Command.Empty() {
		command = fmt.Sprintf("'%s' using ", c.Command.String())
	}
	return fmt.Sprintf("Run %sthe '%s' image%s", command, c.Use, artifact)
}

// Resolve resolves variables in the resource
func (c *RunConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	var err error
	c.Env, err = env.ResolveSlice(c.Env)
	return c, err
}

// ShlexSlice is a type used for config transforming a string into a []string
// using shelx.
type ShlexSlice struct {
	original string
	parsed   []string
}

func (s *ShlexSlice) String() string {
	return s.original
}

// Value returns the slice value
func (s *ShlexSlice) Value() []string {
	return s.parsed
}

// Empty returns true if the instance contains the zero value
func (s *ShlexSlice) Empty() bool {
	return s.original == ""
}

// TransformConfig is used to transform a string from a config file into a
// sliced value, using shlex.
func (s *ShlexSlice) TransformConfig(raw reflect.Value) error {
	var err error
	switch value := raw.Interface().(type) {
	case string:
		s.original = value
		s.parsed, err = shlex.Split(value)
		if err != nil {
			return fmt.Errorf("failed to parse command %q: %s", value, err)
		}
	default:
		return fmt.Errorf("must be a string, not %T", value)
	}
	return nil
}

func runFromConfig(name string, values map[string]interface{}) (Resource, error) {
	cmd := &RunConfig{}
	return cmd, Transform(name, values, cmd)
}

func init() {
	RegisterResource("run", runFromConfig)
}
