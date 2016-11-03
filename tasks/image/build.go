package image

import (
	"io"
	"os"
	"strings"

	"archive/tar"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/utils/dockerignore"
	"github.com/dnephin/dobi/utils/fs"
	docker "github.com/fsouza/go-dockerclient"
)

// RunBuild builds an image if it is out of date
func RunBuild(ctx *context.ExecuteContext, t *Task, hasModifiedDeps bool) (bool, error) {
	if !hasModifiedDeps {
		stale, err := buildIsStale(ctx, t)
		switch {
		case err != nil:
			return false, err
		case !stale:
			t.logger().Info("is fresh")
			return false, nil
		}
	}
	t.logger().Debug("is stale")

	if err := buildImage(ctx, t); err != nil {
		return false, err
	}

	image, err := GetImage(ctx, t.config)
	if err != nil {
		return false, err
	}

	record := imageModifiedRecord{ImageID: image.ID}
	if err := updateImageRecord(recordPath(ctx, t.config), record); err != nil {
		t.logger().Warnf("Failed to update image record: %s", err)
	}
	t.logger().Info("Created")
	return true, nil
}

func buildIsStale(ctx *context.ExecuteContext, t *Task) (bool, error) {
	image, err := GetImage(ctx, t.config)
	switch err {
	case docker.ErrNoSuchImage:
		t.logger().Debug("Image does not exist")
		return true, nil
	case nil:
	default:
		return true, err
	}
	mtime, err := fs.LastModified(t.config.Context)
	if err != nil {
		t.logger().Warnf("Failed to get last modified time of context.")
		return true, err
	}

	record, err := getImageRecord(recordPath(ctx, t.config))
	if err != nil {
		t.logger().Warnf("Failed to get image record: %s", err)
		if image.Created.Before(mtime) {
			t.logger().Debug("Image older than context")
			return true, nil
		}
		return false, nil
	}

	if image.ID != record.ImageID || record.Info.ModTime().Before(mtime) {
		t.logger().Debug("Image record older than context")
		return true, nil
	}
	return false, nil
}

func buildImage(ctx *context.ExecuteContext, t *Task) error {
	var err error
	if t.hasSteps() {
		err = t.buildImageFromTarBall(ctx)
	} else {
		err = t.buildImageFromDockerfile(ctx)
	}
	if err != nil {
		return err
	}
	image, err := GetImage(ctx, t.config)
	if err != nil {
		return err
	}
	record := imageModifiedRecord{ImageID: image.ID}
	return updateImageRecord(recordPath(ctx, t.config), record)
}

func buildArgs(args map[string]string) []docker.BuildArg {
	out := []docker.BuildArg{}
	for key, value := range args {
		out = append(out, docker.BuildArg{Name: key, Value: value})
	}
	return out
}

func (t *Task) buildImageFromDockerfile(ctx *context.ExecuteContext) error {
	err := Stream(os.Stdout, func(out io.Writer) error {
		return ctx.Client.BuildImage(docker.BuildImageOptions{
			Name:           GetImageName(ctx, t.config),
			Dockerfile:     t.config.Dockerfile,
			BuildArgs:      buildArgs(t.config.Args),
			Pull:           t.config.PullBaseImageOnBuild,
			RmTmpContainer: true,
			ContextDir:     t.config.Context,
			OutputStream:   out,
			RawJSONStream:  true,
			SuppressOutput: ctx.Settings.Quiet,
			AuthConfigs:    ctx.GetAuthConfigs(),
		})
	})
	return err
}

func (t *Task) buildImageFromTarBall(ctx *context.ExecuteContext) error {
	inputbuf, err := t.writeTarball()
	if err != nil {
		return err
	}
	err = Stream(os.Stdout, func(out io.Writer) error {
		return ctx.Client.BuildImage(docker.BuildImageOptions{
			Name:           GetImageName(ctx, t.config),
			BuildArgs:      buildArgs(t.config.Args),
			Pull:           t.config.PullBaseImageOnBuild,
			RmTmpContainer: true,
			InputStream:    inputbuf,
			OutputStream:   out,
			RawJSONStream:  true,
			SuppressOutput: ctx.Quiet,
		})
	})
	return err
}

func (t *Task) writeTarball() (*bytes.Buffer, error) {
	inputbuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	defer tr.Close()
	err := t.writeDockerfiletoTarBall(tr)
	if err != nil {
		return inputbuf, err
	}
	err = t.writeFilesToTarBall(tr)
	if err != nil {
		return inputbuf, err
	}
	return inputbuf, nil
}

func (t *Task) getTarContext() ([]string, error) {
	type Slice []string
	type Result struct {
		Slice
		error
	}
	ignored := make(chan *Result)
	ctx := make(chan *Result)
	go func() {
		result := new(Result)
		result.Slice, result.error = t.scanContext()
		ctx <- result
	}()
	go func() {
		result := new(Result)
		result.Slice, result.error = t.scanIgnored()
		ignored <- result
	}()
	ctxFiles := <-ctx
	if ctxFiles.error != nil {
		return []string{}, ctxFiles.error
	}
	ignoredFiles := <-ignored
	if ignoredFiles.error != nil {
		return []string{}, ignoredFiles.error
	}
	return dockerignore.Difference(ctxFiles.Slice, ignoredFiles.Slice), nil
}
func (t *Task) writeDockerfiletoTarBall(tr *tar.Writer) error {
	rightNow := time.Now()
	header := &tar.Header{Name: "Dockerfile",
		Size:       t.getStepsSize(),
		ModTime:    rightNow,
		AccessTime: rightNow,
		ChangeTime: rightNow,
	}
	err := tr.WriteHeader(header)
	if err != nil {
		return err
	}
	for _, val := range t.config.Steps {
		// its not a good idea to catch this error?
		tr.Write([]byte(val + "\n"))
	}
	return nil
}

func (t *Task) writeFilesToTarBall(tr *tar.Writer) error {
	paths, err := t.getTarContext()
	if err != nil {
		return err
	}
	for _, file := range paths {
		t.logger().Debugf("is writing %s to tarball", strings.TrimPrefix(file, filepath.Base(t.config.Context)+"/"))
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		byt, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		rightNow := time.Now()
		header := &tar.Header{Name: strings.TrimPrefix(file, filepath.Base(t.config.Context)+"/"),
			Size:       int64(len(byt)),
			ModTime:    rightNow,
			Mode:       int64(fileInfo.Mode()),
			AccessTime: rightNow,
			ChangeTime: rightNow,
		}
		err = tr.WriteHeader(header)
		if err != nil {
			return err
		}

		_, err = tr.Write(byt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Task) getStepsSize() int64 {
	var size int64
	for _, val := range t.config.Steps {
		size = size + int64(len([]byte(val+"\n")))
	}
	return size
}

func (t *Task) scanIgnored() ([]string, error) {
	allIgnored, err := dockerignore.ReadAll()
	if err != nil {
		return []string{}, err
	}
	var resolvedignores []string
	for _, val := range allIgnored {
		resolvedignores, err = scanRoot2Slice(val, resolvedignores)
		if err != nil {
			return resolvedignores, err
		}
	}
	return resolvedignores, nil
}

func (t *Task) scanContext() ([]string, error) {
	return scanRoot2Slice(t.config.Context, []string{})
}

func scanRoot2Slice(root string, placeholder []string) ([]string, error) {
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			placeholder = append(placeholder, path)
		}
		return nil
	})
	return placeholder, err
}

func (t *Task) hasSteps() bool {
	if len(t.config.Steps) != 0 {
		return true
	}
	return false
}
