package mount

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/stretchr/testify/suite"
)

type CreateTaskSuite struct {
	suite.Suite
	task   *Task
	action *createAction
	path   string
	ctx    *context.ExecuteContext
}

func TestCreateTaskSuite(t *testing.T) {
	suite.Run(t, new(CreateTaskSuite))
}

func (s *CreateTaskSuite) SetupTest() {
	var err error
	s.path, err = ioutil.TempDir("", "mount-task-test")
	s.Require().Nil(err)

	s.task = &Task{
		name: common.NewTaskName("resource", "action"),
		config: &config.MountConfig{
			Bind:     filepath.Join("a", "b", "c"),
			Path:     "/target",
			ReadOnly: false,
		},
		run: runCreate,
	}
	s.action = &createAction{task: s.task}

	s.ctx = context.NewExecuteContext(
		&config.Config{WorkingDir: s.path}, nil, nil, false)
}

func (s *CreateTaskSuite) TearDownTest() {
	s.Nil(os.RemoveAll(s.path))
}

func (s *CreateTaskSuite) TestRunPathExists() {
	s.False(s.action.exists(s.ctx))
	s.Require().Nil(os.MkdirAll(AbsBindPath(s.task.config, s.path), 0777))
	s.True(s.action.exists(s.ctx))

	modified, err := s.task.Run(s.ctx, false)
	s.Nil(err)
	s.False(modified)
}

func (s *CreateTaskSuite) TestRunPathIsNew() {
	modified, err := s.task.Run(s.ctx, false)
	s.Nil(err)
	s.True(modified)
	s.True(s.action.exists(s.ctx))
}

func (s *CreateTaskSuite) TestAsBind() {
	s.Equal(
		fmt.Sprintf("%s:/target:rw", AbsBindPath(s.task.config, s.path)),
		AsBind(s.task.config, s.path),
	)
}
