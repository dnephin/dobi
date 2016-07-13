package config

import "fmt"

// MetaConfig is a data object for non-resource configuration
type MetaConfig struct {
	Default string
	Project string
	Include []string
	ExecID  string `config:"exec-id"`
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
