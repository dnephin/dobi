package config

import (
	"fmt"

	"github.com/dnephin/configtf"
)

// MetaConfig Configure **dobi** and include other config files.
// name: meta
// example: Set the the project name to ``mywebapp`` and run the ``all`` task by
// default.
//
// .. code-block:: yaml
//
//     meta:
//         project: mywebapp
//         default: all
//
type MetaConfig struct {
	// Default The name of a task from the ``dobi.yml`` to run when no
	// task name is specified on the command line.
	Default string

	// Project The name of the project. Used to create unique identifiers for
	// image tags and container names.
	// default: *basename of ``dobi.yml``*
	Project string

	// Include A list of dobi configuration files to include. Paths are
	// relative to the current working directory. Includs can be partial
	// configs that depend on resources in any of the other included files.
	// type: list of file paths or glob patterns
	Include PathGlobs

	// ExecID A template value used as part of unique identifiers for image tags
	// and container names. This field supports :doc:`variables`. This value can
	// be overridden with the ``$DOBI_EXEC_ID`` environment variable.
	// default: ``{user.name}``
	ExecID string `config:"exec-id"`
}

// Validate the MetaConfig
func (m *MetaConfig) Validate(config *Config) error {
	if _, ok := config.Resources[m.Default]; m.Default != "" && !ok {
		return fmt.Errorf("undefined default resource: %s", m.Default)
	}
	if err := m.Include.Validate(); err != nil {
		return fmt.Errorf("invalid include: %s", err)
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
	return meta, configtf.Transform(name, values, meta)
}
