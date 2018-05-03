package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/renstrom/dedent"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestLoadFromBytes(t *testing.T) {
	conf := dedent.Dedent(`
		meta:
		  default: alias-def

		image=image-def:
		  image: imagename
		  dockerfile: what
		  args:
		    VERSION: "3.3.3"
		    DEBUG: 'true'

		mount=vol-def:
		  bind: dist/
		  path: /target

		job=cmd-def:
		  use: image-def
		  mounts: [vol-def]

		alias=alias-def:
		  tasks: [vol-def, cmd-def]

		compose=compose-def:
		  files: ['foo.yml']
	`)

	config, err := LoadFromBytes([]byte(conf))
	assert.NilError(t, err)

	expected := &Config{
		Meta: &MetaConfig{
			Default: "alias-def",
		},
		Resources: map[string]Resource{
			"image-def": &ImageConfig{
				Dockerfile: "what",
				Args: map[string]string{
					"VERSION": "3.3.3",
					"DEBUG":   "true",
				},
				Image: "imagename",
			},
			"vol-def": &MountConfig{
				Bind: "dist/",
				Path: "/target",
			},
			"cmd-def": &JobConfig{
				Use:    "image-def",
				Mounts: []string{"vol-def"},
			},
			"alias-def": &AliasConfig{
				Tasks: []string{"vol-def", "cmd-def"},
			},
			"compose-def": &ComposeConfig{
				Files:     []string{"foo.yml"},
				StopGrace: 5,
				Project:   "{unique}",
			},
		},
	}
	assert.DeepEqual(t, config, expected, cmpConfigOpt)
}

var cmpConfigOpt = cmp.AllowUnexported(PathGlobs{}, pull{}, ShlexSlice{})

func TestLoadFromBytesWithReservedName(t *testing.T) {
	conf := dedent.Dedent(`
		image=image-def:
		  image: imagename
		  dockerfile: what

		mount=autoclean:
		  path: dist/
		  mount: /target
	`)

	_, err := LoadFromBytes([]byte(conf))
	assert.Check(t, is.ErrorContains(err, `"autoclean" is reserved`))
}

func TestLoadFromBytesWithInvalidName(t *testing.T) {
	conf := dedent.Dedent(`
		image=image:latest:
		  image: imagename
		  dockerfile: what
	`)

	_, err := LoadFromBytes([]byte(conf))
	assert.Check(t, is.ErrorContains(err, `invalid character ":"`))
}
