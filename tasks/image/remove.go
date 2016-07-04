package image

import (
	"github.com/dnephin/dobi/tasks/context"
)

// RunRemove builds or pulls an image if it is out of date
func RunRemove(ctx *context.ExecuteContext, t *Task) error {
	t.logger().Info("Removing")

	removeTag := func(tag string) error {
		if err := ctx.Client.RemoveImage(tag); err != nil {
			t.logger().Warnf("failed to remove %q: %s", tag, err)
		}
		return nil
	}

	if err := t.ForEachTag(ctx, removeTag); err != nil {
		return err
	}
	ctx.SetModified(t.name)
	t.logger().Info("Removed")
	return nil
}
