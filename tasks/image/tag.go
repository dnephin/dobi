package image

import (
	"fmt"

	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

// RunTag builds or pulls an image if it is out of date
func RunTag(ctx *context.ExecuteContext, t *Task) error {
	tag := func(tag string) error {
		return tagImage(ctx, t, tag)
	}
	if err := t.ForEachTag(ctx, tag); err != nil {
		return err
	}
	ctx.SetModified(t.name)
	t.logger().Info("Tagged")
	return nil
}

func tagImage(ctx *context.ExecuteContext, t *Task, imageTag string) error {
	if imageTag == GetImageName(ctx, t.config) {
		return nil
	}

	repo, tag := docker.ParseRepositoryTag(imageTag)
	err := ctx.Client.TagImage(GetImageName(ctx, t.config), docker.TagImageOptions{
		Repo:  repo,
		Tag:   tag,
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("Failed to add tag %q: %s", imageTag, err)
	}
	return nil
}
