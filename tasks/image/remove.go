package image

import (
	"github.com/dnephin/dobi/tasks/context"
)

// RunRemove removes an image
func RunRemove(ctx *context.ExecuteContext, t *Task, _ bool) (bool, error) {
	removeTag := func(tag string) error {
		if err := ctx.Client.RemoveImage(tag); err != nil {
			t.logger().Warnf("failed to remove %q: %s", tag, err)
		}
		return nil
	}

	if err := t.ForEachTag(ctx, removeTag); err != nil {
		return false, err
	}

	// Clear the image record so the .dobi state does not break for "pull once" images
	if err := updateImageRecord(recordPath(ctx, t.config), imageModifiedRecord{}); err != nil {
		t.logger().Warnf("Failed to clear image record: %s", err)
	}

	t.logger().Info("Removed")
	return true, nil
}
