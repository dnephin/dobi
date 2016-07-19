package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Something an example config struct used for unmarshaling and validating a
// configuration from a file.
// name: special
type Something struct {
	// First the first field in the struct
	// example: "foo"
	// format: free text
	First string `config:"required"`
	// Second the number of things in a something
	// default: 66
	Second int `config:"foo-field"`
	// Third a list of items
	// type: array
	Third []string
	Forth string
	// Fifth still a thing
	Fifth int
}

func TestParse(t *testing.T) {
	config, err := Parse(Something{})

	expected := ConfigType{
		Name:        "special",
		Description: "an example config struct used for unmarshaling and validating a configuration from a file. ",
		Fields: []ConfigField{
			{
				Name:        "first",
				IsRequired:  true,
				Type:        "string",
				Format:      "free text",
				Example:     "\"foo\"",
				Description: "the first field in the struct ",
			},
			{
				Name:        "foo-field",
				Type:        "int",
				Default:     "66",
				Description: "the number of things in a something ",
			},
			{
				Name:        "third",
				Type:        "array",
				Description: "a list of items ",
			},
			{
				Name:        "forth",
				Type:        "string",
				Description: " ",
			},
			{
				Name:        "fifth",
				Type:        "int",
				Description: "still a thing \n\n",
			},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, config)
}
