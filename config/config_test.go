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
}
