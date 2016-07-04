package config

import (
	"github.com/renstrom/dedent"
	"github.com/stretchr/testify/assert"
	"testing"
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

		run=cmd-def:
		  use: image-dev
		  mounts: [vol-def]

		alias=alias-def:
		  tasks: [vol-def, cmd-def]
	`)

	config, err := LoadFromBytes([]byte(conf))
	assert.Nil(t, err)
	assert.Equal(t, 4, len(config.Resources))
	assert.IsType(t, &ImageConfig{}, config.Resources["image-def"])
	assert.IsType(t, &MountConfig{}, config.Resources["vol-def"])
	assert.IsType(t, &RunConfig{}, config.Resources["cmd-def"])
	assert.IsType(t, &AliasConfig{}, config.Resources["alias-def"])

	// Test default value and override
	imageConf := config.Resources["image-def"].(*ImageConfig)
	assert.Equal(t, "what", imageConf.Dockerfile)
	assert.Equal(t, map[string]string{
		"VERSION": "3.3.3",
		"DEBUG":   "true",
	}, imageConf.Args)

	mountConf := config.Resources["vol-def"].(*MountConfig)
	assert.Equal(t, "dist/", mountConf.Bind)
	assert.Equal(t, "/target", mountConf.Path)
	assert.Equal(t, false, mountConf.ReadOnly)

	aliasConf := config.Resources["alias-def"].(*AliasConfig)
	assert.Equal(t, []string{"vol-def", "cmd-def"}, aliasConf.Tasks)

	assert.Equal(t, &MetaConfig{Default: "alias-def"}, config.Meta)
}

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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "\"autoclean\" is reserved")
}

func TestLoadFromBytesWithInvalidName(t *testing.T) {
	conf := dedent.Dedent(`
		image=image:latest:
		  image: imagename
		  dockerfile: what
	`)

	_, err := LoadFromBytes([]byte(conf))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid character \":\"")
}
