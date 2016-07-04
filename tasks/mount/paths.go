package mount

import (
	"fmt"
	"path/filepath"

	"github.com/dnephin/dobi/config"
)

// AsBind returns a MountConfig formatted as a bind mount string
func AsBind(c *config.MountConfig, workingDir string) string {
	var mode string
	if c.ReadOnly {
		mode = "ro"
	} else {
		mode = "rw"
	}
	return fmt.Sprintf("%s:%s:%s", AbsBindPath(c, workingDir), c.Path, mode)
}

// AbsBindPath returns the MountConfig.Bind as an absolute path
func AbsBindPath(c *config.MountConfig, workingDir string) string {
	if filepath.IsAbs(c.Bind) {
		return c.Bind
	}
	return filepath.Join(workingDir, c.Bind)
}
