package image

import (
	"io"
	"os"

	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

// RunPush builds or pulls an image if it is out of date
func RunPush(ctx *context.ExecuteContext, t *Task) error {
	// TODO: triger a RunBuild action if it's not already in the modified set?

	t.logger().Info("Pushing")
	pushTag := func(tag string) error {
		if err := tagImage(ctx, t, tag); err != nil {
			return err
		}
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

func tagImage(ctx *context.ExecuteContext, t *Task, tag string) error {
	if tag == GetImageName(ctx, t.config) {
		return nil
	}

	return ctx.Client.TagImage(GetImageName(ctx, t.config), docker.TagImageOptions{
		// TODO: test this
		Repo:  tag,
		Force: true,
	})
}
