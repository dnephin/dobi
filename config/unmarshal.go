package config

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	// META is the key used for meta config
	META = "meta"
)

var (
	reservedNames = map[string]bool{
		"autoclean": true,
		"list":      true,
		META:        true,
	}

	resourceTypeRegistry = map[string]resourceFactory{}
)

type resourceFactory func(string, map[string]interface{}) (Resource, error)

func isReservedName(name string) bool {
	_, reserved := reservedNames[name]
	return reserved
}

// UnmarshalYAML unmarshals a config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	values := make(map[string]map[string]interface{})

	if err := unmarshal(&values); err != nil {
		// TODO: better error message on unmarshal failure
		return err
	}

	var err error
	for name, value := range values {
		if name == META {
			c.Meta, err = NewMetaConfig(name, value)
			if err != nil {
				return fmt.Errorf("Invalid \"meta\" config: %s", err)
			}
			continue
		}

		resType, resName, err := parseResourceName(name)
		if err != nil {
			return err
		}

		if isReservedName(resName) {
			return fmt.Errorf(
				"Name %q is reserved, please use a different resource name.", resName)
		}

		resource, err := unmarshalResource(name, resType, value)
		if err != nil {
			return fmt.Errorf("Invalid config for resource %q:\n%s", name, err)
		}
		c.Resources[resName] = resource
	}
	return nil
}

func parseResourceName(value string) (string, string, error) {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf(
			"Resource name %q must be in the form \"type=name\"", value)
	}
	return parts[0], parts[1], nil
}

// RegisterResource registers a config type with a function to unmarshal it from
// config values.
func RegisterResource(name string, typeFunc resourceFactory) {
	resourceTypeRegistry[name] = typeFunc
}

func unmarshalResource(name, resType string, value map[string]interface{}) (Resource, error) {
	fromConfigFunc, ok := resourceTypeRegistry[resType]
	if !ok {
		return nil, fmt.Errorf("invalid resource type %q", resType)
	}
	return fromConfigFunc(name, value)
}

// LoadFromBytes loads a configuration from a bytes slice
func LoadFromBytes(data []byte) (*Config, error) {
	config := NewConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
