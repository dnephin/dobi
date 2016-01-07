package tasks

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/buildpipe/config"
)

// VolumeTask is a task which creates a directory on the host
type VolumeTask struct {
	baseTask
	config *config.VolumeConfig
}

// NewVolumeTask creates a new VolumeTask object
func NewVolumeTask(options taskOptions, conf *config.VolumeConfig) *VolumeTask {
	return &VolumeTask{
		baseTask: baseTask{
			name:   options.name,
			client: options.client,
		},
		config: conf,
	}
}

func (t *VolumeTask) String() string {
	return fmt.Sprintf("VolumeTask(name=%s, config=%s)", t.name, t.config)
}

func (t *VolumeTask) logger() *log.Entry {
	return log.WithFields(log.Fields{
		"task":  "Volume",
		"name":  t.name,
		"path":  t.config.Path,
		"mount": t.config.Mount,
		"mode":  t.config.Mode,
	})
}

// Run creates the host path if it doesn't already exist
func (t *VolumeTask) Run(ctx *ExecuteContext) error {
	t.logger().Debug("run")

	if t.exists() {
		t.logger().Debug("exists")
		return nil
	}

	err := os.MkdirAll(t.config.Path, 0777)
	if err != nil {
		return err
	}
	ctx.setModified(t.name)
	t.logger().Info("created")
	return nil
}

func (t *VolumeTask) exists() bool {
	info, err := os.Stat(t.config.Path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
