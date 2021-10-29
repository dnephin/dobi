package task

import (
	"fmt"
	"regexp"
	"strings"
)

// Name is an identifier for a Task
type Name struct {
	resource   string
	action     Action
	captureVar string
}

type Action string

const (
	Remove  Action = "rm"
	Create  Action = "create"
	Pull    Action = "pull"
	Push    Action = "push"
	Tag     Action = "tag"
	Attach  Action = "attach"
	Detach  Action = "detach"
	Capture Action = "capture"
)

func (a Action) String() string {
	return string(a)
}

// nolint: gocyclo
func NewAction(s string) (Action, error) {
	switch strings.ToLower(s) {
	case "remove", "rm", "delete", "down":
		return Remove, nil
	case "", "default", "create", "build", "run", "up":
		return Create, nil
	case "tag":
		return Tag, nil
	case "attach":
		return Attach, nil
	case "detach":
		return Detach, nil
	case "capture":
		return Capture, nil
	case "push", "upload":
		return Push, nil
	case "pull", "download":
		return Pull, nil
	default:
		return "", fmt.Errorf("%q is not a valid action", s)
	}
}

// Name returns the full name of the task in the form 'resource:action'
func (t Name) Name() string {
	if t.action == Capture {
		return t.resource + ":" + t.action.String() + "(" + t.captureVar + ")"
	}
	return t.resource + ":" + t.action.String()
}

func (t Name) String() string {
	return t.Name()
}

// Resource returns the resource name of the task
func (t Name) Resource() string {
	return t.resource
}

// Action returns the action name of the task
func (t Name) Action() Action {
	return t.action
}

// CaptureVar returns the capture variable of the task
func (t Name) CaptureVar() string {
	return t.captureVar
}

// Equal compares two objects and returns true if they are the same
func (t Name) Equal(o Name) bool {
	return t.resource == o.resource && t.action == o.action
}

// Format the name with the name of the task, used for logging
func (t Name) Format(task string) string {
	return fmt.Sprintf("[%s:%s %s]", task, t.Action(), t.Resource())
}

// NewName returns a new task name from parts
func NewName(res string, action Action) Name {
	return Name{
		resource: res,
		action:   action,
	}
}

// NewNameForCapture returns a new task name for a capture task from parts
func NewNameForCapture(res string, captureVar string) Name {
	return Name{
		resource:   res,
		action:     Capture,
		captureVar: captureVar,
	}
}

// ParseName returns a new Name from a task name string
func ParseName(name string) (Name, error) {
	name, actionString := splitTaskActionName(name)
	if strings.HasPrefix(actionString, "capture") {
		variable, err := parseCapture(actionString)
		if err != nil {
			return Name{}, err
		}
		return NewNameForCapture(name, variable), err
	}
	action, err := NewAction(actionString)
	return NewName(name, action), err
}

func ParseNames(names []string) ([]Name, error) {
	var tasks []Name
	for _, name := range names {
		parsed, err := ParseName(name)
		if err != nil {
			return []Name{}, err
		}
		tasks = append(tasks, parsed)
	}
	return tasks, nil
}

var (
	captureRegex = regexp.MustCompile(`^capture\((\w+)\)$`)
)

func parseCapture(action string) (string, error) {
	matches := captureRegex.FindStringSubmatch(action)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("invalid capture format %q", action)
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
