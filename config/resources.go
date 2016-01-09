package config

import (
	"fmt"
	"path/filepath"
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
func (c *ImageConfig) Validate() error {
	// TODO: better way to generate consistent config errors
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
	Use        string
	Artifact   string
	Command    string
	Volumes    []string
	Privileged bool
	Depends    []string
}

// TODO: support interactive/tty

// Dependencies returns the list of implicit and explicit dependencies
func (c *CommandConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Volumes, c.Depends...)...)
}

// Validate checks that all fields have acceptable values
func (c *CommandConfig) Validate() error {
	// TODO: validate required fields are set
	return nil
}

func (c *CommandConfig) String() string {
	artifact, command := "", ""
	if c.Artifact != "" {
		artifact = fmt.Sprintf(" to create '%s'", c.Artifact)
	}
	if c.Command != "" {
		command = fmt.Sprintf("'%s' using '", c.Command)
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
func (c *VolumeConfig) Validate() error {
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
