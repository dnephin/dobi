package tform

import (
	"reflect"
	"testing"

	"github.com/dnephin/dobi/config/tform/path"
	"github.com/stretchr/testify/assert"
)

type TransformStub struct {
	First   string
	Renamed string `config:"foo-field"`
	Default string
}

func TestTransformWithNamedField(t *testing.T) {
	raw := map[string]interface{}{
		"first":     "value",
		"foo-field": "bar",
	}
	target := &TransformStub{Default: "default"}
	err := Transform("res", raw, target)
	assert.Nil(t, err)
	assert.Equal(t, *target, TransformStub{
		First:   "value",
		Renamed: "bar",
		Default: "default",
	})
}

func TestTransformUnexpectedKey(t *testing.T) {
	raw := map[string]interface{}{"bogus": "value"}
	target := &TransformStub{}
	err := Transform("res", raw, target)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.bogus: unexpected key")
}

func TestTransformMap(t *testing.T) {
	values := map[interface{}]interface{}{
		"FOO": "foo",
		"BAR": "bar",
	}
	raw := reflect.ValueOf(values)

	mapping := make(map[string]interface{})
	target := reflect.ValueOf(mapping)
	err := transformMap(path.NewPath("."), raw, target)
	assert.Nil(t, err)
	assert.Equal(t, mapping["FOO"], values["FOO"])
	assert.Equal(t, mapping["BAR"], values["BAR"])
}
