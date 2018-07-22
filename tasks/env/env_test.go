package env

import (
	"testing"

	"os"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/task"
	"gotest.tools/assert"
	"gotest.tools/env"
)

func TestTask_Run(t *testing.T) {
	var testcases = []struct {
		doc      string
		expected bool
		vars     map[string]string
	}{
		{
			doc: "set from variables",
			vars: map[string]string{
				"VAR_ONE": "override",
				"VAR_TWO": "new",
			},
			expected: true,
		},
		{
			doc: "no modifications",
			vars: map[string]string{
				"VAR_ONE": "preset",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.doc, func(t *testing.T) {
			defer env.PatchAll(t, map[string]string{
				"VAR_ONE": "preset",
			})()

			envTask := newTask(task.NewName("foo", ""), &config.EnvConfig{
				Variables: toSlice(tc.vars),
			})

			modified, err := envTask.Run(nil, false)
			assert.NilError(t, err)
			assert.Equal(t, modified, tc.expected)

			for k, v := range tc.vars {
				assert.Equal(t, os.Getenv(k), v)
			}
		})
	}
}

func toSlice(m map[string]string) []string {
	p := []string{}
	for k, v := range m {
		p = append(p, k+"="+v)
	}
	return p
}
