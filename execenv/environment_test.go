package execenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExecEnvSuite struct {
	suite.Suite
	tmpDir string
}

func TestExecEnvSuite(t *testing.T) {
	suite.Run(t, new(ExecEnvSuite))
}

func (s *ExecEnvSuite) SetupTest() {
	var err error
	s.tmpDir, err = ioutil.TempDir("", "environment-test")
	s.Require().Nil(err)
}

func (s *ExecEnvSuite) TearDownTest() {
	s.Nil(os.RemoveAll(s.tmpDir))
}

func (s *ExecEnvSuite) TestNewExecEnvFromConfigDefault() {
	execEnv, err := NewExecEnvFromConfig("", "", s.tmpDir)
	s.Nil(err)
	expected := fmt.Sprintf("%s-root", filepath.Base(s.tmpDir))
	s.Equal(expected, execEnv.Unique())
}

func (s *ExecEnvSuite) TestNewExecEnvFromConfigWithTemplate() {
	os.Setenv("EXEC_ID", "Use-This")
	defer os.Unsetenv("EXEC_ID")

	execEnv, err := NewExecEnvFromConfig("{env.EXEC_ID}", "", s.tmpDir)
	s.Nil(err)
	s.Equal("Use-This", execEnv.ExecID)
}

func (s *ExecEnvSuite) TestNewExecEnvFromConfigWithInvalidTemplate() {
	_, err := NewExecEnvFromConfig("{env.bogus} ", "", s.tmpDir)
	s.Error(err)
	s.Contains(err.Error(), "a value is required for variable \"env.bogus\"")
}

func (s *ExecEnvSuite) TestValidateExecIDEmpty() {
	output, err := validateExecID("")
	s.Equal("", output)
	s.Error(err)
	s.Contains(err.Error(), "exec-id template was empty after rendering")
}

func (s *ExecEnvSuite) TestValidateExecIDTooManyLines() {
	output, err := validateExecID("one\ntwo")
	s.Equal("", output)
	s.Error(err)
	s.Contains(err.Error(), "rendered to 2 lines")
}

func (s *ExecEnvSuite) TestValidateExecIDValid() {
	output, err := validateExecID("one\n")
	s.Nil(err)
	s.Equal("one", output)

	output, err = validateExecID("one")
	s.Nil(err)
	s.Equal("one", output)
}

func (s *ExecEnvSuite) TestResolveUnique() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	tmpl := "ok-{unique}"
	expected := "ok-" + execEnv.Unique()
	value, err := execEnv.Resolve(tmpl)

	s.Nil(err)
	s.Equal(value, expected)
	s.Equal(execEnv.tmplCache[tmpl], expected)
}

func (s *ExecEnvSuite) TestResolveUnknown() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	_, err := execEnv.Resolve("{bogus}")

	s.Error(err)
	s.Contains(err.Error(), "unknown variable \"bogus\"")
}

func (s *ExecEnvSuite) TestResolveBadTemplate() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	_, err := execEnv.Resolve("{bogus{")

	s.Error(err)
	s.Contains(err.Error(), "Cannot find end tag")
}

func (s *ExecEnvSuite) TestResolveEnvironmentNoDefault() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	_, err := execEnv.Resolve("thing-{env.foo}")

	s.Error(err)
	s.Contains(err.Error(), "required for variable \"env.foo\"")
}

func (s *ExecEnvSuite) TestResolveEnvironment() {
	defer os.Unsetenv("FOO")
	os.Setenv("FOO", "stars")
	tmpl := "thing-{env.FOO}"
	expected := "thing-stars"

	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve(tmpl)

	s.Nil(err)
	s.Equal(value, expected)
	s.Equal(execEnv.tmplCache[tmpl], expected)
}

func (s *ExecEnvSuite) TestResolveTime() {
	tmpl := "build-{time.YYYY-MM-DD}"
	expected := "build-2016-04-05"

	execEnv := NewExecEnv("exec", "project", "cwd")
	execEnv.startTime = time.Date(2016, 4, 5, 0, 0, 0, 0, time.UTC)
	value, err := execEnv.Resolve(tmpl)

	s.Nil(err)
	s.Equal(value, expected)
	s.Equal(execEnv.tmplCache[tmpl], expected)
}

func (s *ExecEnvSuite) TestSplitDefault() {
	tag := "time.19:01:01:default"
	value, defVal, hasDefault := splitDefault(tag)
	s.Equal(value, "time.19:01:01")
	s.Equal(defVal, "default")
	s.Equal(hasDefault, true)
}

func (s *ExecEnvSuite) TestSplitDefaultNoDefault() {
	tag := "env.FOO"
	value, defVal, hasDefault := splitDefault(tag)
	s.Equal(value, "env.FOO")
	s.Equal(defVal, "")
	s.Equal(hasDefault, false)
}

func (s *ExecEnvSuite) TestResolveUserName() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.name}")
	s.Nil(err)
	s.Equal(value, "root")
}

func (s *ExecEnvSuite) TestResolveUserUID() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.uid}")
	s.Nil(err)
	s.Equal(value, "0")
}

func (s *ExecEnvSuite) TestResolveUserGroup() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.group}")
	s.Nil(err)
	s.Equal(value, "root")
}

func (s *ExecEnvSuite) TestResolveUserGID() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.gid}")
	s.Nil(err)
	s.Equal(value, "0")
}

func (s *ExecEnvSuite) TestResolveUserHome() {
	execEnv := NewExecEnv("exec", "project", "cwd")
	value, err := execEnv.Resolve("{user.home}")
	s.Nil(err)
	s.Equal(value, "/root")
}

func TestSplitPrefixNoPrefix(t *testing.T) {
	for _, tag := range []string{".foo", "foo.", "foo"} {
		prefix, suffix := splitPrefix(tag)
		assert.Equal(t, prefix, "")
		assert.Equal(t, suffix, tag)
	}
}

func TestSplitPrefix(t *testing.T) {
	prefix, suffix := splitPrefix("fo.o")
	assert.Equal(t, prefix, "fo")
	assert.Equal(t, suffix, "o")
}
