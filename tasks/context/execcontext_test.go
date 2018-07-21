package context

import (
	"testing"

	"github.com/dnephin/dobi/tasks/task"
	docker "github.com/fsouza/go-dockerclient"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestExecuteContext_GetAuthConfig_NoAuthConfig(t *testing.T) {
	context := ExecuteContext{}
	auth := context.GetAuthConfig("https://bogus")
	assert.Check(t, is.DeepEqual(auth, docker.AuthConfiguration{}))
}

func TestExecuteContext_IsModified(t *testing.T) {
	context := &ExecuteContext{modified: make(map[string]bool)}
	context.SetModified(task.ParseName("task1"))
	context.SetModified(task.NewDefaultName("task2", "pull"))
	context.SetModified(task.ParseName("task3:rm"))

	var testcases = []struct {
		doc      string
		name     task.Name
		expected bool
	}{
		{
			doc:      "match no action",
			name:     task.ParseName("task1"),
			expected: true,
		},
		{
			doc:      "match default action",
			name:     task.NewDefaultName("task1", "build"),
			expected: true,
		},
		{
			doc:      "match with action",
			name:     task.NewName("task3", "rm"),
			expected: true,
		},
		{
			doc:      "match with specified default action",
			name:     task.ParseName("task2:pull"),
			expected: true,
		},
		{
			doc:  "no match default action",
			name: task.NewDefaultName("task3", "build"),
		},
		{
			doc:  "no match",
			name: task.NewName("task3", "other"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.doc, func(t *testing.T) {
			assert.Equal(t, context.IsModified(tc.name), tc.expected)
		})
	}
}
