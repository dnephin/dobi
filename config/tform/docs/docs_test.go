package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Something an example config struct used for unmarshaling and validating a
// configuration from a file.
//
// New Paragraph
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
	// example:
	//   {
	//     third: ['one', 'two'],
	//   }
	Third []string

	Forth string

	// Fifth still a thing
	Fifth int
}

func TestParse(t *testing.T) {
	config, err := Parse(Something{})

	expected := ConfigType{
		Name:        "special",
		Description: "an example config struct used for unmarshaling and validating a\nconfiguration from a file.\n\nNew Paragraph",
		Fields: []ConfigField{
			{
				Name:        "first",
				IsRequired:  true,
				Type:        "string",
				Format:      "free text",
				Example:     "\"foo\"",
				Description: "the first field in the struct",
			},
			{
				Name:        "foo-field",
				Type:        "int",
				Default:     "66",
				Description: "the number of things in a something",
			},
			{
				Name:        "third",
				Type:        "array",
				Example:     "\n  {\n    third: ['one', 'two'],\n  }",
				Description: "a list of items",
			},
			{
				Name: "forth",
				Type: "string",
			},
			{
				Name:        "fifth",
				Type:        "int",
				Description: "still a thing",
			},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, config)
}
