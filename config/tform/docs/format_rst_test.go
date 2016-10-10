package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFormatRst(t *testing.T) {
	output, err := Generate(Something{}, ReStructuredText)
	assert.Nil(t, err)
	assert.Contains(t, output, "special")
}
