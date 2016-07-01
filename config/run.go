package config

import (
	"fmt"

	shellquote "github.com/kballard/go-shellquote"
)

// RunConfig is a data object for a command resource
type RunConfig struct {
	Use           string
	Artifact      string
	Command       string
	Mounts        []string
	Privileged    bool
	Interactive   bool
	Depends       []string
	parsedCommand []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *RunConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Mounts, c.Depends...)...)
}

// Validate checks that all fields have acceptable values
func (c *RunConfig) Validate(config *Config) error {
	if err := ValidateResourcesExist(config, c.Dependencies()); err != nil {
		return NewResourceError(c, err.Error())
	}
	if err := c.validateUse(config); err != nil {
		return err
	}
	if err := c.validateMounts(config); err != nil {
		return err
	}

	if c.Command != "" {
		command, err := shellquote.Split(c.Command)
		if err != nil {
			return NewResourceError(c, "Failed to parse command: %s", err)
		}
		c.parsedCommand = command

	}
	// TODO: validate required fields are set
	return nil
}

// ParsedCommand returns the shlex parsed command
func (c *RunConfig) ParsedCommand() []string {
	return c.parsedCommand
}

func (c *RunConfig) validateUse(config *Config) error {
	reason := fmt.Sprintf("%s is not an image resource", c.Use)

	if c.Use == "" {
		return NewResourceError(c, "\"use\" is required")
	}

	res, ok := config.Resources[c.Use]
	if !ok {
		return NewResourceError(c, reason)
	}

	switch res.(type) {
	case *ImageConfig:
	default:
		return NewResourceError(c, reason)
	}

	return nil
}

func (c *RunConfig) validateMounts(config *Config) error {
	for _, mount := range c.Mounts {
		reason := fmt.Sprintf("%s is not a mount resource", mount)

		res, ok := config.Resources[mount]
		if !ok {
			return NewResourceError(c, reason)
		}

		switch res.(type) {
		case *MountConfig:
		default:
			return NewResourceError(c, reason)
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

func runFromConfig(name string, values map[string]interface{}) (Resource, error) {
	cmd := &RunConfig{}
	return cmd, Transform(name, values, cmd)
}

func init() {
	RegisterResource("run", runFromConfig)
}
