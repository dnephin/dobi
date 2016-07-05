package run

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/dnephin/dobi/tasks/mount"
	"github.com/dnephin/dobi/utils/fs"
	"github.com/docker/docker/pkg/term"
	docker "github.com/fsouza/go-dockerclient"
)

// Task is a task which runs a command in a container to produce a
// file or set of files.
type Task struct {
	name   string
	config *config.RunConfig
}

// NewTask creates a new Task object
func NewTask(name string, conf *config.RunConfig) *Task {
	return &Task{name: name, config: conf}
}

// Name returns the name of the task
func (t *Task) Name() string {
	return t.name
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	buff := &bytes.Buffer{}

	if t.config.Command != "" {
		buff.WriteString(" " + t.config.Command)
	}
	if t.config.Command != "" && t.config.Artifact != "" {
		buff.WriteString(" ->")
	}
	if t.config.Artifact != "" {
		buff.WriteString(" " + t.config.Artifact)
	}
	return fmt.Sprintf("[run %v]%v", t.name, buff.String())
}

// Run creates the host path if it doesn't already exist
func (t *Task) Run(ctx *context.ExecuteContext) error {
	t.logger().Debug("Run")
	stale, err := t.isStale(ctx)
	if !stale || err != nil {
		t.logger().Info("is fresh")
		return err
	}
	t.logger().Debug("is stale")

	t.logger().Info("Start")
	err = t.runContainer(ctx)
	if err != nil {
		return err
	}
	ctx.SetModified(t.name)
	t.logger().Info("Done")
	return nil
}

func (t *Task) isStale(ctx *context.ExecuteContext) (bool, error) {
	if ctx.IsModified(t.config.Dependencies()...) {
		return true, nil
	}

	if t.config.Artifact == "" {
		return true, nil
	}

	artifactLastModified, err := t.artifactLastModified()
	if err != nil {
		return true, err
	}

	mountsLastModified, err := t.mountsLastModified(ctx)
	if err != nil {
		return true, err
	}

	if artifactLastModified.Before(mountsLastModified) {
		t.logger().Debug("artifact older than mount files")
		return true, nil
	}

	image, err := image.GetImage(ctx, ctx.Resources.Image(t.config.Use))
	if err != nil {
		return true, err
	}
	if artifactLastModified.Before(image.Created) {
		t.logger().Debug("artifact older than image")
		return true, nil
	}
	return false, nil
}

func (t *Task) artifactLastModified() (time.Time, error) {
	info, err := os.Stat(t.config.Artifact)
	// File or directory doesn't exist
	if err != nil {
		return time.Time{}, nil
	}

	if !info.IsDir() {
		return info.ModTime(), nil
	}

	return fs.LastModified(t.config.Artifact)
}

// TODO: support a .mountignore file used to ignore mtime of files
func (t *Task) mountsLastModified(ctx *context.ExecuteContext) (time.Time, error) {
	mountPaths := []string{}
	ctx.Resources.EachMount(t.config.Mounts, func(name string, mount *config.MountConfig) {
		mountPaths = append(mountPaths, mount.Bind)
	})
	return fs.LastModified(mountPaths...)
}

func (t *Task) bindMounts(ctx *context.ExecuteContext) []string {
	binds := []string{}
	ctx.Resources.EachMount(t.config.Mounts, func(name string, config *config.MountConfig) {
		binds = append(binds, mount.AsBind(config, ctx.WorkingDir))
	})
	return binds
}

func (t *Task) runContainer(ctx *context.ExecuteContext) error {
	interactive := t.config.Interactive
	// TODO: support other run options
	container, err := ctx.Client.CreateContainer(docker.CreateContainerOptions{
		Name: fmt.Sprintf("%s-%s", ctx.Env.Unique(), t.name),
		Config: &docker.Config{
			Cmd:          t.config.ParsedCommand(),
			Image:        image.GetImageName(ctx, ctx.Resources.Image(t.config.Use)),
			OpenStdin:    interactive,
			Tty:          interactive,
			AttachStdin:  interactive,
			StdinOnce:    interactive,
			AttachStderr: true,
			AttachStdout: true,
			Env:          envWithVars(ctx.Env, t.config.Env),
		},
		HostConfig: &docker.HostConfig{
			Binds:      t.bindMounts(ctx),
			Privileged: t.config.Privileged,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed creating container: %s", err)
	}

	chanSig := t.forwardSignals(ctx.Client, container.ID)
	defer signal.Stop(chanSig)
	defer t.remove(ctx.Client, container.ID)

	_, err = ctx.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
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

	if err := ctx.Client.StartContainer(container.ID, nil); err != nil {
		return fmt.Errorf("Failed starting container: %s", err)
	}

	return t.wait(ctx.Client, container.ID)
}

func (t *Task) wait(client *docker.Client, containerID string) error {
	status, err := client.WaitContainer(containerID)
	if err != nil {
		return fmt.Errorf("Failed to wait on container exit: %s", err)
	}
	if status != 0 {
		return fmt.Errorf("Exited with non-zero status code %d", status)
	}
	return nil
}

func (t *Task) remove(client *docker.Client, containerID string) {
	t.logger().Debug("Removing container")
	if err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerID,
		RemoveVolumes: true,
	}); err != nil {
		t.logger().WithFields(log.Fields{"container": containerID}).Warnf(
			"Failed to remove container: %s", err)
	}
}

func (t *Task) forwardSignals(client *docker.Client, containerID string) chan<- os.Signal {
	chanSig := make(chan os.Signal, 128)

	// TODO: not all of these exist on windows?
	signal.Notify(chanSig, syscall.SIGINT, syscall.SIGTERM)

	kill := func(sig os.Signal) {
		t.logger().WithFields(log.Fields{"signal": sig}).Debug("received")

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
func (t *Task) Prepare(ctx *context.ExecuteContext) error {
	for _, env := range t.config.Env {
		if _, err := ctx.Env.Resolve(env); err != nil {
			return err
		}
	}
	return nil
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}

func envWithVars(execEnv *context.ExecEnv, envs []string) []string {
	out := []string{}
	for _, env := range envs {
		out = append(out, execEnv.GetVar(env))
	}
	return out
}
