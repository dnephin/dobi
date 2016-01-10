package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ImageConfig ia a data object for image resource
type ImageConfig struct {
	Image      string
	Dockerfile string
	Context    string
	Args       map[string]string
	Pull       bool
	Depends    []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *ImageConfig) Dependencies() []string {
	return c.Depends
}

// Validate checks that all fields have acceptable values
func (c *ImageConfig) Validate(config *Config) error {
	if missing := config.missingResources(c.Depends); len(missing) != 0 {
		reason := fmt.Sprintf("missing dependencies: %s", strings.Join(missing, ", "))
		return NewResourceError(c, reason)
	}

	// TODO: check context directory exists
	// TODO: check dockerfile exists
	// TODO: validate required fields are set
	// TODO: validate no tag on image name
	return nil
}

func (c *ImageConfig) String() string {
	dir := filepath.Join(c.Context, c.Dockerfile)
	return fmt.Sprintf("Build image '%s' from '%s'", c.Image, dir)
}

// NewImageConfig creates an ImageConfig with default values
func NewImageConfig() *ImageConfig {
	return &ImageConfig{
		Context:    ".",
		Dockerfile: "Dockerfile",
	}
}

// CommandConfig is a data object for a command resource
type CommandConfig struct {
	Use         string
	Artifact    string
	Command     string
	Volumes     []string
	Privileged  bool
	Interactive bool
	Depends     []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *CommandConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Volumes, c.Depends...)...)
}

// Validate checks that all fields have acceptable values
func (c *CommandConfig) Validate(config *Config) error {
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

func (c *CommandConfig) validateUse(config *Config) error {
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

func (c *CommandConfig) validateVolumes(config *Config) error {
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

func (c *CommandConfig) String() string {
	artifact, command := "", ""
	if c.Artifact != "" {
		artifact = fmt.Sprintf(" to create '%s'", c.Artifact)
	}
	if c.Command != "" {
		command = fmt.Sprintf("'%s' using ", c.Command)
	}
	return fmt.Sprintf("Run %sthe '%s' image%s", command, c.Use, artifact)
}

// VolumeConfig is a data object for a volume resource
type VolumeConfig struct {
	Path  string
	Mount string
	Mode  string
}

// Dependencies returns an empty list, Volume resources have no dependencies
func (c *VolumeConfig) Dependencies() []string {
	return []string{}
}

// Validate checks that all fields have acceptable values
func (c *VolumeConfig) Validate(config *Config) error {
	// TODO: validate required fields are set
	return nil
}

func (c *VolumeConfig) String() string {
	return fmt.Sprintf("Create volume '%s' to be mounted at '%s'", c.Path, c.Mount)
}

// NewVolumeConfig creates a VolumeConfig with default values
func NewVolumeConfig() *VolumeConfig {
	return &VolumeConfig{
		Mode: "rw",
	}
}
