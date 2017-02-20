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

//func TestIncludeTransformFromConfigWithSting(t *testing.T) {
//	include := Include{}
//	value := "path/to/config"
//	err := include.TransformConfig(reflect.ValueOf(value))
//
//	expected := Include{
//		include: IncludeFile{
//			File:           value,
//			PathRelativity: "project",
//		},
//	}
//
//	assert.Nil(t, err)
//	assert.Equal(t, include, expected)
//}
//
//func TestIncludeTransformFromConfigWithFile(t *testing.T) {
//	include := Include{}
//	value := map[string]interface{}{
//		"file":            "path/to/config",
//		"path-relativity": "file",
//		"namespace": "foo",
//	}
//	err := include.TransformConfig(reflect.ValueOf(value))
//
//	expected := Include{
//		include: IncludeFile{
//			File:           "path/to/config",
//			PathRelativity: "file",
//		},
//		Namespace: "foo",
//	}
//
//	assert.Nil(t, err)
//	assert.Equal(t, include, expected)
//}
