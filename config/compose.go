package config

import (
	"fmt"
	"strings"

	"github.com/dnephin/dobi/execenv"
)

// ComposeConfig is a data object for a task compose
type ComposeConfig struct {
	Files   []string
	Project string `config:"required"`
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
