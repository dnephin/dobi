package task

import (
	"fmt"
	"strings"
)

// Name is an identifier for a Task
type Name struct {
	resource      string
	action        string
	defaultAction bool
}

// Name returns the full name of the task in the form 'resource:action'
func (t Name) Name() string {
	return t.resource + ":" + t.action
}

func (t Name) String() string {
	return t.Name()
}

// Resource returns the resource name of the task
func (t Name) Resource() string {
	return t.resource
}

// Action returns the action name of the task
func (t Name) Action() string {
	return t.action
}

// Equal compares two objects and returns true if they are the same
func (t Name) Equal(o Name) bool {
	return t.resource == o.resource && (t.action == o.action ||
		(t.defaultAction && o.defaultAction))
}

// MapKey returns a unique key for storing a task.Name in a map.
// Using task.Name as a key will fail when comparing default actions. See
// the comparison logic in Name.Equal().
func (t Name) MapKey() string {
	if t.defaultAction {
		return t.resource + ":DEFAULT"
	}
	return t.Name()
}

// Format the name with the name of the task, used for logging
func (t Name) Format(task string) string {
	return fmt.Sprintf("[%s:%s %s]", task, t.Action(), t.Resource())
}

// NewName returns a new task name from parts
func NewName(res, action string) Name {
	return Name{
		resource:      res,
		action:        action,
		defaultAction: action == "",
	}
}

// NewDefaultName returns a new task name, for a default action
func NewDefaultName(res, action string) Name {
	name := NewName(res, action)
	name.defaultAction = true
	return name
}

// ParseName returns a new Name from a task name string
func ParseName(name string) Name {
	name, action := splitTaskActionName(name)
	return NewName(name, action)
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
