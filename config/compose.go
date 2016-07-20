package config

import (
	"fmt"
	"strings"

	"github.com/dnephin/dobi/execenv"
)

// ComposeConfig A **compose** resource runs ``docker-compose`` to create an
// isolated environment. The **compose** resource keeps containers running
// until **dobi** exits so the containers can be used by other tasks that depend
// on the **compose** resource, or are listed after it in an `alias`_.
//
// .. note::
//
//     `Docker Compose <https://github.com/docker/compose>`_ must be installed
//     and available in ``$PATH`` to use this resource.
//
// name: compose
// example: Start a Compose environment setting the project name to ``web-devenv``
// and using two Compose files.
//
// .. code-block:: yaml
//
//     compose=devenv:
//         files: [docker-compose.yml, docker-compose-dev.yml]
//         project: 'web-devenv'
//
type ComposeConfig struct {
	// Files The Compose files to use. This field supports :doc:`variables`.
	// type: list of filenames
	Files []string
	// Project The project name used by Compose. This field supports
	// :doc:`variables`.
	Project string `config:"required"`
	// Depends The list of resource dependencies
	// type: list of resource names
	Depends []string
}

// Dependencies returns the list of tasks
func (c *ComposeConfig) Dependencies() []string {
	return c.Depends
}

// Validate the resource
func (c *ComposeConfig) Validate(path Path, config *Config) *PathError {
	return nil
}

func (c *ComposeConfig) String() string {
	return fmt.Sprintf("Run Compose project %q from: %v",
		c.Project, strings.Join(c.Files, ", "))
}

// Resolve resolves variables in the resource
func (c *ComposeConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	var err error
	c.Files, err = env.ResolveSlice(c.Files)
	if err != nil {
		return c, err
	}
	c.Project, err = env.Resolve(c.Project)
	return c, err
}

func composeFromConfig(name string, values map[string]interface{}) (Resource, error) {
	compose := &ComposeConfig{Project: "{unique}"}
	return compose, Transform(name, values, compose)
}

func init() {
	RegisterResource("compose", composeFromConfig)
}
