package path

import (
	"fmt"
	"strings"
)

// Path is a dotted path of key names to config values
type Path struct {
	path []string
}

// Add create a new path from the current path, by adding the next path segment
func (p *Path) Add(next string) Path {
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

// Error is an error during config transformation or validation
type Error struct {
	path Path
	msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error at %s: %s", e.path.String(), e.msg)
}

// Path returns the config path where the error occurred
func (e *Error) Path() Path {
	return e.path
}

// Errorf returns a new Error with a formatted message
func Errorf(path Path, msg string, args ...interface{}) *Error {
	return &Error{path: path, msg: fmt.Sprintf(msg, args...)}
}
