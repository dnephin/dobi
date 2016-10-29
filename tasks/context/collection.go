package context

import (
	"github.com/dnephin/dobi/config"
)

// ResourceCollection holds resource configs that are used by other resources
type ResourceCollection struct {
	mounts map[string]*config.MountConfig
	images map[string]*config.ImageConfig
}

// Add a resource to the collection
func (c *ResourceCollection) Add(name string, resource config.Resource) {
	switch resource := resource.(type) {
	case *config.MountConfig:
		c.mounts[name] = resource
	case *config.ImageConfig:
		c.images[name] = resource
	}
}

// Mount returns a config.MountConfig by name
func (c *ResourceCollection) Mount(name string) *config.MountConfig {
	return c.mounts[name]
}

// Image returns an config.ImageConfig by name
func (c *ResourceCollection) Image(name string) *config.ImageConfig {
	return c.images[name]
}

type eachMountFunc func(name string, vol *config.MountConfig)

// EachMount iterates all the mounts in names and calls f for each
func (c *ResourceCollection) EachMount(names []string, f eachMountFunc) {
	for _, name := range names {
		mount, _ := c.mounts[name]
		f(name, mount)
	}
}

func newResourceCollection() *ResourceCollection {
	return &ResourceCollection{
		mounts: make(map[string]*config.MountConfig),
		images: make(map[string]*config.ImageConfig),
	}
}
