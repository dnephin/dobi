package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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

func (s *ExecEnvSuite) TestNewExecEnvDefault() {
	defer os.Setenv("USER", os.Getenv("USER"))
	os.Setenv("USER", "testuser")

	execEnv, err := NewExecEnv(s.cfg)
	s.Nil(err)
	expected := fmt.Sprintf("%s-testuser", filepath.Base(s.tmpDir))
	s.Equal(expected, execEnv.Unique())
}

func (s *ExecEnvSuite) TestNewExecEnvWithCommand() {
	s.cfg.Meta.UniqueExecID = "echo Use-This"

	execEnv, err := NewExecEnv(s.cfg)
	s.Nil(err)
	s.Equal("Use-This", execEnv.ExecID)
}

func (s *ExecEnvSuite) TestNewExecEnvWithInvalidCommand() {
	s.cfg.Meta.UniqueExecID = "bogus Use-This"

	_, err := NewExecEnv(s.cfg)
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
