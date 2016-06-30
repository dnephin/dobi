package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
)

// VolumeTask is a task which creates a directory on the host
type VolumeTask struct {
	baseTask
	config     *config.VolumeConfig
	workingDir string
}

// NewVolumeTask creates a new VolumeTask object
func NewVolumeTask(options taskOptions, conf *config.VolumeConfig) *VolumeTask {
	return &VolumeTask{
		baseTask:   baseTask{name: options.name},
		config:     conf,
		workingDir: options.config.WorkingDir,
	}
}

func (t *VolumeTask) String() string {
	return fmt.Sprintf("VolumeTask(name=%s, config=%s)", t.name, t.config)
}

func (t *VolumeTask) logger() *log.Entry {
	return log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *VolumeTask) Repr() string {
	return fmt.Sprintf("[volume %s] %s:%s", t.name, t.config.Path, t.config.Mount)
}

// Run creates the host path if it doesn't already exist
func (t *VolumeTask) Run(ctx *ExecuteContext) error {
	t.logger().Debug("Run")

	if t.exists() {
		t.logger().Debug("is fresh")
		return nil
	}

	err := os.MkdirAll(t.absPath(), 0777)
	if err != nil {
		return err
	}
	ctx.setModified(t.name)
	t.logger().Info("Created")
	return nil
}

func (t *VolumeTask) absPath() string {
	if filepath.IsAbs(t.config.Path) {
		return t.config.Path
	}
	return filepath.Join(t.workingDir, t.config.Path)
}

func (t *VolumeTask) exists() bool {
	_, err := os.Stat(t.absPath())
	if err != nil {
		return false
	}

	return true
}

func (t *VolumeTask) asBind() string {
	return fmt.Sprintf("%s:%s:%s", t.absPath(), t.config.Mount, t.config.Mode)
}

// Prepare the task
func (t *VolumeTask) Prepare(ctx *ExecuteContext) error {
	return nil
}

// Stop the task
func (t *VolumeTask) Stop(ctx *ExecuteContext) error {
	return nil
}
