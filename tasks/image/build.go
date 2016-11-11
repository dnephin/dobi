package image

import (
	"io"
	"os"

	"archive/tar"
	"bytes"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/utils/dockerignore"
	"github.com/dnephin/dobi/utils/fs"
	docker "github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"log"
	"path/filepath"
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
	switch t.hasSteps() {
	case true:
		err = t.buildImageFromTarBall(ctx)
	default:
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

func (t *Task) hasSteps() bool {
	if len(t.config.Steps) != 0 {
		return true
	}
	return false
}

func (t *Task) writeTarball() (*bytes.Buffer, error) {
	inputbuf := bytes.NewBuffer(nil)
	err := t.writeDockerfiletoTarBall(inputbuf)
	if err != nil {
		return inputbuf, err
	}
	allContext, err := t.scanContext()
	if err != nil {
		return inputbuf, err
	}
	ignored, err := t.scanIgnored()
	if err != nil {
		return inputbuf, err
	}
	diff := dockerignore.Difference(allContext, ignored)
	err = t.writeFilesToTarBall(diff, inputbuf)
	if err != nil {
		return inputbuf, err
	}
	return inputbuf, nil
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
		size = size + int64(len([]byte(val+"\n")))
	}
	return size
}

func (t *Task) scanIgnored() ([]string, error) {
	bytesrece, err := ioutil.ReadFile(".dockerignore")
	if err != nil {
		return []string{}, nil
	}
	r := bytes.NewReader(bytesrece)
	allIgnored, err := dockerignore.ReadAll(r)
	if err != nil {
		return []string{}, err
	}
	var resolvedignores []string
	for _, val := range allIgnored {
		err := filepath.Walk(val, func(path string, f os.FileInfo, err error) error {
			resolvedignores = append(resolvedignores, path)
			return nil
		})
		if err != nil {
			return []string{}, err
		}
	}
	return resolvedignores, nil
}

func (t *Task) scanContext() ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(t.config.Context, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return fileList, nil
}

func (t *Task) writeFilesToTarBall(files []string, inputbuf *bytes.Buffer) error {
	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			log.Println(file)
			log.Println(fileInfo.Name())
			tr := tar.NewWriter(inputbuf)
			rightNow := time.Now()
			header := &tar.Header{Name: fileInfo.Name(),
				Size:       fileInfo.Size(),
				ModTime:    rightNow,
				AccessTime: rightNow,
				ChangeTime: rightNow,
			}
			err := tr.WriteHeader(header)
			if err != nil {
				return err
			}
			byt, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			_, err = tr.Write(byt)
			if err != nil {
				return err
			}
			tr.Close()
		}
	}

	return nil
}

func (t *Task) writeDockerfiletoTarBall(inputbuf *bytes.Buffer) error {
	tr := tar.NewWriter(inputbuf)
	rightNow := time.Now()
	header := &tar.Header{Name: "Dockerfile",
		Size:       t.getContentSize(),
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
	err = tr.Close()
	if err != nil {
		return err
	}
	return nil
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
	if err != nil {
		return err
	}
	return nil
}
