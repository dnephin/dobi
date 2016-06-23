package config

import (
	"fmt"
)

// MetaConfig is a data object for non-resource configuration
type MetaConfig struct {
	Default      string
	UniqueExecID string
}

// Validate the MetaConfig
func (m *MetaConfig) Validate(config *Config) error {
	if _, ok := config.Resources[m.Default]; m.Default != "" && !ok {
		return fmt.Errorf("Undefined default resource: %s", m.Default)
	}
	return nil
}

// NewMetaConfig returns a new MetaConfig from a raw config map
func NewMetaConfig(values map[string]interface{}) (*MetaConfig, error) {
	meta := &MetaConfig{}
	return meta, Transform(values, meta)
}
