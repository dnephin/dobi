package config

import (
	"fmt"
	"strings"
)

// ResourceError represents an error validating a Resource
// TODO: delete this error
type ResourceError struct {
	Resource Resource
	Reason   string
}

func (e ResourceError) Error() string {
	return e.Reason
}

// NewResourceError returns a new error for a reasource
func NewResourceError(resource Resource, reason string, args ...interface{}) *ResourceError {
	return &ResourceError{
		Resource: resource,
		Reason:   fmt.Sprintf(reason, args...),
	}
}

// Path is a dotted path of key names to config values
type Path struct {
	path []string
}

func (p *Path) add(next string) Path {
	return Path{path: append(p.path, next)}
}

// Path returns the config keys in the path
func (p *Path) Path() []string {
	return p.path
}

func (p *Path) String() string {
	return strings.Join(p.path, ".")
}

// NewPath returns a new root Path
func NewPath(root string) Path {
	return Path{path: []string{root}}
}

// PathError is an error during config transformation or validation
type PathError struct {
	path Path
	msg  string
}

func (e *PathError) Error() string {
	return fmt.Sprintf("Error at %s: %s", e.path.String(), e.msg)
}

// Path returns the config path where the error occurred
func (e *PathError) Path() Path {
	return e.path
}

// PathErrorf returns a new PathError with a formatted message
func PathErrorf(path Path, msg string, args ...interface{}) *PathError {
	return &PathError{path: path, msg: fmt.Sprintf(msg, args...)}
}
