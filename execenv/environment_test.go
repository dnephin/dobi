package execenv

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
	"github.com/gotestyourself/gotestyourself/fs"
)

func TestNewExecEnvFromConfigDefault(t *testing.T) {
	tmpDir := fs.NewDir(t, "test-environment")
	defer tmpDir.Remove()
	execEnv, err := NewExecEnvFromConfig("", "", tmpDir.Path())
	assert.NilError(t, err)
	expected := fmt.Sprintf("%s-root", filepath.Base(tmpDir.Path()))
	assert.Equal(t, expected, execEnv.Unique())
}

func TestNewExecEnvFromConfigWithTemplate(t *testing.T) {
	tmpDir := fs.NewDir(t, "test-environment")
	defer tmpDir.Remove()
	os.Setenv("EXEC_ID", "Use-This")
	defer os.Unsetenv("EXEC_ID")

	execEnv, err := NewExecEnvFromConfig("{env.EXEC_ID}", "", tmpDir.Path())
	assert.NilError(t, err)
	assert.Equal(t, "Use-This", execEnv.ExecID)
}

func TestNewExecEnvFromConfigWithInvalidTemplate(t *testing.T) {
	tmpDir := fs.NewDir(t, "test-environment")
	defer tmpDir.Remove()
	_, err := NewExecEnvFromConfig("{env.bogus} ", "", tmpDir.Path())
	expected := `a value is required for variable "env.bogus"`
	assert.Assert(t, is.ErrorContains(err, expected))
}

func TestValidateExecIDEmpty(t *testing.T) {
	output, err := validateExecID("")
	assert.Equal(t, "", output)
	expected := "exec-id template was empty after rendering"
	assert.Assert(t, is.ErrorContains(err, expected))
}

func TestValidateExecIDTooManyLines(t *testing.T) {
	output, err := validateExecID("one\ntwo")
	assert.Equal(t, "", output)
	assert.Assert(t, is.ErrorContains(err, "rendered to 2 lines"))
}

func TestValidateExecIDValid(t *testing.T) {
	output, err := validateExecID("one\n")
	assert.NilError(t, err)
	assert.Equal(t, "one", output)

	output, err = validateExecID("one")
	assert.NilError(t, err)
	assert.Equal(t, "one", output)
}

func TestResolveUnique(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	tmpl := "ok-{unique}"
	expected := "ok-" + execEnv.Unique()
	value, err := execEnv.Resolve(tmpl)

	assert.NilError(t, err)
	assert.Equal(t, value, expected)
	assert.Equal(t, execEnv.tmplCache[tmpl], expected)
}

func TestResolveUnknown(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	_, err := execEnv.Resolve("{bogus}")
	assert.Assert(t, is.ErrorContains(err, `unknown variable "bogus"`))
}

func TestResolveBadTemplate(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	_, err := execEnv.Resolve("{bogus{")

	assert.Assert(t, is.ErrorContains(err, "Cannot find end tag"))
}

func TestResolveEnvironmentNoDefault(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	_, err := execEnv.Resolve("thing-{env.foo}")

	assert.Assert(t, is.ErrorContains(err, `required for variable "env.foo"`))
}

func TestResolveEnvironment(t *testing.T) {
	defer os.Unsetenv("FOO")
	os.Setenv("FOO", "stars")
	tmpl := "thing-{env.FOO}"
	expected := "thing-stars"

	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve(tmpl)

	assert.NilError(t, err)
	assert.Equal(t, value, expected)
	assert.Equal(t, execEnv.tmplCache[tmpl], expected)
}

func TestResolveTime(t *testing.T) {
	tmpl := "build-{time.YYYY-MM-DD}"
	expected := "build-2016-04-05"

	execEnv := NewExecEnv("exec", "project", "cwd")
	execEnv.startTime = time.Date(2016, 4, 5, 0, 0, 0, 0, time.UTC)
	value, err := execEnv.Resolve(tmpl)

	assert.NilError(t, err)
	assert.Equal(t, value, expected)
	assert.Equal(t, execEnv.tmplCache[tmpl], expected)
}

func TestSplitDefault(t *testing.T) {
	tag := "time.19:01:01:default"
	value, defVal, hasDefault := splitDefault(tag)
	assert.Equal(t, value, "time.19:01:01")
	assert.Equal(t, defVal, "default")
	assert.Equal(t, hasDefault, true)
}

func TestSplitDefaultNoDefault(t *testing.T) {
	tag := "env.FOO"
	value, defVal, hasDefault := splitDefault(tag)
	assert.Equal(t, value, "env.FOO")
	assert.Equal(t, defVal, "")
	assert.Equal(t, hasDefault, false)
}

func TestResolveUserName(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.name}")
	assert.NilError(t, err)
	assert.Equal(t, value, "root")
}

func TestResolveUserUID(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.uid}")
	assert.NilError(t, err)
	assert.Equal(t, value, "0")
}

func TestResolveUserGroup(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.group}")
	assert.NilError(t, err)
	assert.Equal(t, value, "root")
}

func TestResolveUserGID(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.gid}")
	assert.NilError(t, err)
	assert.Equal(t, value, "0")
}

func TestResolveUserHome(t *testing.T) {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.home}")
	assert.NilError(t, err)
	assert.Equal(t, value, "/root")
}

func TestSplitPrefixNoPrefix(t *testing.T) {
	for _, tag := range []string{".foo", "foo.", "foo"} {
		prefix, suffix := splitPrefix(tag)
		assert.Check(t, is.Equal(prefix, ""))
		assert.Check(t, is.Equal(suffix, tag))
	}
}

func TestSplitPrefix(t *testing.T) {
	prefix, suffix := splitPrefix("fo.o")
	assert.Check(t, is.Equal(prefix, "fo"))
	assert.Check(t, is.Equal(suffix, "o"))
}
