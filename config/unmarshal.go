package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

var (
	reservedNames = map[string]bool{
		"autoclean": true,
	}
)

type rawMap struct {
	values map[string]stringKeyMap
	meta   *MetaConfig
}

type stringKeyMap struct {
	value    map[string]interface{}
	resource Resource
}

type rawMeta struct {
	Meta *MetaConfig
}

func (m *rawMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	values := make(map[string]interface{})

	if err := unmarshal(values); err != nil {
		return err
	}

	for name := range values {
		if isReservedName(name) {
			return fmt.Errorf(
				"Name '%s' is reserved, please use a different resource name.",
				name)
		}
	}

	metaMap := &rawMeta{}
	if err := unmarshal(metaMap); err != nil {
		return err
	}
	m.meta = metaMap.Meta

	// TODO: there must be a better way to do this without re-parsing the config
	// multiple times.
	return unmarshal(m.values)
}

func newRawMap() *rawMap {
	return &rawMap{
		values: make(map[string]stringKeyMap),
	}
}

func isReservedName(name string) bool {
	_, reserved := reservedNames[name]
	return reserved
}

// UnmarshalYAML unmarshals a raw config resource
func (m *stringKeyMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m.value = make(map[string]interface{})
	err := unmarshal(m.value)
	if err != nil {
		return err
	}

	var conf Resource
	switch {
	case m.hasKeys("path", "mount"):
		conf = NewVolumeConfig()
	case m.hasKeys("image"):
		conf = NewImageConfig()
	case m.hasKeys("use"):
		conf = &CommandConfig{}
	case m.hasKeys("tasks"):
		conf = &AliasConfig{}
	default:
		// TODO: error on unknown resource type unless it's the MetaConfig
		return nil
	}

	// TODO: error on unexpected fields
	err = unmarshal(conf)
	m.resource = conf
	return err
}

func (m *stringKeyMap) hasKeys(keys ...string) bool {
	for _, key := range keys {
		if _, ok := m.value[key]; ok {
			return true
		}
	}
	return false
}

// LoadFromBytes loads a configuration from a bytes slice
func LoadFromBytes(data []byte) (*Config, error) {
	rawMap := newRawMap()
	if err := yaml.Unmarshal(data, rawMap); err != nil {
		return nil, err
	}

	config := NewConfig()
	for name, raw := range rawMap.values {
		// TODO: this is only necessary because of meta
		if raw.resource == nil {
			continue
		}
		config.Resources[name] = raw.resource
	}
	if rawMap.meta != nil {
		config.Meta = rawMap.meta
	} else {
		config.Meta = &MetaConfig{}
	}
	return config, nil
}
