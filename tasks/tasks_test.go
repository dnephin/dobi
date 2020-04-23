package tasks

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
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
	tasks, err := collectTasks(runOptions)
	assert.Check(t, is.Nil(tasks))
	assert.Check(t, is.ErrorContains(err,
		"Invalid dependency cycle: one:run, two:run, three:run"))
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
	tasks, err := collectTasks(runOptions)
	assert.Check(t, is.Nil(err))
	assert.Check(t, is.Len(tasks.All(), 3))
}
