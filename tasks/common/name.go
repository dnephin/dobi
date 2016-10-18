package common

import "strings"

// TaskName is an identifier for a Task
type TaskName struct {
	resource string
	action   string
}

// Name returns the full name of the task in the form 'resource:action'
func (t TaskName) Name() string {
	return t.resource + ":" + t.action
}

func (t TaskName) String() string {
	return t.Name()
}

// Resource returns the resource name of the task
func (t TaskName) Resource() string {
	return t.resource
}

// Action returns the action name of the task
func (t TaskName) Action() string {
	return t.action
}

// Equal compares two objects and returns true if they are the same
func (t TaskName) Equal(o TaskName) bool {
	return t.resource == o.resource && t.action == o.action
}

// NewTaskName returns a new task name from parts
func NewTaskName(res, action string) TaskName {
	return TaskName{resource: res, action: action}
}

// ParseTaskName returns a new TaskName from a task name string
func ParseTaskName(name string) TaskName {
	name, action := splitTaskActionName(name)
	return TaskName{resource: name, action: action}
}

// splitTaskActionName splits a task name into the resource, action pair
func splitTaskActionName(name string) (string, string) {
	parts := strings.SplitN(name, ":", 2)
	switch len(parts) {
	case 2:
		return parts[0], parts[1]
	default:
		return name, ""
	}
}
