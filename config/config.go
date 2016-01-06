package config

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ImageConfig ia a data object for image resource
type ImageConfig struct {
	Image      string
	Dockerfile string
	Context    string
	Args       map[string]string
	Pull       bool
	Tags       []string
	Depends    []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *ImageConfig) Dependencies() []string {
	return c.Depends
}

// Validate checks that all fields have acceptable values
func (c *ImageConfig) Validate() error {
	// TODO: better way to generate consistent config errors
	// TODO: check context directory exists
	// TODO: check dockerfile exists
	// TODO: validate required fields are set
	return nil
}

// NewImageConfig creates an ImageConfig with default values
func NewImageConfig() *ImageConfig {
	return &ImageConfig{
		Context:    ".",
		Dockerfile: "Dockerfile",
		Tags:       []string{"latest"},
	}
}

// CommandConfig is a data object for a command resource
type CommandConfig struct {
	Use        string
	Artifact   string
	Command    string
	Volumes    []string
	Privileged bool
	Depends    []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *CommandConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Volumes, c.Depends...)...)
}

// Validate checks that all fields have acceptable values
func (c *CommandConfig) Validate() error {
	// TODO: validate required fields are set
	return nil
}

// VolumeConfig is a data object for a volume resource
type VolumeConfig struct {
	Path  string
	Mount string
	Mode  string
}

// Dependencies returns an empty list, Volume resources have no dependencies
func (c *VolumeConfig) Dependencies() []string {
	return []string{}
}

// Validate checks that all fields have acceptable values
func (c *VolumeConfig) Validate() error {
	// TODO: validate required fields are set
	return nil
}

// NewVolumeConfig creates a VolumeConfig with default values
func NewVolumeConfig() *VolumeConfig {
	return &VolumeConfig{
		Mode: "rw",
	}
}

// Resource is an interface for each configurable type
type Resource interface {
	Dependencies() []string
	Validate() error
}

// Config is a data object for a full config file
type Config struct {
	Resources map[string]Resource
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
	log.WithFields(log.Fields{"filename": filename}).Info("Configuration loaded")
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
	// TODO: if pull is true, image must be set
	// TODO: run and compose actions are mutually exclusive
	// TODO: compose config and filename are mutually exclusive
	return nil
}
