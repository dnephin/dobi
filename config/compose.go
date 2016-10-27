package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dnephin/configtf"
	pth "github.com/dnephin/configtf/path"
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
	// StopGrace Seconds to wait for containers to stop before killing them.
	// default: ``5``
	StopGrace int
	dependent
	describable
}

// StopGraceString returns StopGrace as a string
func (c *ComposeConfig) StopGraceString() string {
	return strconv.FormatInt(int64(c.StopGrace), 10)
}

// Validate the resource
func (c *ComposeConfig) Validate(path pth.Path, config *Config) *pth.Error {
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
	compose := &ComposeConfig{Project: "{unique}", StopGrace: 5}
	return compose, configtf.Transform(name, values, compose)
}

func init() {
	RegisterResource("compose", composeFromConfig)
}
