package mount

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/stretchr/testify/suite"
)

type CreateTaskSuite struct {
	suite.Suite
	task *CreateTask
	path string
	ctx  *context.ExecuteContext
}

func TestCreateTaskSuite(t *testing.T) {
	suite.Run(t, new(CreateTaskSuite))
}

func (s *CreateTaskSuite) SetupTest() {
	var err error
	s.path, err = ioutil.TempDir("", "mount-task-test")
	s.Require().Nil(err)

	s.task = NewCreateTask(
		"mount-task-def",
		&config.MountConfig{
			Bind:     filepath.Join("a", "b", "c"),
			Path:     "/target",
			ReadOnly: false,
		})

	s.ctx = context.NewExecuteContext(
		&config.Config{WorkingDir: s.path}, nil, nil, false)
}

func (s *CreateTaskSuite) TearDownTest() {
	s.Nil(os.RemoveAll(s.path))
}

func (s *CreateTaskSuite) TestRunPathExists() {
	s.False(s.task.exists(s.ctx))
	s.Require().Nil(os.MkdirAll(AbsBindPath(s.task.config, s.path), 0777))
	s.True(s.task.exists(s.ctx))

	s.Nil(s.task.Run(s.ctx))
	s.False(s.ctx.IsModified("mount-task-def"))
}

func (s *CreateTaskSuite) TestRunPathIsNew() {
	s.Nil(s.task.Run(s.ctx))

	s.True(s.task.exists(s.ctx))
	s.True(s.ctx.IsModified("mount-task-def"))
}

func (s *CreateTaskSuite) TestAsBind() {
	s.Equal(
		fmt.Sprintf("%s:/target:rw", AbsBindPath(s.task.config, s.path)),
		AsBind(s.task.config, s.path),
	)
}
