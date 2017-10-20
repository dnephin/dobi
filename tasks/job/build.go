package job

import (
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
)

func (t *Task) runWithBuildAndCopy(ctx *context.ExecuteContext) error {
	imageID, err := buildImageWithMounts(ctx)
	if err != nil {
		return err
	}

	name := containerName(ctx, t.name.Resource())
	options := t.createOptions(ctx, name, imageID)
	defer removeContainerWithLogging(t.logger(), ctx.Client, name)

	if err := t.runContainer(ctx, options); err != nil {
		return err
	}
	return copyFilesToHost(ctx, name)
}

func buildImageWithMounts(ctx *context.ExecuteContext) (string, error) {
	imageName := image.GetImageName(ctx, ctx.Resources.Image(t.config.Use))
	return "", nil
}

func copyFilesToHost(ctx *context.ExecuteContext, containerID string) error {
	return nil
}
