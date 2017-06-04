package config

import (
	"github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/execenv"
)

// FakeResource is a fake used for testing
type FakeResource struct {
	config.Annotations
	config.Dependent
}

// Validate is a no-op
func (r *FakeResource) Validate(path path.Path, config *config.Config) *path.Error {
	return nil
}

// Resolve is a no-op
func (r *FakeResource) Resolve(env *execenv.ExecEnv) (config.Resource, error) {
	return r, nil
}

func (r *FakeResource) String() string {
	return "The resource string"
}

var _ config.Resource = &FakeResource{}
