package compose

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dnephin/dobi/tasks/context"
	log "github.com/sirupsen/logrus"
)

// RunUpAttached starts the Compose project
func RunUpAttached(ctx *context.ExecuteContext, t *Task) error {
	t.logger().Info("project up")

	cmd := t.buildCommand("up", "-t", t.config.StopGraceString())
	if err := cmd.Start(); err != nil {
		return err
	}

	chanSig := forwardSignals(t, cmd.Process)
	defer signal.Stop(chanSig)

	if err := cmd.Wait(); err != nil {
		return err
	}
	t.logger().Info("Done")
	return nil
}

func forwardSignals(t *Task, proc *os.Process) chan<- os.Signal {
	chanSig := make(chan os.Signal, 128)

	// TODO: not all of these exist on windows?
	signal.Notify(chanSig, syscall.SIGINT, syscall.SIGTERM)

	kill := func(sig os.Signal) {
		t.logger().WithFields(log.Fields{"signal": sig}).Debug("received")

		if err := proc.Signal(sig); err != nil {
			t.logger().WithFields(log.Fields{"pid": proc.Pid}).Warnf(
				"failed to signal process: %s", err)
		}
	}

	go func() {
		for sig := range chanSig {
			kill(sig)
		}
	}()
	return chanSig
}
