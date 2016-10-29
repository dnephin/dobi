package job

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/dnephin/dobi/tasks/mount"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	"github.com/dnephin/dobi/utils/fs"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/go-connections/nat"
	docker "github.com/fsouza/go-dockerclient"
)

// DefaultUnixSocket to connect to the docker API
const DefaultUnixSocket = "/var/run/docker.sock"

func newRunTask(name task.Name, conf config.Resource) types.Task {
	return &Task{name: name, config: conf.(*config.JobConfig)}
}

// Task is a task which runs a command in a container to produce a
// file or set of files.
type Task struct {
	name   task.Name
	config *config.JobConfig
}

// Name returns the name of the task
func (t *Task) Name() task.Name {
	return t.name
}

func (t *Task) logger() *log.Entry {
	return logging.ForTask(t)
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
	return fmt.Sprintf("[job:run %v]%v", t.name.Resource(), buff.String())
}

// Run creates the host path if it doesn't already exist
func (t *Task) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	if !depsModified {
		stale, err := t.isStale(ctx)
		switch {
		case err != nil:
			return false, err
		case !stale:
			t.logger().Info("is fresh")
			return false, nil
		}
	}
	t.logger().Debug("is stale")

	t.logger().Info("Start")
	err := t.runContainer(ctx)
	if err != nil {
		return false, err
	}
	t.logger().Info("Done")
	return true, nil
}

func (t *Task) isStale(ctx *context.ExecuteContext) (bool, error) {
	if t.config.Artifact == "" {
		return true, nil
	}

	artifactLastModified, err := t.artifactLastModified()
	if err != nil {
		t.logger().Warnf("Failed to get artifact last modified: %s", err)
		return true, err
	}

	if len(t.config.Sources) != 0 {
		sourcesLastModified, err := fs.LastModified(t.config.Sources...)
		if err != nil {
			return true, err
		}
		if artifactLastModified.Before(sourcesLastModified) {
			t.logger().Debug("artifact older than sources")
			return true, nil
		}
		return false, nil
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
	taskImage, err := image.GetImage(ctx, imageName)
	if err != nil {
		return true, fmt.Errorf("Failed to get image %q: %s", imageName, err)
	}
	if artifactLastModified.Before(taskImage.Created) {
		t.logger().Debug("artifact older than image")
		return true, nil
	}
	return false, nil
}

func (t *Task) artifactLastModified() (time.Time, error) {
	// File or directory doesn't exist
	if _, err := os.Stat(t.config.Artifact); err != nil {
		return time.Time{}, nil
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
	name := ContainerName(ctx, t.name.Resource())
	container, err := ctx.Client.CreateContainer(t.createOptions(ctx, name))
	if err != nil {
		return fmt.Errorf("Failed creating container %q: %s", name, err)
	}

	chanSig := t.forwardSignals(ctx.Client, container.ID)
	defer signal.Stop(chanSig)
	defer RemoveContainer(t.logger(), ctx.Client, container.ID, true)

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

	portBinds, exposedPorts := asPortBindings(t.config.Ports)
	// TODO: only set Tty if running in a tty
	opts := docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Cmd:          t.config.Command.Value(),
			Image:        imageName,
			User:         t.config.User,
			OpenStdin:    interactive,
			Tty:          interactive,
			AttachStdin:  interactive,
			StdinOnce:    interactive,
			AttachStderr: true,
			AttachStdout: true,
			Env:          t.config.Env,
			Entrypoint:   t.config.Entrypoint.Value(),
			WorkingDir:   t.config.WorkingDir,
			ExposedPorts: exposedPorts,
		},
		HostConfig: &docker.HostConfig{
			Binds:        t.bindMounts(ctx),
			Privileged:   t.config.Privileged,
			NetworkMode:  t.config.NetMode,
			PortBindings: portBinds,
			Devices:      getDevices(t.config.Devices),
		},
	}
	opts = provideDocker(opts)
	return opts
}

func getDevices(devices []config.Device) []docker.Device {
	var dockerdevices []docker.Device
	for _, dev := range devices {
		if dev.Container == "" {
			dev.Container = dev.Host
		}
		if dev.Permissions == "" {
			dev.Permissions = "rwm"
		}
		dockerdevices = append(dockerdevices,
			docker.Device{
				PathInContainer:   dev.Container,
				PathOnHost:        dev.Host,
				CgroupPermissions: dev.Permissions,
			})
	}
	return dockerdevices
}

func asPortBindings(ports []string) (map[docker.Port][]docker.PortBinding, map[docker.Port]struct{}) {
	binds := make(map[docker.Port][]docker.PortBinding)
	exposed := make(map[docker.Port]struct{})
	for _, port := range ports {
		parts := strings.SplitN(port, ":", 2)
		proto, cport := nat.SplitProtoPort(parts[1])
		cport = cport + "/" + proto
		binds[docker.Port(cport)] = []docker.PortBinding{{HostPort: parts[0]}}
		exposed[docker.Port(cport)] = struct{}{}
	}
	return binds, exposed
}

func provideDocker(opts docker.CreateContainerOptions) docker.CreateContainerOptions {
	dockerHostEnv := os.Getenv("DOCKER_HOST")
	switch {
	case dockerHostEnv != "":
		opts.Config.Env = append(opts.Config.Env, "DOCKER_HOST="+dockerHostEnv)
	default:
		path := DefaultUnixSocket
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

// Stop the task
func (t *Task) Stop(ctx *context.ExecuteContext) error {
	return nil
}
