package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	"github.com/docker/docker/runconfig/opts"
)

// GetTaskConfig returns a new task for the action
func GetTaskConfig(name, action string, conf *config.EnvConfig) (types.TaskConfig, error) {
	switch action {
	case "", "set":
		return types.NewTaskConfig(
			task.NewDefaultName(name, "set"), conf, task.NoDependencies, newTask), nil
	case "rm":
		return types.NewTaskConfig(
			task.NewName(name, "rm"), conf, task.NoDependencies, newRemoveTask), nil
	default:
		return nil, fmt.Errorf("Invalid env action %q for task %q", action, name)
	}
}

// Task sets environment variables
type Task struct {
	types.NoStop
	name   task.Name
	config *config.EnvConfig
}

func newTask(name task.Name, conf config.Resource) types.Task {
	return &Task{name: name, config: conf.(*config.EnvConfig)}
}

// Name returns the name of the task
func (t *Task) Name() task.Name {
	return t.name
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	return t.name.Format("env")
}

// Run sets environment variables
func (t *Task) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	for _, filename := range t.config.Files {
		vars, err := opts.ParseEnvFile(filename)
		if err != nil {
			return true, err
		}
		if err := setVariables(vars); err != nil {
			return true, err
		}
	}
	if err := setVariables(t.config.Variables); err != nil {
		return true, err
	}
	logging.ForTask(t).Info("Done")
	return true, nil
}

func setVariables(vars []string) error {
	for _, variable := range vars {
		key, value, err := splitVar(variable)
		if err != nil {
			return err
		}
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}

func splitVar(variable string) (string, string, error) {
	parts := strings.SplitN(variable, "=", 2)
	if len(parts) < 2 {
		return variable, "", fmt.Errorf("invalid variable format %q", variable)
	}
	return parts[0], parts[1], nil
}

func newRemoveTask(name task.Name, conf config.Resource) types.Task {
	return &removeTask{name: name}
}

type removeTask struct {
	types.NoStop
	name task.Name
}

// Name returns the name of the task
func (t *removeTask) Name() task.Name {
	return t.name
}

// Repr formats the task for logging
func (t *removeTask) Repr() string {
	return t.name.Format("env")
}

// Run does nothing
func (t *removeTask) Run(ctx *context.ExecuteContext, _ bool) (bool, error) {
	return false, nil
}
