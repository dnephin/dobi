package image

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"gotest.tools/v3/assert"
)

func TestForEachTag(t *testing.T) {
	task := Task{
		config: &config.ImageConfig{
			Image: "imagename",
			Tags:  []string{"plain", "repo:tag"},
		},
	}

	expected := []string{"imagename:plain", "repo:tag"}
	tags := []string{}

	eachFunc := func(imageTag string) error {
		tags = append(tags, imageTag)
		return nil
	}

	err := task.ForEachTag(nil, eachFunc)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, tags)
}
