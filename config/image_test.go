package config

import (
	"reflect"
	"testing"
	"time"

	pth "github.com/dnephin/configtf/path"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ImageConfigSuite struct {
	suite.Suite
	image *ImageConfig
}

func TestImageConfigSuite(t *testing.T) {
	suite.Run(t, new(ImageConfigSuite))
}

func (s *ImageConfigSuite) SetupTest() {
	s.image = NewImageConfig()
	s.image.Dockerfile = "Dockerfile"
	s.image.Context = "."
	s.image.Image = "example"
}

func (s *ImageConfigSuite) TestString() {
	s.image.Context = "./files"
	s.Equal("Build image 'example' from 'files/Dockerfile'", s.image.String())
}

func (s *ImageConfigSuite) TestValidateMissingDependencies() {
	s.image.Depends = []string{"one", "two"}
	conf := NewConfig()
	err := validateResourcesExist(pth.NewPath(""), conf, s.image.Dependencies())
	s.Error(err)
	s.Contains(err.Error(), "missing dependencies: one, two")
}

func (s *ImageConfigSuite) TestValidateTagsWithValidFirstTag() {
	s.image.Tags = []string{"good"}
	err := s.image.ValidateTags()
	s.NoError(err)
}

func (s *ImageConfigSuite) TestValidateTagsWithBadFirstTag() {
	s.image.Tags = []string{"bad:tag"}
	err := s.image.ValidateTags()
	if s.Error(err) {
		s.Contains(err.Error(), "the first tag \"tag\" may not include an image name")
	}
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
		err := testcase.image.Validate(pth.NewPath("."), NewConfig())
		if testcase.expectedErr != "" {
			if assert.NotNil(t, err, testcase.doc) {
				assert.Contains(t, err.Error(), testcase.expectedErr, testcase.doc)
			}
		} else {
			assert.Nil(t, err, testcase.doc)
			assert.Equal(t,
				testcase.expectedDockerfile, testcase.image.Dockerfile, testcase.doc)
		}
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
	}
	resolved, err := image.Resolve(resolver)
	require.NoError(t, err)
	expected := &ImageConfig{
		Tags:  []string{"foo", "thetag"},
		Image: "last",
		Steps: "theother",
		Args: map[string]string{
			"key1": "thetag",
			"key2": "ok",
		},
	}
	assert.Equal(t, expected, resolved)
}

func TestPullWithDuration(t *testing.T) {
	p := pull{}
	now := time.Now()
	old := now.Add(-time.Duration(32 * 60 * 10e9))
	err := p.TransformConfig(reflect.ValueOf("30m"))
	require.NoError(t, err)

	assert.Equal(t, p.Required(&now), false)
	assert.Equal(t, p.Required(&old), true)
	assert.Equal(t, p.Required(nil), true)
}

func TestPullTransformConfig(t *testing.T) {
	p := pull{}
	zero := reflect.Value{}
	err := p.TransformConfig(zero)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must be a string")
}
