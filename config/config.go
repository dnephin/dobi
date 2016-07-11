package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/execenv"
	"github.com/dnephin/dobi/logging"
)

// Resource is an interface for each configurable type
type Resource interface {
	Dependencies() []string
	Validate(Path, *Config) *PathError
	Resolve(*execenv.ExecEnv) (Resource, error)
}

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

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmtError(err)
	}
	config, err := LoadFromBytes(data)
	if err != nil {
		return nil, fmtError(err)
	}
	logging.Log.WithFields(log.Fields{"filename": filename}).Debug("Configuration loaded")

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

// validate validates all the resources in the config
func validate(config *Config) error {
	for name, resource := range config.Resources {
		path := NewPath(name)

		if err := ValidateFields(path, resource); err != nil {
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
func ValidateResourcesExist(path Path, c *Config, names []string) error {
	missing := []string{}
	for _, name := range names {
		if _, ok := c.Resources[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) != 0 {
		return PathErrorf(path, "missing dependencies: %s", strings.Join(missing, ", "))
	}
	return nil
}

// ValidateFields runs validations as defined by struct tags
func ValidateFields(path Path, resource interface{}) error {
	ptrValue := reflect.ValueOf(resource)
	value := ptrValue.Elem()

	if kind := value.Kind(); kind != reflect.Struct {
		return fmt.Errorf("invalid target type %s, must be a Struct", kind)
	}

	for i := 0; i < value.Type().NumField(); i++ {
		field := value.Type().Field(i)
		tag := field.Tag.Get("config")
		if tag == "" {
			continue
		}
		fieldPath := path.add(titleCaseToDash(field.Name))
		if err := validateField(fieldPath, ptrValue, field, tag); err != nil {
			return err
		}
	}
	return nil
}

func validateField(path Path, structValue reflect.Value, field reflect.StructField, tag string) error {
	value := structValue.Elem().FieldByName(field.Name)
	for _, item := range strings.Split(tag, ",") {
		switch item {
		case "required":
			// TODO: better way to do this?
			if reflect.DeepEqual(value.Interface(), reflect.Zero(field.Type).Interface()) {
				return PathErrorf(path, "a value is required")
			}
		case "validate":
			if err := runValidationFunc(path, structValue, field.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func runValidationFunc(path Path, structValue reflect.Value, field string) error {
	methodName := "Validate" + field
	methodValue, err := getMethodFromStruct(structValue, methodName)
	if err != nil {
		return err
	}

	switch validationFunc := methodValue.Interface().(type) {
	case func() error:
		if err := validationFunc(); err != nil {
			return PathErrorf(path, "failed validation: %s", err)
		}
		return nil
	default:
		return fmt.Errorf("%s.%s must be of type \"func() error\" not %T",
			structValue.Elem().Type(), methodName, validationFunc)
	}
}

func getMethodFromStruct(structValue reflect.Value, methodName string) (reflect.Value, error) {
	// First look for method with non-pointer receiver
	methodValue := structValue.Elem().MethodByName(methodName)
	if methodValue.IsValid() {
		return methodValue, nil
	}

	// Second look for method with pointer receiver
	methodValue = structValue.MethodByName(methodName)
	if methodValue.IsValid() {
		return methodValue, nil
	}

	return reflect.Value{}, fmt.Errorf("%s is missing validation function %q",
		structValue.Type(), methodName)
}
