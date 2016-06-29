package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dnephin/dobi/config"
	"github.com/stretchr/testify/suite"
)

type ExecEnvSuite struct {
	suite.Suite
	cfg    *config.Config
	tmpDir string
}

func TestExecEnvSuite(t *testing.T) {
	suite.Run(t, new(ExecEnvSuite))
}

func (s *ExecEnvSuite) SetupTest() {
	var err error
	s.tmpDir, err = ioutil.TempDir("", "environment-test")
	s.Require().Nil(err)

	s.cfg = &config.Config{
		Meta:       &config.MetaConfig{},
		WorkingDir: s.tmpDir,
	}
}

func (s *ExecEnvSuite) TearDownTest() {
	s.Nil(os.RemoveAll(s.tmpDir))
}

func (s *ExecEnvSuite) TestNewExecEnvFromConfigDefault() {
	defer os.Setenv("USER", os.Getenv("USER"))
	os.Setenv("USER", "testuser")

	execEnv, err := NewExecEnvFromConfig(s.cfg)
	s.Nil(err)
	expected := fmt.Sprintf("%s-testuser", filepath.Base(s.tmpDir))
	s.Equal(expected, execEnv.Unique())
}

func (s *ExecEnvSuite) TestNewExecEnvFromConfigWithCommand() {
	s.cfg.Meta.UniqueExecID = "echo Use-This"

	execEnv, err := NewExecEnvFromConfig(s.cfg)
	s.Nil(err)
	s.Equal("Use-This", execEnv.ExecID)
}

func (s *ExecEnvSuite) TestNewExecEnvFromConfigWithInvalidCommand() {
	s.cfg.Meta.UniqueExecID = "bogus Use-This"

	_, err := NewExecEnvFromConfig(s.cfg)
	s.Error(err)
	s.Contains(err.Error(), "\"bogus\": executable file not found")
}

func (s *ExecEnvSuite) TestValidateExecIDEmpty() {
	output, err := validateExecID("")
	s.Equal("", output)
	s.Error(err)
	s.Contains(err.Error(), "no output")
}

func (s *ExecEnvSuite) TestValidateExecIDTooManyLines() {
	output, err := validateExecID("one\ntwo")
	s.Equal("", output)
	s.Error(err)
	s.Contains(err.Error(), "returned 2 lines")

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
	execEnv := NewExecEnv("exec", "project")
	tmpl := "ok-{unique}"
	expected := "ok-" + execEnv.Unique()
	value, err := execEnv.Resolve(tmpl)

	s.Nil(err)
	s.Equal(value, expected)
	s.Equal(execEnv.tmplCache[tmpl], expected)
}

func (s *ExecEnvSuite) TestResolveUnknown() {
	execEnv := NewExecEnv("exec", "project")
	_, err := execEnv.Resolve("{bogus}")

	s.Error(err)
	s.Contains(err.Error(), "Unknown variable \"bogus\"")
}

func (s *ExecEnvSuite) TestResolveBadTemplate() {
	execEnv := NewExecEnv("exec", "project")
	_, err := execEnv.Resolve("{bogus{")

	s.Error(err)
	s.Contains(err.Error(), "Cannot find end tag")
}

func (s *ExecEnvSuite) TestResolveEnvironmentNoDefault() {
	execEnv := NewExecEnv("exec", "project")
	_, err := execEnv.Resolve("thing-{env.foo}")

	s.Error(err)
	s.Contains(err.Error(), "required for variable \"env.foo\"")
}

func (s *ExecEnvSuite) TestResolveEnvironment() {
	defer os.Unsetenv("FOO")
	os.Setenv("FOO", "stars")
	tmpl := "thing-{env.FOO}"
	expected := "thing-stars"

	execEnv := NewExecEnv("exec", "project")
	value, err := execEnv.Resolve(tmpl)

	s.Nil(err)
	s.Equal(value, expected)
	s.Equal(execEnv.tmplCache[tmpl], expected)
}

func (s *ExecEnvSuite) TestResolveTime() {
	tmpl := "build-{time.YYYY-MM-DD}"
	expected := "build-2016-04-05"

	execEnv := NewExecEnv("exec", "project")
	execEnv.startTime = time.Date(2016, 4, 5, 0, 0, 0, 0, time.UTC)
	value, err := execEnv.Resolve(tmpl)

	s.Nil(err)
	s.Equal(value, expected)
	s.Equal(execEnv.tmplCache[tmpl], expected)
}

func (s *ExecEnvSuite) TestSplitDefault() {
	tag := "time.19:01:01:default"
	value, defVal := splitDefault(tag)
	s.Equal(value, "time.19:01:01")
	s.Equal(defVal, "default")
}

func (s *ExecEnvSuite) TestSplitDefaultNoDefault() {
	tag := "env.FOO"
	value, defVal := splitDefault(tag)
	s.Equal(value, "env.FOO")
	s.Equal(defVal, "")
}
