package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
)

// MountTask is a task which creates a directory on the host
type MountTask struct {
	baseTask
	config     *config.MountConfig
	workingDir string
}

// NewMountTask creates a new MountTask object
func NewMountTask(options taskOptions, conf *config.MountConfig) *MountTask {
	return &MountTask{
		baseTask:   baseTask{name: options.name},
		config:     conf,
		workingDir: options.config.WorkingDir,
	}
}

func (t *MountTask) String() string {
	return fmt.Sprintf("MountTask(name=%s, config=%s)", t.name, t.config)
}

func (t *MountTask) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *MountTask) Repr() string {
	return fmt.Sprintf("[mount %s] %s:%s", t.name, t.config.Bind, t.config.Path)
}

// Run creates the host path if it doesn't already exist
func (t *MountTask) Run(ctx *ExecuteContext) error {
	t.logger().Debug("Run")

	if t.exists() {
		t.logger().Debug("is fresh")
		return nil
	}

	err := os.MkdirAll(t.absBindPath(), 0777)
	if err != nil {
		return err
	}
	ctx.setModified(t.name)
	t.logger().Info("Created")
	return nil
}

func (t *MountTask) absBindPath() string {
	if filepath.IsAbs(t.config.Bind) {
		return t.config.Bind
	}
	return filepath.Join(t.workingDir, t.config.Bind)
}

func (t *MountTask) exists() bool {
	_, err := os.Stat(t.absBindPath())
	if err != nil {
		return false
	}

	return true
}

func (t *MountTask) asBind() string {
	var mode string
	if t.config.ReadOnly {
		mode = "ro"
	} else {
		mode = "rw"
	}
	return fmt.Sprintf("%s:%s:%s", t.absBindPath(), t.config.Path, mode)
}

// Prepare the task
func (t *MountTask) Prepare(ctx *ExecuteContext) error {
	return nil
}

// Stop the task
func (t *MountTask) Stop(ctx *ExecuteContext) error {
	return nil
}
