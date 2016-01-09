package config

// ImageConfig ia a data object for image resource
type ImageConfig struct {
	Image      string
	Dockerfile string
	Context    string
	Args       map[string]string
	Pull       bool
	Depends    []string
}

// Dependencies returns the list of implicit and explicit dependencies
func (c *ImageConfig) Dependencies() []string {
	return c.Depends
}

// Validate checks that all fields have acceptable values
func (c *ImageConfig) Validate() error {
	// TODO: better way to generate consistent config errors
	// TODO: check context directory exists
	// TODO: check dockerfile exists
	// TODO: validate required fields are set
	// TODO: validate no tag on image name
	return nil
}

// NewImageConfig creates an ImageConfig with default values
func NewImageConfig() *ImageConfig {
	return &ImageConfig{
		Context:    ".",
		Dockerfile: "Dockerfile",
	}
}

// CommandConfig is a data object for a command resource
type CommandConfig struct {
	Use        string
	Artifact   string
	Command    string
	Volumes    []string
	Privileged bool
	Depends    []string
}

// TODO: support interactive/tty

// Dependencies returns the list of implicit and explicit dependencies
func (c *CommandConfig) Dependencies() []string {
	return append([]string{c.Use}, append(c.Volumes, c.Depends...)...)
}

// Validate checks that all fields have acceptable values
func (c *CommandConfig) Validate() error {
	// TODO: validate required fields are set
	return nil
}

// VolumeConfig is a data object for a volume resource
type VolumeConfig struct {
	Path  string
	Mount string
	Mode  string
}

// Dependencies returns an empty list, Volume resources have no dependencies
func (c *VolumeConfig) Dependencies() []string {
	return []string{}
}

// Validate checks that all fields have acceptable values
func (c *VolumeConfig) Validate() error {
	// TODO: validate required fields are set
	return nil
}

// NewVolumeConfig creates a VolumeConfig with default values
func NewVolumeConfig() *VolumeConfig {
	return &VolumeConfig{
		Mode: "rw",
	}
}
