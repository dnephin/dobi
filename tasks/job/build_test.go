package job

import (
	"github.com/dnephin/dobi/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetArtifactPath(t *testing.T) {
	workingDir := "/work"
	mounts := []config.MountConfig{
		{
			Bind: ".",
			Path: "/go/src/github.com/dnephin/dobi",
		},
		{
			Bind: "./dist/bin/",
			Path: "/go/bin",
		},
	}

	var testcases = []struct {
		doc      string
		glob     string
		expected artifactPath
	}{
		{
			doc:      "directory glob, exact match with mount",
			glob:     "./dist/bin/",
			expected: newArtifactPath("/work/dist/bin/", "/go/bin", "/work/dist/bin/"),
		},
	}
	for _, testcase := range testcases {
		actual, err := getArtifactPath(workingDir, testcase.glob, mounts)
		if assert.NoError(t, err, testcase.doc) {
			assert.Equal(t, testcase.expected, actual, testcase.doc)
		}
	}
}

func TestHasPathPrefix(t *testing.T) {
	var testcases = []struct {
		doc      string
		path     string
		prefix   string
		expected bool
	}{
		{
			doc:      "identical parts match",
			path:     "/one/two/three",
			prefix:   "/one/two/three",
			expected: true,
		},
		{
			doc:      "parts match with trailing slash",
			path:     "/one/two/three/",
			prefix:   "/one/two/three",
			expected: true,
		},
		{
			doc:      "parts match with trailing slash on prefix",
			path:     "/one/two/three",
			prefix:   "/one/two/three/",
			expected: true,
		},
		{
			doc:      "prefix match",
			path:     "/one/two/three",
			prefix:   "/one/two",
			expected: true,
		},
		{
			doc:    "item mismatch",
			path:   "/one/two/three",
			prefix: "/one/three/three",
		},
		{
			doc:    "prefix longer mismatch",
			path:   "/one/two/three",
			prefix: "/one/two/three/four",
		},
	}
	for _, testcase := range testcases {
		actual := hasPathPrefix(testcase.path, testcase.prefix)
		assert.Equal(t, testcase.expected, actual, testcase.doc)
	}
}
