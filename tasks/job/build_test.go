package job

import (
	"fmt"
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestBuildDockerfileWithCopy(t *testing.T) {
	mounts := []config.MountConfig{
		{
			Bind: ".",
			Path: "/opt/var/foo",
		},
		{
			Bind: "./dist",
			Path: "/go/bin",
		},
	}
	buf := buildDockerfileWithCopy("alpine:3.6", mounts)
	expected := `FROM alpine:3.6
COPY ./dist /go/bin
COPY . /opt/var/foo
`
	assert.Check(t, is.Equal(expected, buf.String()))
}

func TestGetArtifactPath(t *testing.T) {
	workingDir := "/work"
	mounts := []config.MountConfig{
		{
			Bind: ".",
			Path: "/go/src/github.com/dnephin/dobi",
		},
		{
			Bind: "./dist/bin/",
			Path: "/go/bin/",
		},
		{
			Bind: ".env",
			Path: "/code/.env",
			File: true,
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
			expected: newArtifactPath("/work/dist/bin/", "/go/bin/", "/work/dist/bin/"),
		},
		{
			doc:      "file glob, exact match with mount",
			glob:     ".env",
			expected: newArtifactPath("/work/.env", "/code/.env", "/work/.env"),
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.doc, func(t *testing.T) {
			actual, err := getArtifactPath(workingDir, testcase.glob, mounts)
			assert.NilError(t, err)
			assert.Assert(t, is.DeepEqual(testcase.expected, actual, cmpArtifactPathOpt))
		})
	}
}

var cmpArtifactPathOpt = cmp.AllowUnexported(artifactPath{})

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
		assert.Check(t, is.Equal(testcase.expected, actual), testcase.doc)
	}
}

func TestArtifactPathContainerDir(t *testing.T) {
	path := newArtifactPath("/work/dist/bin/", "/go/bin/", "/work/dist/bin/binary*")
	assert.Check(t, is.Equal("/go/bin/", path.containerDir()))

	path = newArtifactPath("/work/dist/bin/", "/go/bin/", "/work/dist/bin/")
	assert.Check(t, is.Equal("/go/bin/", path.containerDir()))
}

func TestArtifactPathContainerGlob(t *testing.T) {
	path := newArtifactPath("/work/dist/bin/", "/go/bin/", "/work/dist/bin/binary*")
	assert.Check(t, is.Equal("/go/bin/binary*", path.containerGlob()))

	path = newArtifactPath("/work/dist/bin/", "/go/bin/", "/work/dist/bin/")
	assert.Check(t, is.Equal("/go/bin/", path.containerGlob()))
}

func TestArtifactPathHostPath(t *testing.T) {
	path := newArtifactPath("/work/dist/bin/", "/go/bin/", "/work/dist/bin/")
	containerPath := "/go/bin/dobi-darwin"
	assert.Check(t, is.Equal("/work/dist/bin/dobi-darwin", path.hostPath(containerPath)))
}

func TestArtifactPathFromArchive(t *testing.T) {
	var testcases = []struct {
		artifactPath artifactPath
		archivePath  string
		expected     string
	}{
		{
			artifactPath: newArtifactPath(
				"/work/dist/bin/",
				"/go/bin/",
				"/work/dist/bin/"),
			archivePath: "bin/dobi-darwin",
			expected:    "/go/bin/dobi-darwin",
		},
		{
			artifactPath: newArtifactPath(
				"/work/",
				"/go/src/github.com/dnephin/dobi/",
				"/work/docs/build/html/"),
			archivePath: "html/",
			expected:    "/go/src/github.com/dnephin/dobi/docs/build/html/",
		},
		{
			artifactPath: newArtifactPath(
				"/work/",
				"/go/src/github.com/dnephin/dobi/",
				"/work/docs/build/html/"),
			archivePath: "html/foo/file",
			expected:    "/go/src/github.com/dnephin/dobi/docs/build/html/foo/file",
		},
	}

	for _, testcase := range testcases {
		actual := testcase.artifactPath.pathFromArchive(testcase.archivePath)
		assert.Check(t, is.Equal(testcase.expected, actual))
	}
}

func TestFileMatchesGlob(t *testing.T) {
	var testcases = []struct {
		glob     string
		path     string
		expected bool
	}{
		{
			glob:     "/go/bin/",
			path:     "/go/bin/foo/bar",
			expected: true,
		},
		{
			glob: "/go/bin",
			path: "/go/bin/foo/bar",
		},
		{
			glob:     "/work/foo-*",
			path:     "/work/foo-one",
			expected: true,
		},
		{
			glob: "/work/foo-*",
			path: "/work/foo-one/two",
		},
	}

	for _, testcase := range testcases {
		doc := fmt.Sprintf("path: %s glob: %s", testcase.path, testcase.glob)
		match, err := fileMatchesGlob(testcase.path, testcase.glob)
		if assert.Check(t, err, doc) {
			assert.Check(t, is.Equal(testcase.expected, match), doc)
		}
	}
}
