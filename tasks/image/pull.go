package image

import (
	"io"
	"os"
	"time"

	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

// RunPull builds or pulls an image if it is out of date
func RunPull(ctx *context.ExecuteContext, t *Task, _ bool) (bool, error) {
	record, err := getImageRecord(recordPath(ctx, t.config))
	switch {
	case !t.config.Pull.Required(record.LastPull):
		t.logger().Debugf("Pull not required")
		return false, nil
	case err != nil:
		t.logger().Warnf("Failed to get image record: %s", err)
	}

	pullTag := func(tag string) error {
		return pullImage(ctx, t, tag)
	}
	if err := t.ForEachTag(ctx, pullTag); err != nil {
		return false, err
	}

	record = imageModifiedRecord{LastPull: now()}
	if err := updateImageRecord(recordPath(ctx, t.config), record); err != nil {
		t.logger().Warnf("Failed to update image record: %s", err)
	}

	t.logger().Info("Pulled")
	return true, nil
}

func now() *time.Time {
	now := time.Now()
	return &now
}

func pullImage(ctx *context.ExecuteContext, t *Task, imageTag string) error {
	registry := parseAuthRepo(t.config.Image)
	repo, tag := docker.ParseRepositoryTag(imageTag)
	return Stream(os.Stdout, func(out io.Writer) error {
		return ctx.Client.PullImage(docker.PullImageOptions{
			Repository:    repo,
			Tag:           tag,
			OutputStream:  out,
			RawJSONStream: true,
			// TODO: timeout
		}, ctx.GetAuthConfig(registry))
	})
}
