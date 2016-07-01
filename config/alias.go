package config

import (
	"fmt"
	"strings"
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
func (c *AliasConfig) Validate(config *Config) error {
	return nil
}

func (c *AliasConfig) String() string {
	return fmt.Sprintf("Run tasks: %v", strings.Join(c.Tasks, ", "))
}

func aliasFromConfig(name string, values map[string]interface{}) (Resource, error) {
	alias := &AliasConfig{}
	return alias, Transform(name, values, alias)
}

func init() {
	RegisterResource("alias", aliasFromConfig)
}
