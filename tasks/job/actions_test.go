package job

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestParseCapture(t *testing.T) {
	variable, err := parseCapture("capture(FOO)")
	assert.Check(t, is.Nil(err))
	assert.Check(t, is.Equal("FOO", variable))
}

func TestParseCaptureInvalid(t *testing.T) {
	_, err := parseCapture("capture")
	assert.Check(t, is.ErrorContains(err, "invalid capture format"))
}
