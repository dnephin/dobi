package config

import (
	"fmt"
	"strings"

	"github.com/dnephin/dobi/execenv"
)

// AliasConfig An **alias** resource is a list of other tasks which will be run
// in the order they are listed.
// name: alias
type AliasConfig struct {
	// Tasks The list of tasks
	// type: list of tasks
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
