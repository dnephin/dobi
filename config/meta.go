package config

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)

// MetaConfig is a data object for non-resource configuration
type MetaConfig struct {
	Default      string
	Project      string
	UniqueExecID string
}

// Validate the MetaConfig
func (m *MetaConfig) Validate(config *Config) error {
	if _, ok := config.Resources[m.Default]; m.Default != "" && !ok {
		return fmt.Errorf("Undefined default resource: %s", m.Default)
	}

	if m.Project == "" {
		log.Warn(
			"meta.project is not set. Defauling to working directory basename.")
	}
	return nil
}

// NewMetaConfig returns a new MetaConfig from config values
func NewMetaConfig(name string, values map[string]interface{}) (*MetaConfig, error) {
	meta := &MetaConfig{}
	return meta, Transform(name, values, meta)
}
