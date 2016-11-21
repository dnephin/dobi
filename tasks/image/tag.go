package image

import (
	"fmt"

	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

// RunTag builds or pulls an image if it is out of date
func RunTag(ctx *context.ExecuteContext, t *Task, _ bool) (bool, error) {
	tag := func(tag string) error {
		return tagImage(ctx, t, tag)
	}
	if err := t.ForEachTag(ctx, tag); err != nil {
		return false, err
	}
	t.logger().Info("Tagged")
	return true, nil
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
		return fmt.Errorf("failed to add tag %q: %s", imageTag, err)
	}
	return nil
}
