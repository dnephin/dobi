package config

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Resource is an interface for each configurable type
type Resource interface {
	Dependencies() []string
	Validate(config *Config) error
}

// Config is a data object for a full config file
type Config struct {
	Resources  map[string]Resource
	WorkingDir string
}

// NewConfig returns a new Config object
func NewConfig() *Config {
	return &Config{
		Resources: make(map[string]Resource),
	}
}

// UnmarshalYAML unmarshals a config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := newRawMap()
	if err := unmarshal(raw.value); err != nil {
		return err
	}
	for name, raw := range raw.value {
		c.Resources[name] = raw.resource
	}
	return nil
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

// LoadFromBytes loads a configuration from a bytes slice
func LoadFromBytes(data []byte) (*Config, error) {
	config := NewConfig()
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func validate(config *Config) error {
	for _, resource := range config.Resources {
		if err := resource.Validate(config); err != nil {
			return err
		}
	}

	// TODO: validate references between resources
	return nil
}
