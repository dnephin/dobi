package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

// PathGlobs is a list of path globs
type PathGlobs struct {
	globs []string
}

// Validate the globs
func (p *PathGlobs) Validate() error {
	_, err := p.all()
	return err
}

// TransformConfig from raw value to paths
func (p *PathGlobs) TransformConfig(raw reflect.Value) error {
	if !raw.IsValid() {
		return fmt.Errorf("must be a path glob, was undefined")
	}

	switch value := raw.Interface().(type) {
	case string:
		p.globs = []string{value}
	case []string:
		p.globs = value
	default:
		return fmt.Errorf("must be a string or list of strings, not %T", value)
	}
	return nil

}

func (p *PathGlobs) all() ([]string, error) {
	all := []string{}
	for _, glob := range p.globs {
		paths, err := filepath.Glob(glob)
		if err != nil {
			return all, err
		}
		all = append(all, paths...)
	}
	return all, nil
}

// Paths returns all the paths matched by the glob
func (p *PathGlobs) Paths() []string {
	all, err := p.all()
	if err != nil {
		// Error hould have already been returned during Validate()
		panic(err)
	}
	return all
}

// Empty returns true if there are no globs
func (p *PathGlobs) Empty() bool {
	return len(p.globs) == 0
}

func (p *PathGlobs) String() string {
	return strings.Join(p.globs, ", ")
}

// NoMatches returns true if there are globs defined, but none are valid paths
func (p *PathGlobs) NoMatches() bool {
	return !p.Empty() && len(p.Paths()) == 0
}

type validator struct {
	name     string
	validate func() error
}

func newValidator(name string, validate func() error) validator {
	return validator{name: name, validate: validate}
}
