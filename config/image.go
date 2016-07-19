package config

import (
	"fmt"
	"path/filepath"

	"github.com/dnephin/dobi/execenv"
)

// ImageConfig An **image** resource provides actions for working with a Docker
// image. If an image is buildable it is considered up-to-date if all files in
// the build context have a modified time older than the created time of the
// image.
// name: image
type ImageConfig struct {
	// Image The name of the **image** without any tags
	Image string `config:"required"`
	// Dockerfile The path to the ``Dockerfile`` used to build the image. This
	// path is relative to the **context**.
	Dockerfile string
	// Context The build context used to build the image.
	// default: ``.``
	Context string
	// Args Build args used to build the image.
	// type: mapping ``key: value``
	Args map[string]string
	// PullBaseImageOnBuild If **true** the base image used in the
	// ``Dockerfile`` will be pulled before building the image.
	PullBaseImageOnBuild bool
	// Pull Not implemented yet
	Pull string
	// Tags The image tags applied to the image before pushing the image to a
	// registry.  The first tag in the list is used when the image is built.
	// Each item in the list supports :doc:`variables`.
	// default: ``['{unique}']``
	// type: list of tags
	Tags []string
	// Depends The list of resource dependencies
	// type: list of resources
	Depends []string
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
