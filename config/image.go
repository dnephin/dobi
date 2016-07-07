package config

import (
	"fmt"
	"path/filepath"

	"github.com/dnephin/dobi/execenv"
)

// ImageConfig ia a data object for image resource
type ImageConfig struct {
	Image                string `config:"required"`
	Dockerfile           string
	Context              string
	Args                 map[string]string
	PullBaseImageOnBuild bool
	Pull                 string
	Tags                 []string
	Depends              []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *ImageConfig) Dependencies() []string {
	return c.Depends
}

// Validate checks that all fields have acceptable values
func (c *ImageConfig) Validate(path Path, config *Config) *PathError {
	// TODO: validate no tag on image name

	if err := c.validateBuildOrPull(); err != nil {
		return PathErrorf(path, err.Error())
	}
	return nil
}

func (c *ImageConfig) validateBuildOrPull() error {
	if c.Dockerfile == "" && c.Context == "" && c.Pull == "" {
		return fmt.Errorf("one of dockerfile, context, or pull is required")
	}
	switch {
	case c.Dockerfile == "" && c.Context != "":
		c.Dockerfile = "Dockerfile"
	case c.Context == "" && c.Dockerfile != "":
		c.Context = "."
	}
	return nil
}

func (c *ImageConfig) String() string {
	dir := filepath.Join(c.Context, c.Dockerfile)
	return fmt.Sprintf("Build image '%s' from '%s'", c.Image, dir)
}

// Resolve resolves variables in the resource
func (c *ImageConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	var err error
	c.Tags, err = env.ResolveSlice(c.Tags)
	return c, err
}

// NewImageConfig creates a new ImageConfig with default values
func NewImageConfig() *ImageConfig {
	return &ImageConfig{}
}

func imageFromConfig(name string, values map[string]interface{}) (Resource, error) {
	image := NewImageConfig()
	return image, Transform(name, values, image)
}

func init() {
	RegisterResource("image", imageFromConfig)
}
