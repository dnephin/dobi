package config

import (
	"github.com/renstrom/dedent"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadFromBytes(t *testing.T) {
	conf := dedent.Dedent(`
		image-def:
		  image: imagename
		  dockerfile: what

		vol-def:
		  path: dist/
		  mount: /target

		cmd-def:
		  use: image-dev
		  volumes: [vol-def]

	`)

	config, err := LoadFromBytes([]byte(conf))
	assert.Nil(t, err)
	assert.Equal(t, len(config.Resources), 3)
	assert.IsType(t, &ImageConfig{}, config.Resources["image-def"])
	assert.IsType(t, &VolumeConfig{}, config.Resources["vol-def"])
	assert.IsType(t, &CommandConfig{}, config.Resources["cmd-def"])

	// Test default value and override
	imageConf := config.Resources["image-def"].(*ImageConfig)
	assert.Equal(t, ".", imageConf.Context)
	assert.Equal(t, "what", imageConf.Dockerfile)
}

func TestLoadFromBytesWithReservedName(t *testing.T) {
	conf := dedent.Dedent(`
		image-def:
		  image: imagename
		  dockerfile: what

		meta:
		  path: dist/
		  mount: /target
	`)

	_, err := LoadFromBytes([]byte(conf))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Name 'meta' is reserved")
}
