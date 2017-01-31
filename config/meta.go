package config

import (
	"fmt"
	"reflect"

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
	// type: list of file path globs, or include objects
	Include []Include

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
	for _, include := range m.Include {
		if err := include.Validate(); err != nil {
			return fmt.Errorf("invalid include: %s", err)
		}
	}
	return nil
}

type includeable interface {
	Load(path string) (*Config, error)
}

// Include is either a filepath glob or url to a dobi config file
type Include struct {
	include includeable
}

// TransformConfig from raw value to an include object
func (i *Include) TransformConfig(raw reflect.Value) error {
	if !raw.IsValid() {
		return fmt.Errorf("must be a include, was undefined")
	}

	switch value := raw.Interface().(type) {
	case string:
		i.include = includeFile{File: value, Relativity: "project"}
		return nil
	case map[string]interface{}:
		if _, ok := value["file"]; ok {
			include := includeFile{}
			return configtf.Transform("meta.include", value, include)
		}
		if _, ok := value["url"]; ok {
			return fmt.Errorf("url includes not yet implemented")
		}
	}
	return fmt.Errorf("must be a string or list of strings, not %T", raw.Interface())
}

// Validate the include
func (i *Include) Validate() error {
	return nil
}

// Load configuration for the include
func (i *Include) Load(path string) (*Config, error) {
	return nil, nil
}

type includeFile struct {
	File       string
	Relativity string
	Optional   bool
}

func (f includeFile) Load(path string) (*Config, error) {
	return nil, nil
}

// NewMetaConfig returns a new MetaConfig from config values
func NewMetaConfig(name string, values map[string]interface{}) (*MetaConfig, error) {
	meta := &MetaConfig{}
	return meta, configtf.Transform(name, values, meta)
}
