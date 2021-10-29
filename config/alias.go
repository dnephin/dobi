package config

import (
	"fmt"
	"strings"

	"github.com/dnephin/configtf"
	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/tasks/task"
)

// AliasConfig An **alias** resource is a list of other tasks which will be run
// in the order they are listed.
// example: An alias that runs three other tasks:
//
// .. code-block:: yaml
//
//     alias=test
//         tasks: [test-unit, test-integration, test-acceptance]
//
// name: alias
type AliasConfig struct {
	// Tasks The list of tasks
	// type: list of tasks
	Tasks []string `config:"required"`
	Annotations
}

// Dependencies returns the list of tasks
func (c *AliasConfig) Dependencies() ([]task.Name, error) {
	return task.ParseNames(c.Tasks)
}

// Validate the resource
func (c *AliasConfig) Validate(path pth.Path, config *Config) *pth.Error {
	return nil
}

func (c *AliasConfig) String() string {
	return fmt.Sprintf("Run tasks: %v", strings.Join(c.Tasks, ", "))
}

// Resolve resolves variables in the resource
func (c *AliasConfig) Resolve(_ Resolver) (Resource, error) {
	copy := *c
	return &copy, nil
}

func aliasFromConfig(name string, values map[string]interface{}) (Resource, error) {
	alias := &AliasConfig{}
	return alias, configtf.Transform(name, values, alias)
}

func init() {
	RegisterResource("alias", aliasFromConfig)
}
