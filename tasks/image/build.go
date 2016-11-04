package image

import (
	"io"
	"os"

	"github.com/dnephin/dobi/tasks/context"
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

	dockerfile, err := dobiOrDocker(t)
	if err != nil {
		return err
	}
	if isDobi(t){
		defer os.Remove(dockerfile)
	}

	if err := Stream(os.Stdout, func(out io.Writer) error {
		return ctx.Client.BuildImage(docker.BuildImageOptions{
			Name:           GetImageName(ctx, t.config),
			Dockerfile:     dockerfile,
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

func isDobi(t *Task) (bool){
	if len(t.config.Dobifile) != 0 && t.config.Dockerfile == "Dockerfile" {
		return true
	}
	return false
}

func dobiOrDocker(t  *Task) (string, error) {
	if isDobi(t) {
		return parseDobifile(t)
	}
	return t.config.Dockerfile, nil
}

func parseDobifile(t *Task) (string, error) {
	var dobifile string
	for _, val := range t.config.Dobifile {
		for key, value := range val {
			dobifile = dobifile + key + " " + value + "\n"
		}
	}
	return createDockerfilefromString(".dobi/Dockerfile." + t.name.Resource(), dobifile)
}

func createDockerfilefromString(path, str string) (string, error) {
	_, err := os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}
	tempfile, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer tempfile.Close()
	_, err = tempfile.WriteString(str)
	if err != nil {
		return "", err
	}
	err = tempfile.Sync()
	if err != nil {
		return "", err
	}
	return path, nil
}

