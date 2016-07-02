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

type MountTaskSuite struct {
	suite.Suite
	task *MountTask
	path string
	ctx  *ExecuteContext
}

func TestMountTaskSuite(t *testing.T) {
	suite.Run(t, new(MountTaskSuite))
}

func (s *MountTaskSuite) SetupTest() {
	var err error
	s.path, err = ioutil.TempDir("", "mount-task-test")
	s.Require().Nil(err)

	s.task = NewMountTask(
		taskOptions{
			name:   "mount-task-def",
			config: &config.Config{WorkingDir: s.path},
		},
		&config.MountConfig{
			Bind:     filepath.Join("a", "b", "c"),
			Path:     "/target",
			ReadOnly: false,
		})

	s.ctx = NewExecuteContext(nil, nil, nil, false)
}

func (s *MountTaskSuite) TearDownTest() {
	s.Nil(os.RemoveAll(s.path))
}

func (s *MountTaskSuite) TestRunPathExists() {
	s.False(s.task.exists())
	s.Require().Nil(os.MkdirAll(s.task.absBindPath(), 0777))
	s.True(s.task.exists())

	s.Nil(s.task.Run(s.ctx))
	s.False(s.ctx.isModified("mount-task-def"))
}

func (s *MountTaskSuite) TestRunPathIsNew() {
	s.Nil(s.task.Run(s.ctx))

	s.True(s.task.exists())
	s.True(s.ctx.isModified("mount-task-def"))
}

func (s *MountTaskSuite) TestAsBind() {
	s.Equal(fmt.Sprintf("%s:/target:rw", s.task.absBindPath()), s.task.asBind())
}
