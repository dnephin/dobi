package job

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/fsouza/go-dockerclient"
)

func (t *Task) runWithBuildAndCopy(ctx *context.ExecuteContext) error {
	imageID, err := t.buildImageWithMounts(ctx)
	if err != nil {
		return err
	}

	name := containerName(ctx, t.name.Resource())
	options := t.createOptions(ctx, name, imageID)
	defer removeContainerWithLogging(t.logger(), ctx.Client, name)

	if err := t.runContainer(ctx, options); err != nil {
		return err
	}
	return copyFilesToHost(ctx, name)
}

func (t *Task) buildImageWithMounts(ctx *context.ExecuteContext) (string, error) {
	dockerfile := buildDockerfileWithCopy(ctx, t.config)
	buildContext, err := buildTarFromDockerfile(dockerfile)
	if err != nil {
		return "", err
	}

	err = image.Stream(os.Stdout, func(out io.Writer) error {
		opts := buildImageOptions(ctx, out)
		opts.InputStream = buildContext
		return ctx.Client.BuildImage(opts)
	})
	// TODO: imageID
	return "", err
}

func buildDockerfileWithCopy(ctx *context.ExecuteContext, cfg *config.JobConfig) *bytes.Buffer {
	baseImage := image.GetImageName(ctx, ctx.Resources.Image(cfg.Use))
	buf := bytes.NewBufferString("FROM " + baseImage + "\n")
	// TODO: sort by shortest path first
	for _, mountName := range cfg.Mounts {
		mount := ctx.Resources.Mount(mountName)
		if !mount.IsBind() {
			continue
		}
		buf.WriteString(fmt.Sprintf("COPY %s %s\n", mount.Bind, mount.Path))
	}
	return buf
}

func buildTarFromDockerfile(dockerfile *bytes.Buffer) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)
	if err := tarWriter.WriteHeader(&tar.Header{
		Name: "Dockerfile",
		Mode: 644,
		Size: int64(dockerfile.Len()),
	}); err != nil {
		return nil, err
	}
	_, err := io.Copy(tarWriter, dockerfile)
	return buf, err
}

func buildImageOptions(ctx *context.ExecuteContext, out io.Writer) docker.BuildImageOptions {
	return docker.BuildImageOptions{
		RmTmpContainer: true,
		OutputStream:   out,
		RawJSONStream:  true,
		SuppressOutput: ctx.Settings.Quiet,
		AuthConfigs:    ctx.GetAuthConfigs(),
		Dockerfile:     "Dockerfile",
	}
}

func copyFilesToHost(ctx *context.ExecuteContext, containerID string) error {
	return nil
}
