package image

import (
	"io"
	"os"

	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

// RunPush builds or pulls an image if it is out of date
func RunPush(ctx *context.ExecuteContext, t *Task) error {
	t.logger().Info("Pushing")
	pushTag := func(tag string) error {
		return pushImage(ctx, t, tag)
	}
	if err := t.ForEachTag(ctx, pushTag); err != nil {
		return err
	}
	ctx.SetModified(t.name)
	t.logger().Info("Pushed")
	return nil
}

func pushImage(ctx *context.ExecuteContext, t *Task, tag string) error {
	repo, err := parseRepo(t.config.Image)
	if err != nil {
		return err
	}

	return Stream(os.Stdout, func(out io.Writer) error {
		return ctx.Client.PushImage(docker.PushImageOptions{
			Name:          tag,
			OutputStream:  out,
			RawJSONStream: true,
			// TODO: timeout
		}, ctx.GetAuthConfig(repo))
	})
}
