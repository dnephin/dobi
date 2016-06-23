package config

import (
	"fmt"
)

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
func (c *VolumeConfig) Validate(config *Config) error {
	// TODO: validate required fields are set
	return nil
}

func (c *VolumeConfig) String() string {
	return fmt.Sprintf("Create volume '%s' to be mounted at '%s'", c.Path, c.Mount)
}

// NewVolumeConfig creates a VolumeConfig from a raw config map
func NewVolumeConfig(values map[string]interface{}) (*VolumeConfig, error) {
	volume := &VolumeConfig{Mode: "rw"}
	return volume, Transform(values, volume)
}
