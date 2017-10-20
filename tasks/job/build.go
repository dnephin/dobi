package job

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Microsoft/go-winio/archive/tar"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/pkg/archive"
	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sort"
)

func (t *Task) runWithBuildAndCopy(ctx *context.ExecuteContext) error {
	name := containerName(ctx, t.name.Resource())
	imageName := fmt.Sprintf("%s:job-%s",
		ctx.Resources.Image(t.config.Use).Image, name)

	if err := t.buildImageWithMounts(ctx, imageName); err != nil {
		return err
	}
	defer removeImage(t.logger(), ctx.Client, imageName)

	defer removeContainerWithLogging(t.logger(), ctx.Client, name)
	options := t.createOptions(ctx, name, imageName)
	if err := t.runContainer(ctx, options); err != nil {
		return err
	}
	return copyFilesToHost(t.logger(), ctx, t.config, name)
}

func (t *Task) buildImageWithMounts(ctx *context.ExecuteContext, imageName string) error {
	baseImage := image.GetImageName(ctx, ctx.Resources.Image(t.config.Use))
	mounts := getBindMounts(ctx, t.config)

	dockerfile := buildDockerfileWithCopy(baseImage, mounts)
	buildContext, dockerfileName, err := buildTarContext(dockerfile, mounts)
	if err != nil {
		return err
	}
	return image.Stream(os.Stdout, func(out io.Writer) error {
		opts := buildImageOptions(ctx, out)
		opts.InputStream = buildContext
		opts.Name = imageName
		opts.Dockerfile = dockerfileName
		return ctx.Client.BuildImage(opts)
	})
}

func getBindMounts(ctx *context.ExecuteContext, cfg *config.JobConfig) []config.MountConfig {
	mounts := []config.MountConfig{}
	ctx.Resources.EachMount(cfg.Mounts, func(_ string, mount *config.MountConfig) {
		if !mount.IsBind() {
			return
		}
		mounts = append(mounts, *mount)
	})
	return mounts
}

func buildDockerfileWithCopy(baseImage string, mounts []config.MountConfig) *bytes.Buffer {
	buf := bytes.NewBufferString("FROM " + baseImage + "\n")
	// TODO: sort by shortest path first
	for _, mount := range mounts {
		buf.WriteString(fmt.Sprintf("COPY %s %s\n", mount.Bind, mount.Path))
	}
	return buf
}

func buildTarContext(
	dockerfile io.Reader,
	mounts []config.MountConfig,
) (io.Reader, string, error) {
	paths := []string{}
	for _, mount := range mounts {
		paths = append(paths, mount.Bind)
	}
	buildCtx, err := archive.TarWithOptions(".", &archive.TarOptions{
		IncludeFiles: paths,
	})
	if err != nil {
		return nil, "", err
	}
	return build.AddDockerfileToBuildContext(ioutil.NopCloser(dockerfile), buildCtx)
}

func buildImageOptions(ctx *context.ExecuteContext, out io.Writer) docker.BuildImageOptions {
	return docker.BuildImageOptions{
		RmTmpContainer: true,
		OutputStream:   out,
		RawJSONStream:  true,
		SuppressOutput: ctx.Settings.Quiet,
		AuthConfigs:    ctx.GetAuthConfigs(),
	}
}

func removeImage(logger *log.Entry, client client.DockerClient, imageID string) {
	if err := client.RemoveImage(imageID); err != nil {
		logger.Warnf("failed to remove %q: %s", imageID, err)
	}
}

// TODO: optimize by performing only one copy from container per directory
// if there are overlapping artifact paths
func copyFilesToHost(
	logger *log.Entry,
	ctx *context.ExecuteContext,
	cfg *config.JobConfig,
	containerID string,
) error {
	mounts := getBindMounts(ctx, cfg)
	for _, artifact := range cfg.Artifact.Globs() {
		artifactPath, err := getArtifactPath(ctx.WorkingDir, artifact, mounts)
		if err != nil {
			return err
		}
		logger.Debugf("Copying %s from container directory %s",
			artifact, artifactPath.containerDir())
		buf := new(bytes.Buffer)
		opts := docker.DownloadFromContainerOptions{
			Path:         artifactPath.containerDir(),
			OutputStream: buf,
		}
		if err := ctx.Client.DownloadFromContainer(containerID, opts); err != nil {
			return err
		}
		if err := unpack(buf, artifactPath); err != nil {
			return err
		}
	}
	return nil
}

// artifactPath stores the absolute paths of an artifact
type artifactPath struct {
	mountBind    string
	mountPath    string
	artifactGlob string
}

func newArtifactPath(mountBind, mountPath, glob string) artifactPath {
	return artifactPath{
		mountBind:    mountBind,
		mountPath:    mountPath,
		artifactGlob: glob,
	}
}

// containerDir used as the path for a container copy API call
func (p artifactPath) containerDir() string {
	return filepath.Dir(p.containerGlob())
}

// containerGlob used to match files in the archive returned by the API
func (p artifactPath) containerGlob() string {
	return rebasePath(p.artifactGlob, p.mountBind, p.mountPath)
}

// the host prefix to prepend to the archive paths
func (p artifactPath) hostBase() string {
	return p.mountBind
}

// the container path to strip from archive paths
// func (p artifactPath) containerPrefix() string { }

func getArtifactPath(
	workingDir string,
	glob string,
	mounts []config.MountConfig,
) (artifactPath, error) {
	absGlob := filepathJoinPreserveDirectorySlash(workingDir, glob)

	sortMountsByLongestBind(mounts)
	for _, mount := range mounts {
		absBindPath := filepathJoinPreserveDirectorySlash(workingDir, mount.Bind)

		if !hasPathPrefix(filepath.Dir(absGlob), absBindPath) {
			continue
		}
		return newArtifactPath(absBindPath, mount.Path, absGlob), nil
	}
	return artifactPath{}, errors.Errorf("no mount found for artifact %s", glob)
}

func sortMountsByLongestBind(mounts []config.MountConfig) {
	sort.Slice(mounts, func(i, j int) bool {
		return mounts[i].Bind >= mounts[j].Bind
	})
}

// hasPathPrefix returns true if path is under the directory prefix
func hasPathPrefix(path, prefix string) bool {
	sep := string(filepath.Separator)
	pathParts := strings.Split(filepath.Clean(path), sep)
	prefixParts := strings.Split(filepath.Clean(prefix), sep)

	if len(prefixParts) > len(pathParts) {
		return false
	}
	for index, prefixItem := range prefixParts {
		if prefixItem != pathParts[index] {
			return false
		}
	}
	return true
}

func filepathJoinPreserveDirectorySlash(elem ...string) string {
	sep := string(filepath.Separator)
	trailingSlash := ""
	if strings.HasSuffix(elem[len(elem)-1], sep) {
		trailingSlash = sep
	}
	return filepath.Join(elem...) + trailingSlash
}

func rebasePath(path, oldPrefix, newPrefix string) string {
	return filepathJoinPreserveDirectorySlash(
		newPrefix,
		strings.TrimPrefix(path, oldPrefix))
}

func unpack(source io.Reader, path artifactPath) error {
	tarReader := tar.NewReader(source)

	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}

		// TODO: remove debug
		fmt.Println("Archive Path: ", header.Name)
		match, err := filepath.Match(path.containerGlob(), header.Name)
		switch {
		case err != nil:
			return err
		case !match:
			continue
		}

		if err := createFromTar(tarReader, header, path); err != nil {
			return err
		}
	}
}

// create files and directories from tar archive entries
func createFromTar(tarReader *tar.Reader, header *tar.Header, path artifactPath) error {
	// If directory, create it if it doesn't exist

	// If file, create the parent directories if they don't exist
	// then create the file and write to it
	return nil
}
