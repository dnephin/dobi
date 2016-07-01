package config

import (
	"fmt"
	"strings"
)

// AliasConfig is a data object for a task alias
type AliasConfig struct {
	Tasks []string
}

// Dependencies returns the list of tasks
func (c *AliasConfig) Dependencies() []string {
	return c.Tasks
}

// Validate the aliased resources exist
func (c *AliasConfig) Validate(config *Config) error {
	if err := ValidateResourcesExist(config, c.Dependencies()); err != nil {
		return NewResourceError(c, err.Error())
	}
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
