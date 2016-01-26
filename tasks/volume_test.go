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

type VolumeTaskSuite struct {
	suite.Suite
	task *VolumeTask
	path string
	ctx  *ExecuteContext
}

func TestVolumeTaskSuite(t *testing.T) {
	suite.Run(t, new(VolumeTaskSuite))
}

func (s *VolumeTaskSuite) SetupTest() {
	var err error
	s.path, err = ioutil.TempDir("", "volume-task-test")
	s.Require().Nil(err)

	s.task = NewVolumeTask(
		taskOptions{
			name:   "volume-task-def",
			config: &config.Config{WorkingDir: s.path},
		},
		&config.VolumeConfig{
			Path:  filepath.Join("a", "b", "c"),
			Mount: "/target",
			Mode:  "rw",
		})

	s.ctx = NewExecuteContext(nil, nil, nil)
}

func (s *VolumeTaskSuite) TearDownTest() {
	s.Nil(os.RemoveAll(s.path))
}

func (s *VolumeTaskSuite) TestRunPathExists() {
	s.False(s.task.exists())
	s.Require().Nil(os.MkdirAll(s.task.absPath(), 0777))
	s.True(s.task.exists())

	s.Nil(s.task.Run(s.ctx))
	s.False(s.ctx.isModified("volume-task-def"))
}

func (s *VolumeTaskSuite) TestRunPathIsNew() {
	s.Nil(s.task.Run(s.ctx))

	s.True(s.task.exists())
	s.True(s.ctx.isModified("volume-task-def"))
}

func (s *VolumeTaskSuite) TestAsBind() {
	s.Equal(fmt.Sprintf("%s:/target:rw", s.task.absPath()), s.task.asBind())
}
