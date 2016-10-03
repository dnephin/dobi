package image

import (
	"fmt"
	"net/url"
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

func parseRepo(image string) (string, error) {
	repo, _ := docker.ParseRepositoryTag(image)
	if !strings.HasPrefix(repo, "http") {
		repo = "https://" + repo
	}

	addr, err := url.Parse(repo)
	if err != nil {
		return "", fmt.Errorf("Failed to parse repo name from %q: %s", repo, err)
	}

	if addr.Host == "" {
		logging.Log.Debugf("Using default registry %q", defaultRepo)
		return defaultRepo, nil
	}

	// TODO: what about v2?
	addr.Path = "v1"
	logging.Log.Debugf("Using registry %q", addr)
	return addr.String(), nil
}
