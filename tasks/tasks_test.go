package tasks

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/stretchr/testify/assert"
)

func aliasWithDeps(deps []string) config.Resource {
	return &config.AliasConfig{Tasks: deps}
}

func TestCollectTasksErrorsOnCyclicDependencies(t *testing.T) {
	runOptions := RunOptions{
		Config: &config.Config{
			Resources: map[string]config.Resource{
				"one":   aliasWithDeps([]string{"two"}),
				"two":   aliasWithDeps([]string{"three"}),
				"three": aliasWithDeps([]string{"four", "one"}),
				"four":  aliasWithDeps([]string{"five"}),
				"five":  aliasWithDeps([]string{}),
			},
		},
		Tasks: []string{"one"},
	}
	tasks, err := collectTasks(runOptions, nil)
	assert.Nil(t, tasks)
	assert.Error(t, err)
	assert.Contains(t,
		err.Error(), "Invalid dependency cycle: one:run, two:run, three:run")
}

func TestCollectTasksDoesNotErrorOnDuplicateTask(t *testing.T) {
	runOptions := RunOptions{
		Config: &config.Config{
			Resources: map[string]config.Resource{
				"one": &config.ImageConfig{},
				"two": aliasWithDeps([]string{"one"}),
			},
		},
		Tasks: []string{"one", "two"},
	}
	tasks, err := collectTasks(runOptions, nil)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(tasks.All()))
}
