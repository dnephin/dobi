package job

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
)

func newCaptureTask(variable string) types.TaskBuilder {
	return func(name task.Name, conf config.Resource) types.Task {
		buffer := bytes.NewBufferString("")
		return &captureTask{
			runTask: &Task{
				name:      name,
				config:    conf.(*config.JobConfig),
				outStream: buffer,
			},
			variable: variable,
			buffer:   buffer,
		}
	}
}

type captureTask struct {
	types.NoStop
	runTask  *Task
	variable string
	buffer   *bytes.Buffer
}

// Name returns the name of the task
func (t *captureTask) Name() task.Name {
	return t.runTask.name
}

// Repr formats the task for logging
func (t *captureTask) Repr() string {
	return fmt.Sprintf("%s capture %s", t.runTask.name.Format("job"), t.variable)
}

// Run the job to capture the output in a variable
func (t *captureTask) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	modified, err := t.runTask.Run(ctx, depsModified)
	if err != nil {
		return modified, err
	}

	out := strings.TrimSpace(t.buffer.String())
	if current, ok := os.LookupEnv(t.variable); ok && current == out {
		return false, nil
	}

	logging.ForTask(t).Debug("Setting %q to: %s", t.variable, out)
	return true, os.Setenv(t.variable, out)
}
