package job

import (
	"os"
	"syscall"
)

// SIGWINCH does not exist on windows, create a fake signal
const SIGWINCH = syscall.Signal(0xffffff)

func initWindow(chanSig chan<- os.Signal) {
}
