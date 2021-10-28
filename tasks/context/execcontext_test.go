package context

import (
	"testing"

	"github.com/dnephin/dobi/tasks/task"
	docker "github.com/fsouza/go-dockerclient"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestExecuteContext_GetAuthConfig_NoAuthConfig(t *testing.T) {
	context := ExecuteContext{}
	auth := context.GetAuthConfig("https://bogus")
	assert.Check(t, is.DeepEqual(auth, docker.AuthConfiguration{}))
}

func parseNameNoErrors(name string) task.Name {
	tName, _ := task.ParseName(name)
	return tName
}

func TestExecuteContext_IsModified(t *testing.T) {
	context := &ExecuteContext{modified: make(map[string]bool)}
	context.SetModified(parseNameNoErrors("task1"))
	context.SetModified(parseNameNoErrors("task2:create"))
	context.SetModified(parseNameNoErrors("task3:rm"))

	var testcases = []struct {
		doc      string
		name     task.Name
		expected bool
	}{
		{
			doc:      "match no action",
			name:     parseNameNoErrors("task1"),
			expected: true,
		},
		{
			doc:      "match default action",
			name:     parseNameNoErrors("task1:build"),
			expected: true,
		},
		{
			doc:      "match with action",
			name:     parseNameNoErrors("task3:rm"),
			expected: true,
		},
		{
			doc:      "match default with specified action",
			name:     parseNameNoErrors("task2"),
			expected: true,
		},
		{
			doc:  "no match valid wrong action",
			name: parseNameNoErrors("task3:build"),
		},
		{
			doc:  "no match wrong default action",
			name: parseNameNoErrors("task3"),
		},
		{
			doc:  "no match invalid",
			name: task.NewName("task3", "other"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.doc, func(t *testing.T) {
			assert.Equal(t, context.IsModified(tc.name), tc.expected)
		})
	}
}
