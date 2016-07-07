package config

import (
	"fmt"

	"github.com/dnephin/dobi/execenv"
	shlex "github.com/kballard/go-shellquote"
)

// RunConfig is a data object for a command resource
type RunConfig struct {
	Use         string `config:"required"`
	Artifact    string
	Command     string `config:"validate"`
	Mounts      []string
	Privileged  bool
	Interactive bool
	Depends     []string
	Env         []string

	parsedCommand []string
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

// ValidateCommand validates the Command field
func (c *RunConfig) ValidateCommand() error {
	if c.Command != "" {
		command, err := shlex.Split(c.Command)
		if err != nil {
			return fmt.Errorf("failed to parse command %q: %s", c.Command, err)
		}
		c.parsedCommand = command
	}
	return nil
}

// ParsedCommand returns the shlex parsed command
func (c *RunConfig) ParsedCommand() []string {
	return c.parsedCommand
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
	if c.Command != "" {
		command = fmt.Sprintf("'%s' using ", c.Command)
	}
	return fmt.Sprintf("Run %sthe '%s' image%s", command, c.Use, artifact)
}

// Resolve resolves variables in the resource
func (c *RunConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	var err error
	c.Env, err = env.ResolveSlice(c.Env)
	return c, err
}

func runFromConfig(name string, values map[string]interface{}) (Resource, error) {
	cmd := &RunConfig{}
	return cmd, Transform(name, values, cmd)
}

func init() {
	RegisterResource("run", runFromConfig)
}
