package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	log "github.com/Sirupsen/logrus"
)

// Resource is an interface for each configurable type
type Resource interface {
	Dependencies() []string
	Validate(config *Config) error
}

// MetaConfig is a data object for non-resource configuration
type MetaConfig struct {
	Default string
}

// Validate the MetaConfig
func (m *MetaConfig) Validate(config *Config) error {
	if _, ok := config.Resources[m.Default]; m.Default != "" && !ok {
		return fmt.Errorf("Undefined default resource: %s", m.Default)
	}
	return nil
}

// Config is a data object for a full config file
type Config struct {
	Meta       *MetaConfig
	Resources  map[string]Resource
	WorkingDir string
}

// NewConfig returns a new Config object
func NewConfig() *Config {
	return &Config{
		Resources: make(map[string]Resource),
	}
}

// Sorted returns the list of resource names in alphabetical sort order
func (c *Config) Sorted() []string {
	names := []string{}
	for name := range c.Resources {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (c *Config) missingResources(names []string) []string {
	missing := []string{}
	for _, name := range names {
		if _, ok := c.Resources[name]; !ok {
			missing = append(missing, name)
		}
	}
	return missing
}

// Load a configuration from a filename
func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := LoadFromBytes(data)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"filename": filename}).Debug("Configuration loaded")

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	config.WorkingDir = filepath.Dir(absPath)

	if err = validate(config); err != nil {
		return nil, err
	}
	return config, nil
}

func validate(config *Config) error {
	for _, resource := range config.Resources {
		if err := resource.Validate(config); err != nil {
			return err
		}
	}
	config.Meta.Validate(config)
	return nil
}
