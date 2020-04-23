package config

import (
	"reflect"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPathGlobsTransformConfigFromSlice(t *testing.T) {
	globs := PathGlobs{}

	value := []interface{}{"one", "two", "three"}
	err := globs.TransformConfig(reflect.ValueOf(value))
	assert.NilError(t, err)
	assert.DeepEqual(t, []string{"one", "two", "three"}, globs.globs)
}
