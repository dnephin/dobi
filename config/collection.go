package config

// ResourceCollection holds resource configs that are used by other resources
type ResourceCollection struct {
	mounts map[string]*MountConfig
	images map[string]*ImageConfig
}

func (c *ResourceCollection) add(name string, resource Resource) {
	switch resource := resource.(type) {
	case *MountConfig:
		c.mounts[name] = resource
	case *ImageConfig:
		c.images[name] = resource
	}
}

// Mount returns a MountConfig by name
func (c *ResourceCollection) Mount(name string) *MountConfig {
	return c.mounts[name]
}

// Image returns an ImageConfig by name
func (c *ResourceCollection) Image(name string) *ImageConfig {
	return c.images[name]
}

type eachMountFunc func(name string, vol *MountConfig)

// EachMount iterates all the mounts in names and calls f for each
func (c *ResourceCollection) EachMount(names []string, f eachMountFunc) {
	for _, name := range names {
		mount, _ := c.mounts[name]
		f(name, mount)
	}
}

func newResourceCollection() *ResourceCollection {
	return &ResourceCollection{
		mounts: make(map[string]*MountConfig),
		images: make(map[string]*ImageConfig),
	}
}
