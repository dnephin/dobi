package config

import "fmt"

// ResourceError represents an error validating a Resource
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
