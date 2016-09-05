package job

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
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/common"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/dnephin/dobi/tasks/mount"
	"github.com/dnephin/dobi/utils/fs"
	dopts "github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/term"
	docker "github.com/fsouza/go-dockerclient"
)

// Task is a task which runs a command in a container to produce a
// file or set of files.
type Task struct {
	name   string
	config *config.JobConfig
}

// NewTask creates a new Task object
func NewTask(name string, conf *config.JobConfig) *Task {
	return &Task{name: name, config: conf}
}

// Name returns the name of the task
func (t *Task) Name() common.TaskName {
	return common.NewTaskName(t.name, "run")
}

func (t *Task) logger() *log.Entry {
	return logging.Log.WithFields(log.Fields{"task": t})
}

// Repr formats the task for logging
func (t *Task) Repr() string {
	buff := &bytes.Buffer{}

	if !t.config.Command.Empty() {
		buff.WriteString(" " + t.config.Command.String())
	}
	if !t.config.Command.Empty() && t.config.Artifact != "" {
		buff.WriteString(" ->")
	}
	if t.config.Artifact != "" {
		buff.WriteString(" " + t.config.Artifact)
	}
	return fmt.Sprintf("[run %v]%v", t.name, buff.String())
}

// Run creates the host path if it doesn't already exist
func (t *Task) Run(ctx *context.ExecuteContext) error {
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
		t.logger().Warnf("Failed to get artifact last modified: %s", err)
		return true, err
	}

	mountsLastModified, err := t.mountsLastModified(ctx)
	if err != nil {
		t.logger().Warnf("Failed to get mounts last modified: %s", err)
		return true, err
	}

	if artifactLastModified.Before(mountsLastModified) {
		t.logger().Debug("artifact older than mount files")
		return true, nil
	}

	imageName := ctx.Resources.Image(t.config.Use)
	image, err := image.GetImage(ctx, imageName)
	if err != nil {
		return true, fmt.Errorf("Failed to get image %q: %s", imageName, err)
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
	name := ContainerName(ctx, t.name)
	container, err := ctx.Client.CreateContainer(t.createOptions(ctx, name))
	if err != nil {
		return fmt.Errorf("Failed creating container %q: %s", name, err)
	}

	chanSig := t.forwardSignals(ctx.Client, container.ID)
	defer signal.Stop(chanSig)
	defer RemoveContainer(t.logger(), ctx.Client, container.ID)

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
		return fmt.Errorf("Failed attaching to container %q: %s", name, err)
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
		return fmt.Errorf("Failed starting container %q: %s", name, err)
	}

	return t.wait(ctx.Client, container.ID)
}

func (t *Task) createOptions(ctx *context.ExecuteContext, name string) docker.CreateContainerOptions {
	interactive := t.config.Interactive

	imageName := image.GetImageName(ctx, ctx.Resources.Image(t.config.Use))
	t.logger().Debugf("Image name %q", imageName)
	// TODO: only set Tty if running in a tty
	opts := docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Cmd:          t.config.Command.Value(),
			Image:        imageName,
			OpenStdin:    interactive,
			Tty:          interactive,
			AttachStdin:  interactive,
			StdinOnce:    interactive,
			AttachStderr: true,
			AttachStdout: true,
			Env:          t.config.Env,
			Entrypoint:   t.config.Entrypoint.Value(),
		},
		HostConfig: &docker.HostConfig{
			Binds:       t.bindMounts(ctx),
			Privileged:  t.config.Privileged,
			NetworkMode: t.config.NetMode,
		},
	}
	opts = provideDocker(opts)
	return opts
}

func provideDocker(opts docker.CreateContainerOptions) docker.CreateContainerOptions {
	dockerHostEnv := os.Getenv("DOCKER_HOST")
	switch {
	case dockerHostEnv != "":
		opts.Config.Env = append(opts.Config.Env, "DOCKER_HOST="+dockerHostEnv)
	default:
		path := dopts.DefaultUnixSocket
		opts.HostConfig.Binds = append(opts.HostConfig.Binds, path+":"+path)
	}
	return opts
}

func (t *Task) wait(client client.DockerClient, containerID string) error {
	status, err := client.WaitContainer(containerID)
	if err != nil {
		return fmt.Errorf("Failed to wait on container exit: %s", err)
	}
	if status != 0 {
		return fmt.Errorf("Exited with non-zero status code %d", status)
	}
	return nil
}

func (t *Task) forwardSignals(client client.DockerClient, containerID string) chan<- os.Signal {
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

// Dependencies returns the list of dependencies
func (t *Task) Dependencies() []string {
	return t.config.Dependencies()
}

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}
