package config

import (
	"fmt"
	"reflect"

	"github.com/dnephin/dobi/execenv"
	shlex "github.com/kballard/go-shellquote"
)

// RunConfig A **run** resource uses an `image`_ to run a conatiner.
// A **run** resource that doesn't have an **artifact** is never considered
// up-to-date and will always run.  If a run resource has an **artifact**
// the last modified time of that file will be used as the modified time for the
// **run** resource.
//
// The `image`_ specified in **use** and any `mount`_ resources listed in
// **mounts** are automatically added as dependencies and will always be
// created first.
// name: run
// example: Run a container using the ``builder`` image to compile some source
// code to ``./dist/app-binary``.
//
// .. code-block:: yaml
//
//     run=compile:
//         use: builder
//         mounts: [source, dist]
//         artifact: dist/app-binary
//
type RunConfig struct {
	// Use The name of an `image`_ resource. The referenced image is used
	// to created the container for the **run**.
	Use string `config:"required"`
	// Artifact A host path to a file or directory that is the output of this
	// **run**. Paths are relative to the current working directory.
	Artifact string
	// Command The command to run in the container.
	// type: shell quoted string
	// example: ``"bash -c 'echo something'"``
	Command ShlexSlice
	// Entrypoint Override the image entrypoint
	// type: shell quoted string
	Entrypoint ShlexSlice
	// Mounts A list of `mount`_ resources to use when creating the container.
	// type: list of mount resources
	Mounts []string
	// Privileged Gives extended privileges to the container
	Privileged bool
	// Interactive Makes the container interative and enables a tty.
	Interactive bool
	// Depends The list of resources dependencies
	// type: list of resource names
	Depends []string
	// Env Environment variables to pass to the container. This field
	// supports :doc:`variables`.
	// type: list of ``key=value`` strings
	Env []string
	// ProvideDocker Exposes the docker engine to the container by either
	// mounting the unix socket or setting the **DOCKER_HOST** environment
	// variable.
	ProvideDocker bool
	// NetMode The network mode to use
	NetMode string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *RunConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Depends, c.Mounts...)...)
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
