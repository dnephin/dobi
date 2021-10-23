package job

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/dnephin/dobi/tasks/mount"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	"github.com/dnephin/dobi/utils/fs"
	"github.com/docker/go-connections/nat"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/moby/term"
	log "github.com/sirupsen/logrus"
)

// DefaultUnixSocket to connect to the docker API
const DefaultUnixSocket = "/var/run/docker.sock"

func newRunTask(name task.Name, conf config.Resource) types.Task {
	return &Task{name: name, config: conf.(*config.JobConfig)}
}

// Task is a task which runs a command in a container to produce a
// file or set of files.
type Task struct {
	types.NoStop
	name      task.Name
	config    *config.JobConfig
	outStream io.Writer
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
	if !t.config.Command.Empty() && !t.config.Artifact.Empty() {
		buff.WriteString(" ->")
	}
	if !t.config.Artifact.Empty() {
		buff.WriteString(" " + t.config.Artifact.String())
	}
	return fmt.Sprintf("%s%v", t.name.Format("job"), buff.String())
}

// Run the job command in a container
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
	var err error
	if ctx.Settings.BindMount {
		err = t.runContainerWithBinds(ctx)
	} else {
		err = t.runWithBuildAndCopy(ctx)
	}
	if err != nil {
		return false, err
	}
	t.logger().Info("Done")
	return true, nil
}

// nolint: gocyclo
func (t *Task) isStale(ctx *context.ExecuteContext) (bool, error) {
	if t.config.Artifact.Empty() {
		return true, nil
	}

	artifactLastModified, err := t.artifactLastModified(ctx.WorkingDir)
	if err != nil {
		t.logger().Warnf("Failed to get artifact last modified: %s", err)
		return true, err
	}

	if t.config.Sources.NoMatches() {
		t.logger().Warnf("No sources found matching: %s", &t.config.Sources)
		return true, nil
	}

	if len(t.config.Sources.Paths()) != 0 {
		sourcesLastModified, err := fs.LastModified(&fs.LastModifiedSearch{
			Root:  ctx.WorkingDir,
			Paths: t.config.Sources.Paths(),
		})
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
		return true, fmt.Errorf("failed to get image %q: %s", imageName, err)
	}
	if artifactLastModified.Before(taskImage.Created) {
		t.logger().Debug("artifact older than image")
		return true, nil
	}
	return false, nil
}

func (t *Task) artifactLastModified(workDir string) (time.Time, error) {
	paths := t.config.Artifact.Paths()
	// File or directory doesn't exist
	if len(paths) == 0 {
		return time.Time{}, nil
	}
	return fs.LastModified(&fs.LastModifiedSearch{Root: workDir, Paths: paths})
}

// TODO: support a .mountignore file used to ignore mtime of files
func (t *Task) mountsLastModified(ctx *context.ExecuteContext) (time.Time, error) {
	mountPaths := []string{}
	ctx.Resources.EachMount(t.config.Mounts, func(name string, mount *config.MountConfig) {
		mountPaths = append(mountPaths, mount.Bind)
	})
	return fs.LastModified(&fs.LastModifiedSearch{Root: ctx.WorkingDir, Paths: mountPaths})
}

func (t *Task) runContainerWithBinds(ctx *context.ExecuteContext) error {
	name := containerName(ctx, t.name.Resource())
	imageName := image.GetImageName(ctx, ctx.Resources.Image(t.config.Use))
	options := t.createOptions(ctx, name, imageName)

	defer removeContainerWithLogging(t.logger(), ctx.Client, name)
	return t.runContainer(ctx, options)
}

func removeContainerWithLogging(
	logger *log.Entry,
	client client.DockerClient,
	containerID string,
) {
	removed, err := removeContainer(logger, client, containerID)
	if !removed && err == nil {
		logger.WithFields(log.Fields{"container": containerID}).Warn(
			"Container does not exist")
	}
}

func (t *Task) runContainer(
	ctx *context.ExecuteContext,
	options docker.CreateContainerOptions,
) error {
	name := options.Name
	container, err := ctx.Client.CreateContainer(options)
	if err != nil {
		return fmt.Errorf("failed creating container %q: %s", name, err)
	}

	chanSig := t.forwardSignals(ctx.Client, container.ID)
	defer signal.Stop(chanSig)

	closeWaiter, err := ctx.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: t.output(),
		ErrorStream:  os.Stderr,
		InputStream:  ioutil.NopCloser(os.Stdin),
		Stream:       true,
		Stdin:        t.config.Interactive,
		RawTerminal:  t.config.Interactive,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		return fmt.Errorf("failed attaching to container %q: %s", name, err)
	}
	defer closeWaiter.Wait() // nolint: errcheck

	if t.config.Interactive {
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
		return fmt.Errorf("failed starting container %q: %s", name, err)
	}

	initWindow(chanSig)
	return t.wait(ctx.Client, container.ID)
}

func (t *Task) output() io.Writer {
	if t.outStream == nil {
		return os.Stdout
	}
	return io.MultiWriter(t.outStream, os.Stdout)
}

func (t *Task) createOptions(
	ctx *context.ExecuteContext,
	name string,
	imageName string,
) docker.CreateContainerOptions {
	t.logger().Debugf("Image name %q", imageName)

	interactive := t.config.Interactive
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
			Labels:       t.config.Labels,
			AttachStderr: true,
			AttachStdout: true,
			Env:          t.config.Env,
			Entrypoint:   t.config.Entrypoint.Value(),
			WorkingDir:   t.config.WorkingDir,
			ExposedPorts: exposedPorts,
		},
		HostConfig: &docker.HostConfig{
			Binds:        getMountsForHostConfig(ctx, t.config.Mounts),
			Privileged:   t.config.Privileged,
			NetworkMode:  t.config.NetMode,
			PortBindings: portBinds,
			Devices:      getDevices(t.config.Devices),
		},
	}
	if t.config.ProvideDocker {
		opts = provideDocker(opts)
	}
	return opts
}

func getMountsForHostConfig(ctx *context.ExecuteContext, mounts []string) []string {
	binds := []string{}
	ctx.Resources.EachMount(mounts, func(name string, mountConfig *config.MountConfig) {
		if !ctx.Settings.BindMount && mountConfig.IsBind() {
			return
		}
		binds = append(binds, mount.AsBind(mountConfig, ctx.WorkingDir))
	})
	return binds
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

func asPortBindings(ports []string) (map[docker.Port][]docker.PortBinding, map[docker.Port]struct{}) { // nolint: lll
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
	if os.Getenv("DOCKER_HOST") == "" {
		path := DefaultUnixSocket
		opts.HostConfig.Binds = append(opts.HostConfig.Binds, path+":"+path)
	}
	for _, envVar := range os.Environ() {
		if strings.HasPrefix(envVar, "DOCKER_") {
			opts.Config.Env = append(opts.Config.Env, envVar)
		}
	}
	return opts
}

func (t *Task) wait(client client.DockerClient, containerID string) error {
	status, err := client.WaitContainer(containerID)
	if err != nil {
		return fmt.Errorf("failed to wait on container exit: %s", err)
	}
	if status != 0 {
		return fmt.Errorf("exited with non-zero status code %d", status)
	}
	return nil
}

func (t *Task) forwardSignals(
	client client.DockerClient,
	containerID string,
) chan<- os.Signal {
	chanSig := make(chan os.Signal, 128)

	signal.Notify(chanSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)

	go func() {
		for sig := range chanSig {
			logger := t.logger().WithField("signal", sig)
			logger.Debug("received")

			sysSignal, ok := sig.(syscall.Signal)
			if !ok {
				logger.Warnf("Failed to convert signal from %T", sig)
				return
			}

			switch sysSignal {
			case syscall.SIGWINCH:
				handleWinSizeChangeSignal(logger, client, containerID)
			default:
				handleShutdownSignals(logger, client, containerID, sysSignal)
			}
		}
	}()
	return chanSig
}

func handleWinSizeChangeSignal(
	logger log.FieldLogger,
	client client.DockerClient,
	containerID string,
) {
	winsize, err := term.GetWinsize(os.Stdin.Fd())
	if err != nil {
		logger.WithError(err).
			Error("Failed to get host's TTY window size")
		return
	}

	err = client.ResizeContainerTTY(containerID, int(winsize.Height), int(winsize.Width))
	if err != nil {
		logger.WithError(err).
			Warning("Failed to set container's TTY window size.")
	}
}

func handleShutdownSignals(
	logger log.FieldLogger,
	client client.DockerClient,
	containerID string,
	sig syscall.Signal,
) {
	if err := client.KillContainer(docker.KillContainerOptions{
		ID:     containerID,
		Signal: docker.Signal(sig),
	}); err != nil {
		logger.WithError(err).
			Warn("Failed to forward signal")
	}
}
