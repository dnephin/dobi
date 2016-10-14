package image

import (
	"fmt"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	defaultRepo = "https://index.docker.io/v1/"
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
		return conf.Tags[0]
	}
	return ctx.Env.Unique()
}

func parseAuthRepo(image string) (string, error) {
	// This is the approximate logic from
	// github.com/docker/docker/reference.splitHostname(). That package is
	// conflicting with other dependencies, so it can't be imported at this time.
	parts := strings.SplitN(image, "/", 3)
	switch len(parts) {
	case 1, 2:
		logging.Log.Debugf("Using default registry %q", defaultRepo)
		return defaultRepo, nil
	default:
		logging.Log.Debugf("Using registry %q", parts[0])
		return parts[0], nil
	}
}
