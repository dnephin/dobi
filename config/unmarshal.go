package config

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	typeImage   = "image"
	typeVolume  = "volume"
	typeCommand = "command"
	typeAlias   = "alias"
)

var (
	reservedNames = map[string]bool{
		"autoclean": true,
		"list":      true,
	}
)

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
		if name == "meta" {
			c.Meta, err = NewMetaConfig(value)
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

		resource, err := unmarshalResource(resType, value)
		if err != nil {
			return fmt.Errorf("Invalid config for resource %q:\n%s", name, err.Error())
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

func unmarshalResource(resType string, value map[string]interface{}) (Resource, error) {
	switch resType {
	case typeImage:
		return NewImageConfig(value)
	case typeVolume:
		return NewVolumeConfig(value)
	case typeCommand:
		return NewCommandConfig(value)
	case typeAlias:
		return NewAliasConfig(value)
	default:
		return nil, fmt.Errorf("invalid resource type %q", resType)
	}
}

// LoadFromBytes loads a configuration from a bytes slice
func LoadFromBytes(data []byte) (*Config, error) {
	config := NewConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
