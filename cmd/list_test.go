package cmd

import (
	"testing"

	"github.com/dnephin/dobi/config"
	testconfig "github.com/dnephin/dobi/internal/test/config"
	"github.com/stretchr/testify/assert"
)

func TestInclude(t *testing.T) {
	var testcases = []struct {
		doc      string
		opts     listOptions
		resource config.Resource
		expected bool
	}{
		{
			resource: &testconfig.FakeResource{},
			expected: false,
		},
		{
			opts:     listOptions{all: true},
			resource: &testconfig.FakeResource{},
			expected: true,
		},
		{
			opts:     listOptions{tags: []string{"one"}},
			resource: &testconfig.FakeResource{},
			expected: false,
		},
		{
			opts: listOptions{tags: []string{"one"}},
			resource: &testconfig.FakeResource{
				Annotations: config.Annotations{
					Annotations: config.AnnotationFields{Description: "foo"},
				},
			},
			expected: false,
		},
		{
			opts: listOptions{tags: []string{"one"}},
			resource: &testconfig.FakeResource{
				Annotations: config.Annotations{
					Annotations: config.AnnotationFields{
						Tags: []string{"one", "two"},
					},
				},
			},
			expected: true,
		},
	}

	for _, testcase := range testcases {
		actual := include(testcase.resource, testcase.opts)
		assert.Equal(t, testcase.expected, actual)
	}
}
