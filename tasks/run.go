package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/docker/docker/pkg/term"
	docker "github.com/fsouza/go-dockerclient"
)

// RunTask is a task which runs a command in a container to produce a
// file or set of files.
type RunTask struct {
	baseTask
	config *config.RunConfig
}

// NewRunTask creates a new RunTask object
func NewRunTask(options taskOptions, conf *config.RunConfig) *RunTask {
	return &RunTask{
		baseTask: baseTask{name: options.name},
		config:   conf,
	}
}

func (t *RunTask) String() string {
	return fmt.Sprintf("RunTask(name=%s, config=%s)", t.name, t.config)
}

func (t *RunTask) logger() *log.Entry {
	return log.WithFields(log.Fields{
		"task":     "Command",
		"name":     t.name,
		"use":      t.config.Use,
		"command":  t.config.Command,
		"artifact": t.config.Artifact,
	})
}

// Run creates the host path if it doesn't already exist
func (t *RunTask) Run(ctx *ExecuteContext) error {
	t.logger().Info("run")
	stale, err := t.isStale(ctx)
	if !stale || err != nil {
		return err
	}
	t.logger().Debug("artifact is stale")

	err = t.runContainer(ctx)
	if err != nil {
		return err
	}
	ctx.setModified(t.name)
	t.logger().Info("done")
	return nil
}

func (t *RunTask) isStale(ctx *ExecuteContext) (bool, error) {
	if ctx.isModified(t.config.Dependencies()...) {
		return true, nil
	}

	if t.config.Artifact == "" {
		return true, nil
	}

	artifactLastModified, err := t.artifactLastModified()
	if err != nil {
		return true, err
	}

	volumesLastModified, err := t.volumesLastModified(ctx)
	if err != nil {
		return true, err
	}

	if artifactLastModified.Before(volumesLastModified) {
		t.logger().Debug("artifact older than volume files")
		return true, nil
	}

	image, err := ctx.tasks.images[t.config.Use].GetImage(ctx)
	if err != nil {
		return true, err
	}
	if artifactLastModified.Before(image.Created) {
		t.logger().Debug("artifact older than image")
		return true, nil
	}
	return false, nil
}

func (t *RunTask) artifactLastModified() (time.Time, error) {
	info, err := os.Stat(t.config.Artifact)
	// File or directory doesn't exist
	if err != nil {
		return time.Time{}, nil
	}

	if !info.IsDir() {
		return info.ModTime(), nil
	}

	return lastModified(t.config.Artifact)
}

// TODO: support a .volumeignore file used to ignore mtime of files
func (t *RunTask) volumesLastModified(ctx *ExecuteContext) (time.Time, error) {
	volumePaths := []string{}
	ctx.tasks.EachVolume(t.config.Volumes, func(name string, volume *VolumeTask) {
		volumePaths = append(volumePaths, volume.config.Path)
	})
	return lastModified(volumePaths...)
}

func (t *RunTask) volumeBinds(ctx *ExecuteContext) []string {
	binds := []string{}
	ctx.tasks.EachVolume(t.config.Volumes, func(name string, volume *VolumeTask) {
		binds = append(binds, volume.asBind())
	})
	return binds
}

func (t *RunTask) runContainer(ctx *ExecuteContext) error {
	interactive := t.config.Interactive
	// TODO: support other run options
	container, err := ctx.client.CreateContainer(docker.CreateContainerOptions{
		Name: fmt.Sprintf("%s-%s", ctx.Env.Unique(), t.name),
		Config: &docker.Config{
			Cmd:          t.config.ParsedCommand(),
			Image:        ctx.tasks.images[t.config.Use].getImageName(ctx),
			OpenStdin:    interactive,
			Tty:          interactive,
			AttachStdin:  interactive,
			StdinOnce:    interactive,
			AttachStderr: true,
			AttachStdout: true,
		},
		HostConfig: &docker.HostConfig{
			Binds:      t.volumeBinds(ctx),
			Privileged: t.config.Privileged,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed creating container: %s", err)
	}

	chanSig := t.forwardSignals(ctx.client, container.ID)
	defer signal.Stop(chanSig)

	_, err = ctx.client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		InputStream:  ioutil.NopCloser(os.Stdin),
		Stream:       true,
		Stdin:        t.config.Interactive,
		RawTerminal:  t.config.Interactive,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		return fmt.Errorf("Failed attaching to container: %s", err)
	}

	if interactive {
		inFd, _ := term.GetFdInfo(os.Stdin)
		state, err := term.SetRawTerminal(inFd)
		if err != nil {
			return err
		}
		defer func() {
			if err := term.RestoreTerminal(inFd, state); err != nil {
				t.logger().Warnf("Failed to restore fd %v: %s", inFd, err)
			}
		}()
	}

	if err := ctx.client.StartContainer(container.ID, nil); err != nil {
		return fmt.Errorf("Failed starting container: %s", err)
	}

	t.waitAndRemove(ctx.client, container.ID)
	return nil
}

func (t *RunTask) waitAndRemove(client *docker.Client, containerID string) {
	status, err := client.WaitContainer(containerID)
	if err != nil {
		t.logger().Warnf("Failed to wait on container exit: %s", err)
	}
	if status != 0 {
		t.logger().WithFields(log.Fields{"status": status}).Warn(
			"Container exited with non-zero status code")
	}

	if err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerID,
		RemoveVolumes: true,
	}); err != nil {
		t.logger().WithFields(log.Fields{"container": containerID}).Warnf(
			"Failed to remove container: %s", err)
	}
}

func (t *RunTask) forwardSignals(client *docker.Client, containerID string) chan<- os.Signal {
	chanSig := make(chan os.Signal, 128)

	// TODO: not all of these exist on windows?
	signal.Notify(chanSig, syscall.SIGINT, syscall.SIGTERM)

	kill := func(sig os.Signal) {
		t.logger().WithFields(log.Fields{"signal": sig}).Debug("Received signal")

		intSig, ok := sig.(syscall.Signal)
		if !ok {
			t.logger().WithFields(log.Fields{"signal": sig}).Warnf(
				"Failed to convert signal from %T", sig)
			return
		}

		if err := client.KillContainer(docker.KillContainerOptions{
			ID:     containerID,
			Signal: docker.Signal(intSig),
		}); err != nil {
			t.logger().WithFields(log.Fields{"signal": sig}).Warnf(
				"Failed to send signal: %s", err)
		}
	}

	go func() {
		for sig := range chanSig {
			kill(sig)
		}
	}()
	return chanSig
}

// Prepare the task
func (t *RunTask) Prepare(ctx *ExecuteContext) error {
	return nil
}

// Stop the task
func (t *RunTask) Stop(ctx *ExecuteContext) error {
	return nil
}
