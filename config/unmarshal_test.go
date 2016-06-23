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

		volume=vol-def:
		  path: dist/
		  mount: /target

		command=cmd-def:
		  use: image-dev
		  volumes: [vol-def]

		alias=alias-def:
		  tasks: [vol-def, cmd-def]
	`)

	config, err := LoadFromBytes([]byte(conf))
	assert.Nil(t, err)
	assert.Equal(t, 4, len(config.Resources))
	assert.IsType(t, &ImageConfig{}, config.Resources["image-def"])
	assert.IsType(t, &VolumeConfig{}, config.Resources["vol-def"])
	assert.IsType(t, &CommandConfig{}, config.Resources["cmd-def"])
	assert.IsType(t, &AliasConfig{}, config.Resources["alias-def"])

	// Test default value and override
	imageConf := config.Resources["image-def"].(*ImageConfig)
	assert.Equal(t, ".", imageConf.Context)
	assert.Equal(t, "what", imageConf.Dockerfile)

	volumeConf := config.Resources["vol-def"].(*VolumeConfig)
	assert.Equal(t, "dist/", volumeConf.Path)
	assert.Equal(t, "/target", volumeConf.Mount)
	assert.Equal(t, "rw", volumeConf.Mode)

	aliasConf := config.Resources["alias-def"].(*AliasConfig)
	assert.Equal(t, []string{"vol-def", "cmd-def"}, aliasConf.Tasks)

	assert.Equal(t, &MetaConfig{Default: "alias-def"}, config.Meta)
}

func TestLoadFromBytesWithReservedName(t *testing.T) {
	conf := dedent.Dedent(`
		image=image-def:
		  image: imagename
		  dockerfile: what

		volume=autoclean:
		  path: dist/
		  mount: /target
	`)

	_, err := LoadFromBytes([]byte(conf))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Name \"autoclean\" is reserved")
}
