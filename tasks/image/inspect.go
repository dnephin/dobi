package image

import (
	"fmt"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

// GetImage returns the image created by an image config
func GetImage(ctx *context.ExecuteContext, conf *config.ImageConfig) (*docker.Image, error) {
	return ctx.Client.InspectImage(GetImageName(ctx, conf))
}

// GetImageName returns the image name for an image config
func GetImageName(ctx *context.ExecuteContext, conf *config.ImageConfig) string {
	return fmt.Sprintf("%s:%s", conf.Image, GetCanonicalTag(ctx, conf))
}

// GetCanonicalTag returns the canonical tag for an image config
func GetCanonicalTag(ctx *context.ExecuteContext, conf *config.ImageConfig) string {
	if len(conf.Tags) > 0 {
		return ctx.Env.GetVar(conf.Tags[0])
	}
	return ctx.Env.Unique()
}
