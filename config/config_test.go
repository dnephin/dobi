package config

import (
	"fmt"
	"testing"

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

func (r StubResource) Validate(config *Config) error {
	return nil
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
	First  string   `config:"required"`
	Second string   `config:"required,validate"`
	Third  []string `config:"required"`
	Forth  string
}

func (s *Something) ValidateSecond() error {
	if s.Second == "bad" {
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
	err := ValidateFields(obj)
	assert.Nil(t, err)
}

func TestValidateFieldsMissingRequiredString(t *testing.T) {
	obj := &Something{
		Second: "b",
		Third:  []string{"one"},
	}
	err := ValidateFields(obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field \"first\" requires a value")
}

func TestValidateFieldsMissingRequiredStringSlice(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "b",
	}
	err := ValidateFields(obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field \"third\" requires a value")
}

func TestValidateFieldsValidationFailed(t *testing.T) {
	obj := &Something{
		First:  "a",
		Second: "bad",
		Third:  []string{"one"},
	}
	err := ValidateFields(obj)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "\"second\" failed validation: validation error")
}

type Sample struct {
	ThisNameIsLong string `config:"required"`
}

func TestValidateFieldsCorrectFieldName(t *testing.T) {
	err := ValidateFields(&Sample{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Field \"this-name-is-long\" requires a value")
}

type BadValidation struct {
	Foo string `config:"validate"`
}

func (t *BadValidation) ValidateFoo() (string, error) {
	return "", nil
}

func TestValidateFieldsBadValidationFuncReturnValue(t *testing.T) {
	err := ValidateFields(&BadValidation{})
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
	err := ValidateFields(&BadValidationTwo{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be of type \"func() error\"")
}
