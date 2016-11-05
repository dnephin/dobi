package image

import (
	"io"
	"os"

	"archive/tar"
	"bytes"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/utils/fs"
	docker "github.com/fsouza/go-dockerclient"
	"time"
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
	switch t.hasContent() {
	case true:
		err = t.runContainerFromTarBall(ctx)
	default:
		err = t.runContainerFromDockerfile(ctx)
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

func (t *Task) hasContent() bool {
	if len(t.config.Steps) != 0 {
		return true
	}
	return false
}

func (t *Task) writeTarball() (*bytes.Buffer, error) {
	inputbuf := bytes.NewBuffer(nil)
	rightNow := time.Now()
	tr := tar.NewWriter(inputbuf)
	header := &tar.Header{Name: "Dockerfile",
		Size:       t.getContentSize(),
		ModTime:    rightNow,
		AccessTime: rightNow,
		ChangeTime: rightNow,
	}
	err := tr.WriteHeader(header)
	if err != nil {
		return inputbuf, err
	}
	for _, val := range t.config.Steps {
		for key, value := range val {
			// its not a good idea to catch this error?
			tr.Write([]byte(key + " " + value + "\n"))
		}
	}
	tr.Close()
	return inputbuf, nil
}

func (t *Task) runContainerFromDockerfile(ctx *context.ExecuteContext) error {
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
			SuppressOutput: ctx.Quiet,
		})
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) getContentSize() int64 {
	var size int64
	for _, val := range t.config.Steps {
		for key, value := range val {
			size = size + int64(len([]byte(key+" "+value+"\n")))
		}
	}
	return size
}

func (t *Task) runContainerFromTarBall(ctx *context.ExecuteContext) error {
	err := t.replaceCustomSteps()
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
	return nil
}
