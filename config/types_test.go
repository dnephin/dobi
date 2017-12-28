package config

import (
	"reflect"
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
)

func TestPathGlobsTransformConfigFromSlice(t *testing.T) {
	globs := PathGlobs{}

	value := []interface{}{"one", "two", "three"}
	err := globs.TransformConfig(reflect.ValueOf(value))
	assert.Check(t, is.Nil(err))
	assert.Check(t, is.Compare([]string{"one", "two", "three"}, globs.globs))
}
