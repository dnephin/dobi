package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCapture(t *testing.T) {
	variable, err := parseCapture("capture(FOO)")
	assert.Nil(t, err)
	assert.Equal(t, "FOO", variable)
}

func TestParseCaptureInvalid(t *testing.T) {
	_, err := parseCapture("capture")
	assert.Error(t, err)
}
