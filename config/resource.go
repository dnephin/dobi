package config

import (
	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/execenv"
)

// Resource is an interface for each configurable type
type Resource interface {
	Dependencies() []string
	Validate(pth.Path, *Config) *pth.Error
	Resolve(*execenv.ExecEnv) (Resource, error)
	//PrependPath(string)
	Describe() string
	String() string
}

// describebale can be used to provide part of the Resource interface
type describable struct {
	// Description Description of the resource. Adding a description to a
	// resource makes it visible from ``dobi list``.
	Description string
}

func (d *describable) Describe() string {
	return d.Description
}

// dependent can be used to provide part of the Resource interface
type dependent struct {
	// Depends The list of task dependencies.
	// type: list of tasks
	Depends []string
}

// Dependencies returns the list of tasks
func (d *dependent) Dependencies() []string {
	return d.Depends
}
