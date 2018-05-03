package config

import (
	"reflect"
	"testing"

	pth "github.com/dnephin/configtf/path"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestJobConfigString(t *testing.T) {
	job := &JobConfig{
		Use:      "builder",
		Command:  ShlexSlice{original: "run"},
		Artifact: PathGlobs{globs: []string{"foo"}},
	}
	assert.Equal(t, job.String(), "Run 'run' using the 'builder' image to create 'foo'")
}

func TestJobConfigValidateMissingUse(t *testing.T) {
	conf := NewConfig()
	conf.Resources["example"] = &AliasConfig{}
	job := &JobConfig{Use: "example"}
	err := job.Validate(pth.NewPath(""), conf)
	assert.Assert(t, is.ErrorContains(err, "example is not an image resource"))
}

func TestJobConfigValidateMissingMount(t *testing.T) {
	conf := NewConfig()
	conf.Resources["one"] = NewImageConfig()
	conf.Resources["two"] = NewImageConfig()
	conf.Resources["example"] = NewImageConfig()
	job := &JobConfig{}
	job.Use = "example"
	job.Mounts = []string{"one", "two"}

	err := job.Validate(pth.NewPath(""), conf)
	assert.Assert(t, is.ErrorContains(err, "one is not a mount resource"))
}

func TestJobConfigRunFromConfig(t *testing.T) {
	values := map[string]interface{}{
		"use":        "image-res",
		"command":    "echo foo",
		"entrypoint": "bash -c",
	}
	res, err := jobFromConfig("foo", values)
	job, ok := res.(*JobConfig)
	assert.Assert(t, ok)
	assert.NilError(t, err)
	// TODO: compare against the entire struct
	assert.Equal(t, job.Use, "image-res")
	assert.Assert(t, is.DeepEqual(job.Command.Value(), []string{"echo", "foo"}))
	assert.Assert(t, is.DeepEqual(job.Entrypoint.Value(), []string{"bash", "-c"}))
}

func TestShlexSliceTransformConfig(t *testing.T) {
	s := ShlexSlice{}
	zero := reflect.Value{}
	err := s.TransformConfig(zero)

	assert.Check(t, is.ErrorContains(err, "must be a string"))
}
