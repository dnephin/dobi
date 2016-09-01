package tasks

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/stretchr/testify/assert"
)

func TestCollectTasksErrorsOnCyclicDependencies(t *testing.T) {
	runOptions := RunOptions{
		Config: &config.Config{
			Resources: map[string]config.Resource{
				"one":   &config.ImageConfig{Depends: []string{"two"}},
				"two":   &config.ImageConfig{Depends: []string{"three"}},
				"three": &config.ImageConfig{Depends: []string{"four", "one"}},
				"four":  &config.ImageConfig{Depends: []string{"five"}},
				"five":  &config.ImageConfig{},
			},
		},
		Tasks: []string{"one"},
	}
	tasks, err := collectTasks(runOptions, nil)
	assert.Nil(t, tasks)
	assert.Error(t, err)
	assert.Contains(t,
		err.Error(), "Invalid dependency cycle: one:build, two:build, three:build")
}

func TestCollectTasksDoesNotErrorOnDuplicateTask(t *testing.T) {
	runOptions := RunOptions{
		Config: &config.Config{
			Resources: map[string]config.Resource{
				"one": &config.ImageConfig{},
				"two": &config.ImageConfig{Depends: []string{"one"}},
			},
		},
		Tasks: []string{"one", "two"},
	}
	tasks, err := collectTasks(runOptions, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tasks.All()))
}
