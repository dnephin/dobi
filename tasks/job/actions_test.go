package job

import (
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
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
