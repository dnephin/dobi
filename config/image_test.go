package config

import (
	"reflect"
	"testing"
	"time"

	pth "github.com/dnephin/configtf/path"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func sampleImageConfig() *ImageConfig {
	return &ImageConfig{
		Dockerfile: "Dockerfile",
		Context:    ".",
		Image:      "example",
	}
}

func TestImageConfigString(t *testing.T) {
	image := sampleImageConfig()
	image.Context = "./files"
	assert.Equal(t, "Build image 'example' from 'files/Dockerfile'", image.String())
}

func TestImageConfigValidateMissingDependencies(t *testing.T) {
	image := sampleImageConfig()
	image.Depends = []string{"one", "two"}
	conf := NewConfig()
	err := validateResourcesExist(pth.NewPath(""), conf, image.Dependencies())
	assert.Assert(t, is.ErrorContains(err, "missing dependencies: one, two"))
}

func TestImageConfigValidateTagsWithValidFirstTag(t *testing.T) {
	image := sampleImageConfig()
	image.Tags = []string{"good"}
	err := image.ValidateTags()
	assert.NilError(t, err)
}

func TestImageConfigValidateTagsWithBadFirstTag(t *testing.T) {
	image := sampleImageConfig()
	image.Tags = []string{"bad:tag"}
	err := image.ValidateTags()
	expected := "the first tag \"tag\" may not include an image name"
	assert.Assert(t, is.ErrorContains(err, expected))
}

func TestImageConfigValidate(t *testing.T) {
	var testcases = []struct {
		doc                string
		image              *ImageConfig
		expectedErr        string
		expectedDockerfile string
	}{
		{
			doc: "dockerfile and steps both set",
			image: &ImageConfig{
				Steps:      "FROM alpine:3.6",
				Dockerfile: "Dockerfile",
				Context:    ".",
			},
			expectedErr: "dockerfile can not be used with steps",
		},
		{
			doc:         "missing required field",
			image:       &ImageConfig{Steps: "FROM alpine:3.6"},
			expectedErr: "one of context, or pull is required",
		},
		{
			doc:                "just context",
			image:              &ImageConfig{Context: "."},
			expectedDockerfile: "Dockerfile",
		},
		{
			doc:   "just pull",
			image: &ImageConfig{Pull: pull{action: pullAlways}},
		},
		{
			doc:                "just dockerfile",
			image:              &ImageConfig{Dockerfile: "Dockerfile"},
			expectedDockerfile: "Dockerfile",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.doc, func(t *testing.T) {
			err := testcase.image.Validate(pth.NewPath("."), NewConfig())
			if testcase.expectedErr != "" {
				assert.Assert(t, is.ErrorContains(err, testcase.expectedErr))
				return
			}

			assert.Assert(t, err == nil)
			assert.Assert(t, is.Equal(testcase.expectedDockerfile, testcase.image.Dockerfile))
		})
	}
}

func TestImageConfigResolve(t *testing.T) {
	resolver := newFakeResolver(map[string]string{
		"{one}":   "thetag",
		"{two}":   "theother",
		"{three}": "last",
	})

	image := &ImageConfig{
		Tags:  []string{"foo", "{one}"},
		Image: "{three}",
		Steps: "{two}",
		Args: map[string]string{
			"key1": "{one}",
			"key2": "ok",
		},
		CacheFrom: []string{"{one}", "two"},
	}
	resolved, err := image.Resolve(resolver)
	assert.NilError(t, err)
	expected := &ImageConfig{
		Tags:  []string{"foo", "thetag"},
		Image: "last",
		Steps: "theother",
		Args: map[string]string{
			"key1": "thetag",
			"key2": "ok",
		},
		CacheFrom: []string{"thetag", "two"},
	}
	assert.Check(t, is.DeepEqual(expected, resolved, cmpConfigOpt))
}

func TestPullWithDuration(t *testing.T) {
	p := pull{}
	now := time.Now()
	old := now.Add(-time.Duration(32 * 60 * 10e9))
	err := p.TransformConfig(reflect.ValueOf("30m"))
	assert.NilError(t, err)

	assert.Check(t, !p.Required(&now))
	assert.Check(t, p.Required(&old))
	assert.Check(t, p.Required(nil))
}

func TestPullTransformConfig(t *testing.T) {
	p := pull{}
	zero := reflect.Value{}
	err := p.TransformConfig(zero)

	assert.Check(t, is.ErrorContains(err, "must be a string"))
}
