//go:build !windows
// +build !windows

package job

import (
	"os"
	"syscall"
)

// Send a SIGWINCH signal to make sure terminals to have the correct
// window dimensions
func initWindow(chanSig chan<- os.Signal) {
	chanSig <- syscall.SIGWINCH
}
