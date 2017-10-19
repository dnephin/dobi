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
	Describe() string
	CategoryTags() []string
	String() string
}

// Annotations provides a description and tags to a resource
type Annotations struct {
	Annotations AnnotationFields
}

// Describe returns the resource description
func (a *Annotations) Describe() string {
	return a.Annotations.Description
}

// CategoryTags tags returns the list of tags
func (a *Annotations) CategoryTags() []string {
	return a.Annotations.Tags
}

// AnnotationFields used to annotate a resource
type AnnotationFields struct {
	// Description Description of the resource. Adding a description to a
	// resource makes it visible from ``dobi list``.
	Description string
	// Tags
	Tags []string
}

// Dependent can be used to provide part of the Resource interface
type Dependent struct {
	// Depends The list of task dependencies.
	// type: list of tasks
	Depends []string
}

// Dependencies returns the list of tasks
func (d *Dependent) Dependencies() []string {
	return d.Depends
}
