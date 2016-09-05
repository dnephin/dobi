package image

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/utils/fs"
	docker "github.com/fsouza/go-dockerclient"
	yaml "gopkg.in/yaml.v2"
)

const (
	imageRecordDir = ".dobi/images"
	all            = -1
)

// RunBuild builds or pulls an image if it is out of date
func RunBuild(ctx *context.ExecuteContext, t *Task) error {
	stale, err := buildIsStale(ctx, t)
	if !stale || err != nil {
		t.logger().Info("is fresh")
		return err
	}
	t.logger().Debug("is stale")

	if err := buildImage(ctx, t); err != nil {
		return err
	}

	image, err := GetImage(ctx, t.config)
	if err != nil {
		return err
	}
	if err := updateImageRecord(recordPath(ctx, t.config), image.ID); err != nil {
		t.logger().Warnf("Failed to update image record: %s", err)
	}
	ctx.SetModified(t.name)
	t.logger().Info("Created")
	return nil
}

func buildIsStale(ctx *context.ExecuteContext, t *Task) (bool, error) {
	if ctx.IsModified(t.config.Dependencies()...) {
		return true, nil
	}

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

type imageModifiedRecord struct {
	ImageID string
	Info    os.FileInfo `yaml:",omitempty"`
}

func updateImageRecord(path string, imageID string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	record := imageModifiedRecord{ImageID: imageID}
	bytes, err := yaml.Marshal(record)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

// TODO: verify error message are sufficient
func getImageRecord(filepath string) (*imageModifiedRecord, error) {
	record := &imageModifiedRecord{}
	var err error

	record.Info, err = os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	recordBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return record, yaml.Unmarshal(recordBytes, record)
}

func recordPath(ctx *context.ExecuteContext, conf *config.ImageConfig) string {
	imageName := strings.Replace(GetImageName(ctx, conf), "/", " ", all)
	return filepath.Join(ctx.WorkingDir, imageRecordDir, imageName)
}

func buildImage(ctx *context.ExecuteContext, t *Task) error {
	if err := Stream(os.Stdout, func(out io.Writer) error {
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
	}); err != nil {
		return err
	}

	image, err := GetImage(ctx, t.config)
	if err != nil {
		return err
	}
	return updateImageRecord(recordPath(ctx, t.config), image.ID)
}

func buildArgs(args map[string]string) []docker.BuildArg {
	out := []docker.BuildArg{}
	for key, value := range args {
		out = append(out, docker.BuildArg{Name: key, Value: value})
	}
	return out
}
