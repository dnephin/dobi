package iface

import (
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
)

// Task interface performs some operation with a resource config
type Task interface {
	logging.LogRepresenter
	Name() common.TaskName
	Run(*context.ExecuteContext) error
	Stop(*context.ExecuteContext) error
}

// TaskConfig is a data object which stores the full configuration of a Task
type TaskConfig interface {
	Name() common.TaskName
	Resource() config.Resource
	Dependencies() []string
	Task(config.Resource) Task
}
