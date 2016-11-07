package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathGlobsTransformConfigFromSlice(t *testing.T) {
	globs := PathGlobs{}

	value := []interface{}{"one", "two", "three"}
	err := globs.TransformConfig(reflect.ValueOf(value))
	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two", "three"}, globs.globs)
}
