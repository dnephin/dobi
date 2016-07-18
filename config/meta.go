package config

import "fmt"

// MetaConfig **dobi** meta config allows you to configure how **dobi** runs
// name: meta
type MetaConfig struct {
	// Default The name of a resource to run when no resource name is given.
	Default string

	// Project The name of the project. Used to create unique identifiers for
	// image tags and container names.
	// default: *basename of the working directory*
	Project string

	// Include A list of other dobi configuration files to include. Paths are
	// relative to the current working directory.
	// type: list of filepaths
	Include []string

	// ExecID A template value used as part of unique identifiers for image tags
	// and container names. Supports variables.
	// default: ``{env.USER}``
	ExecID string `config:"exec-id"`
}

// Validate the MetaConfig
func (m *MetaConfig) Validate(config *Config) error {
	if _, ok := config.Resources[m.Default]; m.Default != "" && !ok {
		return fmt.Errorf("Undefined default resource: %s", m.Default)
	}
	return nil
}

// IsZero returns true if the struct contains only zero values, except for
// Includes which is ignored
func (m *MetaConfig) IsZero() bool {
	return m.Default == "" && m.Project == "" && m.ExecID == ""
}

// NewMetaConfig returns a new MetaConfig from config values
func NewMetaConfig(name string, values map[string]interface{}) (*MetaConfig, error) {
	meta := &MetaConfig{}
	return meta, Transform(name, values, meta)
}
