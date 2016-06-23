package config

import (
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

// Validate does nothing
func (c *AliasConfig) Validate(config *Config) error {
	return nil
}

func (c *AliasConfig) String() string {
	return strings.Join(c.Tasks, ", ")
}

// NewAliasConfig returns a new AliasConfig from a raw config map
func NewAliasConfig(values map[string]interface{}) (*AliasConfig, error) {
	alias := &AliasConfig{}
	return alias, Transform(values, alias)
}
