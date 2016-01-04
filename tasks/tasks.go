package tasks

import (
	"fmt"

	//	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/buildpipe/config"
	"github.com/fsouza/go-dockerclient"
	//	"github.com/hashicorp/errwrap"
)

// Task is an interface implemented by all tasks
type Task interface {
	Run() error
	Dependencies() []string
}

// Task
type baseTask struct {
	name   string
	client *docker.Client
}

func prepareTasks(options RunOptions) (*[]Task, error) {
	tasks := []Task{}
	prepared := make(map[string]bool)

	for _, name := range options.Pipelines {
		if _, ok := prepared[name]; ok {
			continue
		}

		resource, ok := options.Config.Resources[name]
		if !ok {
			return nil, fmt.Errorf("Resource not defined: %s", name)
		}

		task := buildTaskFromResource(taskOptions{
			name:     name,
			client:   options.Client,
			resource: &resource,
		})

		// TODO: recursively build tasks for dependencies first

		tasks = append(tasks, task)
		prepared[name] = true
	}
	return &tasks, nil
}

type taskOptions struct {
	name     string
	client   *docker.Client
	resource config.Resource
}

func buildTaskFromResource(options taskOptions) Task {
	switch conf := options.resource.(type) {
	case config.ImageConfig:

	case config.CommandConfig:

	case config.VolumeConfig:
		return NewVolumeTask(options, conf)
	default:
		panic(fmt.Sprintf("Unexpected config type %T", conf))
	}
	return nil
}

func executeTasks(tasks *[]Task) error {
	return nil
}

// RunOptions are the options supported by Run
type RunOptions struct {
	Client    *docker.Client
	Config    *config.Config
	Pipelines []string
}

// Run one or more pipelines
func Run(options RunOptions) error {
	tasks, err := prepareTasks(options)
	if err != nil {
		return err
	}
	return executeTasks(tasks)
}
