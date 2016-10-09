package tform

import (
	"fmt"
	"testing"

	pth "github.com/dnephin/dobi/config/tform/path"
	"github.com/stretchr/testify/assert"
)

type Something struct {
	First     string   `config:"required"`
	Second    string   `config:"required,validate"`
	Third     []string `config:"required"`
	Forth     string
	Renamed   string `config:"foo-field,validate"`
	validated bool
}

func (s *Something) ValidateSecond() error {
	if s.Second == "bad" {
		return fmt.Errorf("validation error")
	}
	s.validated = true
	return nil
}

func (s *Something) ValidateRenamed() error {
	if s.Renamed == "bad" {
		return fmt.Errorf("validation error")
	}
	return nil
}

func TestValidateFieldsIsValid(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "b",
		Third:  []string{"one"},
	}
	err := ValidateFields(pth.NewPath(""), obj)
	assert.Nil(t, err)
	assert.Equal(t, obj.validated, true)
}

func TestValidateFieldsMissingRequiredString(t *testing.T) {
	obj := &Something{
		Second: "b",
		Third:  []string{"one"},
	}
	err := ValidateFields(pth.NewPath("res"), obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.first: a value is required")
}

func TestValidateFieldsMissingRequiredStringSlice(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "b",
	}
	err := ValidateFields(pth.NewPath("res"), obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.third: a value is required")
}

func TestValidateFieldsValidationFailed(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "bad",
		Third:  []string{"one"},
	}
	err := ValidateFields(pth.NewPath("res"), obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.second: failed validation: validation error")
}

func TestFieldWithDefinedName(t *testing.T) {
	obj := &Something{
		First:   "a",
		Second:  "ok",
		Third:   []string{"one"},
		Renamed: "bad",
	}
	err := ValidateFields(pth.NewPath("res"), obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.foo-field: failed validation: validation error")
}

type Sample struct {
	ThisNameIsLong string `config:"required"`
}

func TestValidateFieldsCorrectFieldName(t *testing.T) {
	err := ValidateFields(pth.NewPath("res"), &Sample{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.this-name-is-long: a value is required")
}

type BadValidation struct {
	Foo string `config:"validate"`
}

func (t *BadValidation) ValidateFoo() (string, error) {
	return "", nil
}

func TestValidateFieldsBadValidationFuncReturnValue(t *testing.T) {
	err := ValidateFields(pth.NewPath(""), &BadValidation{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be of type \"func() error\"")
}

type BadValidationTwo struct {
	Foo string `config:"validate"`
}

func (t *BadValidationTwo) ValidateFoo(a string) error {
	return nil
}

func TestValidateFieldsBadValidationFuncArgs(t *testing.T) {
	err := ValidateFields(pth.NewPath(""), &BadValidationTwo{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be of type \"func() error\"")
}
