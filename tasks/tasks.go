package tasks

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/utils/stack"
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
	tasks   []Task
	volumes map[string]*VolumeTask
	images  map[string]*ImageTask
}

func (c *TaskCollection) add(task Task) {
	c.tasks = append(c.tasks, task)
	switch typedTask := task.(type) {
	case *VolumeTask:
		c.volumes[task.Name()] = typedTask
	case *ImageTask:
		c.images[task.Name()] = typedTask
	}
}

func (c *TaskCollection) contains(name string) bool {
	for _, task := range c.tasks {
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
	modified    map[string]bool
	tasks       *TaskCollection
	client      *docker.Client
	environment *ExecEnv
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
func NewExecuteContext(
	tasks *TaskCollection,
	client *docker.Client,
	execEnv *ExecEnv,
) *ExecuteContext {
	return &ExecuteContext{
		modified:    make(map[string]bool),
		tasks:       tasks,
		client:      client,
		environment: execEnv,
	}
}

func prepareTasks(options RunOptions) (*TaskCollection, error) {
	return prepare(options, newTaskCollection(), stack.NewStringStack())
}

func prepare(
	options RunOptions,
	tasks *TaskCollection,
	taskStack *stack.StringStack,
) (*TaskCollection, error) {
	for _, name := range options.Tasks {
		if tasks.contains(name) {
			continue
		}

		if taskStack.Contains(name) {
			return nil, fmt.Errorf(
				"Invalid dependency cycle: %s",
				strings.Join(taskStack.Items(), ", "))
		}

		resource, ok := options.Config.Resources[name]
		if !ok {
			return nil, fmt.Errorf("Resource %q does not exist", name)
		}

		task := buildTaskFromResource(taskOptions{
			name:     name,
			resource: resource,
			config:   options.Config,
		})

		taskStack.Push(name)
		options.Tasks = resource.Dependencies()
		if _, err := prepare(options, tasks, taskStack); err != nil {
			return nil, err
		}
		tasks.add(task)
		taskStack.Pop()
	}
	return tasks, nil
}

type taskOptions struct {
	name     string
	client   *docker.Client
	resource config.Resource
	config   *config.Config
}

// TODO: some way to make this a registry
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
	for _, task := range ctx.tasks.tasks {
		if err := task.Run(ctx); err != nil {
			return fmt.Errorf("Failed to execute task '%s': %s", task.Name(), err)
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

func getTaskNames(options RunOptions) []string {
	if len(options.Tasks) > 0 {
		return options.Tasks
	}

	if options.Config.Meta.Default != "" {
		return []string{options.Config.Meta.Default}
	}

	return options.Tasks
}

// Run one or more tasks
func Run(options RunOptions) error {
	options.Tasks = getTaskNames(options)
	if len(options.Tasks) == 0 {
		return fmt.Errorf("No task to run, and no default task defined.")
	}

	tasks, err := prepareTasks(options)
	if err != nil {
		return err
	}

	execEnv, err := NewExecEnv(options.Config)
	if err != nil {
		return err
	}

	ctx := NewExecuteContext(tasks, options.Client, execEnv)
	return executeTasks(ctx)
}
