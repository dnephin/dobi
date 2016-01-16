package tasks

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/fsouza/go-dockerclient"
)

// Task is an interface implemented by all tasks
type Task interface {
	Run(ctx *ExecuteContext) error
	Name() string
}

// Task
type baseTask struct {
	name string
}

func (t *baseTask) Name() string {
	return t.name
}

// TaskCollection is a collection of Task objects
type TaskCollection struct {
	allTasks []Task
	volumes  map[string]*VolumeTask
	images   map[string]*ImageTask
}

func (c *TaskCollection) add(task Task) {
	c.allTasks = append(c.allTasks, task)
	switch typedTask := task.(type) {
	case *VolumeTask:
		c.volumes[task.Name()] = typedTask
	case *ImageTask:
		c.images[task.Name()] = typedTask
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

type eachVolumeFunc func(name string, vol *VolumeTask)

// EachVolume iterates all the volumes in names and calls f for each
func (c *TaskCollection) EachVolume(names []string, f eachVolumeFunc) {
	for _, name := range names {
		volume, _ := c.volumes[name]
		f(name, volume)
	}
}

func newTaskCollection() *TaskCollection {
	return &TaskCollection{
		volumes: make(map[string]*VolumeTask),
		images:  make(map[string]*ImageTask),
	}
}

// ExecuteContext contains all the context for task execution
type ExecuteContext struct {
	modified map[string]bool
	tasks    *TaskCollection
	client   *docker.Client
}

func (ctx *ExecuteContext) isModified(names ...string) bool {
	for _, name := range names {
		if modified, _ := ctx.modified[name]; modified {
			return true
		}
	}
	return false
}

func (ctx *ExecuteContext) setModified(name string) {
	ctx.modified[name] = true
}

// NewExecuteContext craetes a new empty ExecuteContext
func NewExecuteContext(tasks *TaskCollection, client *docker.Client) *ExecuteContext {
	return &ExecuteContext{
		modified: make(map[string]bool),
		tasks:    tasks,
		client:   client,
	}
}

func prepareTasks(options RunOptions) (*TaskCollection, error) {
	tasks := newTaskCollection()

	var prepare func(resourceNames []string) error
	// TODO: detect cyclic dependencies
	prepare = func(resourceNames []string) error {
		for _, name := range resourceNames {
			if tasks.contains(name) {
				continue
			}

			resource, ok := options.Config.Resources[name]
			if !ok {
				panic(fmt.Sprintf("Resource not defined: %s", name))
			}

			task := buildTaskFromResource(taskOptions{
				name:     name,
				resource: resource,
				config:   options.Config,
			})

			prepare(resource.Dependencies())
			tasks.add(task)
		}
		return nil
	}

	if err := prepare(options.Tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

type taskOptions struct {
	name     string
	client   *docker.Client
	resource config.Resource
	config   *config.Config
}

func buildTaskFromResource(options taskOptions) Task {
	switch conf := options.resource.(type) {
	case *config.ImageConfig:
		return NewImageTask(options, conf)
	case *config.CommandConfig:
		return NewCommandTask(options, conf)
	case *config.VolumeConfig:
		return NewVolumeTask(options, conf)
	case *config.AliasConfig:
		return NewAliasTask(options, conf)
	default:
		panic(fmt.Sprintf("Unexpected config type %T", conf))
	}
}

func executeTasks(ctx *ExecuteContext) error {
	log.Debug("executing tasks")
	for _, task := range ctx.tasks.allTasks {
		if err := task.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunOptions are the options supported by Run
type RunOptions struct {
	Client *docker.Client
	Config *config.Config
	Tasks  []string
}

// Run one or more tasks
func Run(options RunOptions) error {
	// TODO: handle empty options.Tasks, run default, or err on no default

	tasks, err := prepareTasks(options)
	if err != nil {
		return err
	}

	ctx := NewExecuteContext(tasks, options.Client)
	return executeTasks(ctx)
}
