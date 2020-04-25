package job

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/pkg/archive"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	runErr := t.runContainer(ctx, options)
	copyErr := copyFilesToHost(t.logger(), ctx, t.config, name)
	if runErr != nil {
		return runErr
	}
	return copyErr
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

	sortMountsByShortestPath(mounts)
	for _, mount := range mounts {
		buf.WriteString(fmt.Sprintf("COPY %s %s\n", mount.Bind, mount.Path))
	}
	return buf
}

func sortMountsByShortestPath(mounts []config.MountConfig) {
	sort.Slice(mounts, func(i, j int) bool {
		return mounts[i].Path < mounts[j].Path
	})
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
	return filepathDirWithDirectorySlash(p.containerGlob())
}

// containerGlob used to match files in the archive returned by the API
func (p artifactPath) containerGlob() string {
	return rebasePath(p.artifactGlob, p.mountBind, p.mountPath)
}

// the host prefix to prepend to the archive paths
func (p artifactPath) hostPath(path string) string {
	return rebasePath(path, p.mountPath, p.mountBind)
}

// pathFromArchive strips the archive directory from the path and returns the
// absolute path to the file in a container
func (p artifactPath) pathFromArchive(path string) string {
	parts := strings.SplitN(path, string(filepath.Separator), 2)
	if len(parts) == 1 || parts[1] == "" {
		return p.containerDir()
	}
	return filepathJoinPreserveDirectorySlash(p.containerDir(), parts[1])
}

func getArtifactPath(
	workingDir string,
	glob string,
	mounts []config.MountConfig,
) (artifactPath, error) {
	absGlob := filepathJoinPreserveDirectorySlash(workingDir, glob)

	sortMountsByLongestBind(mounts)
	for _, mount := range mounts {
		absBindPath := filepathJoinPreserveDirectorySlash(workingDir, mount.Bind)

		if mount.File && hasPathPrefix(absGlob, absBindPath) {
			return newArtifactPath(absBindPath, mount.Path, absGlob), nil
		}

		if hasPathPrefix(filepathDirWithDirectorySlash(absGlob), absBindPath) {
			return newArtifactPath(absBindPath, mount.Path, absGlob), nil
		}
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
	trailingSlash := ""
	if endsWithSlash(elem[len(elem)-1]) {
		trailingSlash = string(filepath.Separator)
	}
	return filepath.Join(elem...) + trailingSlash
}

func filepathDirWithDirectorySlash(path string) string {
	return filepath.Dir(path) + string(filepath.Separator)
}

func rebasePath(path, oldPrefix, newPrefix string) string {
	relativePath := strings.TrimPrefix(path, oldPrefix)
	if relativePath == "" {
		return newPrefix
	}
	return filepathJoinPreserveDirectorySlash(newPrefix, relativePath)
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

		containerPath := path.pathFromArchive(header.Name)
		match, err := fileMatchesGlob(containerPath, path.containerGlob())
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

func fileMatchesGlob(path string, glob string) (bool, error) {
	// Directory glob should match entire tree
	if endsWithSlash(glob) && strings.HasPrefix(path, glob) {
		return true, nil
	}

	return filepath.Match(glob, path)
}

func endsWithSlash(path string) bool {
	return strings.HasSuffix(path, string(filepath.Separator))
}

// create files and directories from tar archive entries
func createFromTar(tarReader io.Reader, header *tar.Header, path artifactPath) error {
	hostPath := path.hostPath(path.pathFromArchive(header.Name))
	fileMode := header.FileInfo().Mode()

	switch header.Typeflag {
	case tar.TypeDir:
		logging.Log.Debugf("Creating dir %s", hostPath)
		return os.MkdirAll(hostPath, fileMode)

	case tar.TypeReg, tar.TypeRegA:
		logging.Log.Debugf("Creating file %s", hostPath)
		if err := os.MkdirAll(filepath.Dir(hostPath), 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(hostPath, os.O_RDWR|os.O_CREATE, fileMode)
		if err != nil {
			return err
		}
		_, err = io.Copy(file, tarReader)
		return err

	case tar.TypeSymlink:
		logging.Log.Debugf("Creating symlink %s", hostPath)
		if err := os.MkdirAll(filepath.Dir(hostPath), 0755); err != nil {
			return err
		}

		return os.Symlink(header.Linkname, hostPath)

	default:
		logging.Log.Warnf("Unhandled file type from archive %s: %s",
			string(header.Typeflag),
			header.Name)
	}

	return nil
}
