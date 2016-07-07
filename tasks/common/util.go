package common

import "strings"

// SplitTaskActionName splits a task name into the resource, action pair
func SplitTaskActionName(name string) (string, string) {
	parts := strings.SplitN(name, ":", 2)
	switch len(parts) {
	case 2:
		return parts[0], parts[1]
	default:
		return name, ""
	}
}
