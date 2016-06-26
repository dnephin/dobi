package config

import (
	"fmt"
	"strings"
)

// RunConfig is a data object for a command resource
type RunConfig struct {
	Use         string
	Artifact    string
	Command     string
	Volumes     []string
	Privileged  bool
	Interactive bool
	Depends     []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *RunConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Volumes, c.Depends...)...)
}

// Validate checks that all fields have acceptable values
func (c *RunConfig) Validate(config *Config) error {
	if missing := config.missingResources(c.Depends); len(missing) != 0 {
		reason := fmt.Sprintf("missing dependencies: %s", strings.Join(missing, ", "))
		return NewResourceError(c, reason)
	}
	if err := c.validateUse(config); err != nil {
		return err
	}
	if err := c.validateVolumes(config); err != nil {
		return err
	}

	// TODO: validate required fields are set
	return nil
}

func (c *RunConfig) validateUse(config *Config) error {
	reason := fmt.Sprintf("%s is not an image resource", c.Use)

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

func (c *RunConfig) validateVolumes(config *Config) error {
	for _, volume := range c.Volumes {
		reason := fmt.Sprintf("%s is not a volume resource", volume)

		res, ok := config.Resources[volume]
		if !ok {
			return NewResourceError(c, reason)
		}

		switch res.(type) {
		case *VolumeConfig:
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

func runFromConfig(values map[string]interface{}) (Resource, error) {
	cmd := &RunConfig{}
	return cmd, Transform(values, cmd)
}

func init() {
	RegisterResource("run", runFromConfig)
}
