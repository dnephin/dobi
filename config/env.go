package config

import (
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
//         captures:
//             - job: generate-version
//               variable: VERSION
//
// name: env
type EnvConfig struct {
	// Files List of files which contain environment variables
	// type: list of filenames
	Files []string
	// Variables List of environment variable ``key=value`` pairs
	// type: list of environment variables
	Variables []string
	// Capture Capture the stdout of **job** resources in an environment
	// variable
	// type: list of capture definitions
	// example: ``[{job: name-of-job, variable: TARGET_VAR}]``
	Captures []VariableCaptures
	describable
}

// VariableCaptures defines the structure for capturing variables from the
// stdout of a job resource
type VariableCaptures struct {
	// Job name of the job to capture
	Job string
	// Variable used to store the captured value
	Variable string
}

// Dependencies returns the list of job dependencies
func (c *EnvConfig) Dependencies() []string {
	jobs := []string{}
	for _, capture := range c.Captures {
		jobs = append(jobs, capture.Job)
	}
	return jobs
}

// Validate runs config validation
func (c *EnvConfig) Validate(pth.Path, *Config) *pth.Error {
	return nil
}

// Resolve resolves variables in the config
func (c *EnvConfig) Resolve(*execenv.ExecEnv) (Resource, error) {
	// TODO: support variables in Variables, Files, Captures.Variables
	return c, nil
}

func (c *EnvConfig) String() string {
	// TODO: more specific string
	return "Set variables"
}

func envFromConfig(name string, values map[string]interface{}) (Resource, error) {
	cnf := &EnvConfig{}
	return cnf, configtf.Transform(name, values, cnf)
}

func init() {
	RegisterResource("env", envFromConfig)
}
