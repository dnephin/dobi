package config

import (
	"fmt"
	"strings"

	"github.com/dnephin/configtf"
	pth "github.com/dnephin/configtf/path"
	"github.com/dnephin/dobi/execenv"
)

// EnvConfig An **env** resource provides environment variables to **job** and
// **compose** resources.
//
// example: Define some variables for a ``job``
//
// .. code-block:: yaml
//
//     env=settings:
//         files: [local.env]
//         variables: [PORT=3838, HOST=stage]
//
// name: env
type EnvConfig struct {
	// Files List of files which contain environment variables
	// type: list of filenames
	Files []string
	// Variables List of environment variable ``key=value`` pairs
	// type: list of environment variables
	Variables []string
	describable
}

// Dependencies returns the list of job dependencies
func (c *EnvConfig) Dependencies() []string {
	return []string{}
}

// Validate runs config validation
func (c *EnvConfig) Validate(pth.Path, *Config) *pth.Error {
	return nil
}

// Resolve resolves variables in the config
func (c *EnvConfig) Resolve(*execenv.ExecEnv) (Resource, error) {
	// TODO: support variables in Variables, Files
	return c, nil
}

func (c *EnvConfig) String() string {
	return fmt.Sprintf(
		"Set vars from: %s and set: %s",
		strings.Join(c.Files, ", "), strings.Join(c.Variables, ", "))
}

func envFromConfig(name string, values map[string]interface{}) (Resource, error) {
	cnf := &EnvConfig{}
	return cnf, configtf.Transform(name, values, cnf)
}

func init() {
	RegisterResource("env", envFromConfig)
}
