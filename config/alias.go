package config

import (
	"fmt"
	"strings"

	"github.com/dnephin/dobi/execenv"
)

// AliasConfig is a data object for a task alias
type AliasConfig struct {
	Tasks []string `config:"required"`
}

// Dependencies returns the list of tasks
func (c *AliasConfig) Dependencies() []string {
	return c.Tasks
}

// Validate the resource
func (c *AliasConfig) Validate(path Path, config *Config) *PathError {
	return nil
}

func (c *AliasConfig) String() string {
	return fmt.Sprintf("Run tasks: %v", strings.Join(c.Tasks, ", "))
}

// Resolve resolves variables in the resource
func (c *AliasConfig) Resolve(env *execenv.ExecEnv) (Resource, error) {
	return c, nil
}

func aliasFromConfig(name string, values map[string]interface{}) (Resource, error) {
	alias := &AliasConfig{}
	return alias, Transform(name, values, alias)
}

func init() {
	RegisterResource("alias", aliasFromConfig)
}
