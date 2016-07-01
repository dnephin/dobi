package config

import (
	"fmt"
	"path/filepath"
)

// ImageConfig ia a data object for image resource
type ImageConfig struct {
	Image      string `config:"required"`
	Dockerfile string `config:"required"`
	Context    string `config:"required"`
	Args       map[string]string
	Pull       bool
	Tags       []string
	Depends    []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *ImageConfig) Dependencies() []string {
	return c.Depends
}

// Validate checks that all fields have acceptable values
func (c *ImageConfig) Validate(path Path, config *Config) *PathError {
	// TODO: validate no tag on image name
	return nil
}

func (c *ImageConfig) String() string {
	dir := filepath.Join(c.Context, c.Dockerfile)
	return fmt.Sprintf("Build image '%s' from '%s'", c.Image, dir)
}

// NewImageConfig creates a new ImageConfig with default values
func NewImageConfig() *ImageConfig {
	return &ImageConfig{Dockerfile: "Dockerfile", Context: "."}
}

func imageFromConfig(name string, values map[string]interface{}) (Resource, error) {
	image := NewImageConfig()
	return image, Transform(name, values, image)
}

func init() {
	RegisterResource("image", imageFromConfig)
}
