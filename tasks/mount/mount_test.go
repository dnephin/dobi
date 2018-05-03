package mount

import (
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"gotest.tools/assert"
	"gotest.tools/fs"
)

func defaultExecContext(path string) *context.ExecuteContext {
	return context.NewExecuteContext(
		&config.Config{WorkingDir: path},
		nil,
		nil,
		context.Settings{})
}

func TestTaskRun(t *testing.T) {
	dir := fs.NewDir(t, "test-mount-task")
	defer dir.Remove()

	ctx := defaultExecContext(dir.Path())
	task := &Task{
		name: task.NewName("resource", "action"),
		config: &config.MountConfig{
			Bind: "a/b/c",
			Path: "/target",
		},
		run: runCreate,
	}

	modified, err := task.Run(ctx, false)
	assert.NilError(t, err)
	assert.Assert(t, modified)

	action := &createAction{task: task}
	assert.Assert(t, action.exists(ctx))

	// Next run is a no-op
	modified, err = task.Run(ctx, false)
	assert.NilError(t, err)
	assert.Assert(t, !modified)
}

func TestAsBind(t *testing.T) {
	workDir := "/working"
	mountConf := &config.MountConfig{
		Path: "/target",
		Bind: "./a/b/c",
	}
	expected := "/working/a/b/c:/target:rw"
	assert.Equal(t, AsBind(mountConf, workDir), expected)
}
