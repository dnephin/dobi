package config

import (
	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/logging"
	"github.com/pkg/errors"
)

// Resource is an interface for each configurable type
type Resource interface {
	Dependencies() []string
	Validate(pth.Path, *Config) *pth.Error
	Resolve(Resolver) (Resource, error)
	Describe() string
	Group() string
	CategoryTags() []string
	String() string
}

// Annotations provides a description and tags to a resource
type Annotations struct {
	// Description of a resource
	// Deprecated use Annotations.Description
	Description string `config:"validate"`
	Annotations AnnotationFields
}

// Describe returns the resource description
func (a *Annotations) Describe() string {
	if a.Annotations.Description != "" {
		return a.Annotations.Description
	}
	// fall back to old deprecated field
	return a.Description
}

// Group returns the group the resource belongs to
func (a *Annotations) Group() string {
	if a.Annotations.Group == "" {
		return "none"
	}
	return a.Annotations.Group
}

// CategoryTags tags returns the list of tags
func (a *Annotations) CategoryTags() []string {
	return a.Annotations.Tags
}

// ValidateDescription prints a warning if set
func (a *Annotations) ValidateDescription() error {
	if a.Description != "" && a.Annotations.Description != "" {
		return errors.Errorf(
			"deprecated description will be ignored in favor of annotations.description")
	}
	if a.Description != "" {
		logging.Log.Warn("description is deprecated. Use annotations.description")
	}
	return nil
}

// AnnotationFields used to annotate a resource
type AnnotationFields struct {
	// Description Description of the resource. Adding a description to a
	// resource makes it visible from ``dobi list``.
	Description string
	// Group Group defines whether the resource can be listed by group.
	// All resources sharing the same group are listed together.
	// Don't set this value if the resource shall be not grouped.
	// The grouped view is visible from ``dobi list -g``.
	Group string
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

// Resolver is an interface for a type that returns values for variables
type Resolver interface {
	Resolve(tmpl string) (string, error)
	ResolveSlice(tmpls []string) ([]string, error)
}
