package config

import (
	"fmt"
	"testing"

	"github.com/dnephin/dobi/execenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
	config *Config
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) SetupTest() {
	s.config = NewConfig()
}

type StubResource struct{}

func (r StubResource) Dependencies() []string {
	return nil
}

func (r StubResource) Validate(path Path, config *Config) *PathError {
	return nil
}

func (r StubResource) Resolve(env *execenv.ExecEnv) (Resource, error) {
	return r, nil
}

func (s *ConfigSuite) TestSorted() {
	s.config.Resources = map[string]Resource{
		"beta":  StubResource{},
		"alpha": StubResource{},
		"cabo":  StubResource{},
	}
	sorted := s.config.Sorted()
	s.Equal([]string{"alpha", "beta", "cabo"}, sorted)
}

type Something struct {
	First     string   `config:"required"`
	Second    string   `config:"required,validate"`
	Third     []string `config:"required"`
	Forth     string
	validated bool
}

func (s *Something) ValidateSecond() error {
	if s.Second == "bad" {
		return fmt.Errorf("validation error")
	}
	s.validated = true
	return nil
}

func TestValidateFieldsIsValid(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "b",
		Third:  []string{"one"},
	}
	err := ValidateFields(NewPath(""), obj)
	assert.Nil(t, err)
	assert.Equal(t, obj.validated, true)
}

func TestValidateFieldsMissingRequiredString(t *testing.T) {
	obj := &Something{
		Second: "b",
		Third:  []string{"one"},
	}
	err := ValidateFields(NewPath("res"), obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.first: a value is required")
}

func TestValidateFieldsMissingRequiredStringSlice(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "b",
	}
	err := ValidateFields(NewPath("res"), obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.third: a value is required")
}

func TestValidateFieldsValidationFailed(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "bad",
		Third:  []string{"one"},
	}
	err := ValidateFields(NewPath("res"), obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error at res.second: failed validation: validation error")
}

type Sample struct {
	ThisNameIsLong string `config:"required"`
}

func TestValidateFieldsCorrectFieldName(t *testing.T) {
	err := ValidateFields(NewPath("res"), &Sample{})
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
	err := ValidateFields(NewPath(""), &BadValidation{})
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
	err := ValidateFields(NewPath(""), &BadValidationTwo{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be of type \"func() error\"")
}
