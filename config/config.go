package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/configtf"
	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/common"
)

// Config is a data object for a full config file
type Config struct {
	Meta       *MetaConfig
	Resources  map[string]Resource
	WorkingDir string
	Collection *ResourceCollection
}

// NewConfig returns a new Config object
func NewConfig() *Config {
	return &Config{
		Resources:  make(map[string]Resource),
		Meta:       &MetaConfig{},
		Collection: newResourceCollection(),
	}
}

func (c *Config) add(name string, resource Resource) error {
	if c.contains(name) {
		return fmt.Errorf("duplicate resource name %q", name)
	}
	c.Resources[name] = resource
	c.Collection.add(name, resource)
	return nil
}

func (c *Config) contains(name string) bool {
	_, exists := c.Resources[name]
	return exists
}

// Sorted returns the list of resource names in alphabetical sort order
func (c *Config) Sorted() []string {
	names := []string{}
	for name := range c.Resources {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Load a configuration from a filename
func Load(filename string) (*Config, error) {
	fmtError := func(err error) error {
		return fmt.Errorf("Failed to load config from %q: %s", filename, err)
	}

	config, err := loadConfig(filename)
	if err != nil {
		return nil, fmtError(err)
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmtError(err)
	}
	config.WorkingDir = filepath.Dir(absPath)

	if err = validate(config); err != nil {
		return nil, fmtError(err)
	}
	return config, nil
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := LoadFromBytes(data)
	if err != nil {
		return nil, err
	}
	logging.Log.WithFields(log.Fields{"filename": filename}).Debug("Configuration loaded")
	return config, nil
}

// validate validates all the resources in the config
func validate(config *Config) error {
	for name, resource := range config.Resources {
		path := pth.NewPath(name)

		if err := configtf.ValidateFields(path, resource); err != nil {
			return err
		}
		if err := ValidateResourcesExist(path, config, resource.Dependencies()); err != nil {
			return err
		}
		if err := resource.Validate(path, config); err != nil {
			return err
		}
	}
	config.Meta.Validate(config)
	return nil
}

// ValidateResourcesExist checks that the list of resources is defined in the
// config and returns an error if a resources is not defined.
func ValidateResourcesExist(path pth.Path, c *Config, names []string) error {
	missing := []string{}
	for _, name := range names {
		resource := common.ParseTaskName(name).Resource()
		if _, ok := c.Resources[resource]; !ok {
			missing = append(missing, resource)
		}
	}
	if len(missing) != 0 {
		return pth.Errorf(path, "missing dependencies: %s", strings.Join(missing, ", "))
	}
	return nil
}
