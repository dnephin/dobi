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
	Name() string
}

// Task
type baseTask struct {
	name   string
	client *docker.Client
}

func (t *baseTask) Name() string {
	return t.name
}

// TaskCollection is a collection of Task objects
type TaskCollection struct {
	allTasks []Task
	volumes  map[string]*VolumeTask
}

func (c *TaskCollection) add(task Task) {
	c.allTasks = append(c.allTasks, task)
	switch typedTask := task.(type) {
	case *VolumeTask:
		c.volumes[task.Name()] = typedTask
	}
}

func (c *TaskCollection) contains(name string) bool {
	for _, task := range c.allTasks {
		if task.Name() == name {
			return true
		}
	}
	return false
}

func newTaskCollection() *TaskCollection {
	return &TaskCollection{
		volumes: make(map[string]*VolumeTask),
	}
}

func prepareTasks(options RunOptions) (*TaskCollection, error) {
	tasks := newTaskCollection()

	for _, name := range options.Pipelines {
		if tasks.contains(name) {
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
		tasks.add(task)
	}
	return tasks, nil
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

func executeTasks(tasks *TaskCollection) error {
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
