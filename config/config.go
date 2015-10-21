package config

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
)

type BuildConfig struct {
	Args       map[string]string
	Context    string
	Dockerfile string
	Pull       bool
	Tags       []string
}

type ComposeConfig struct {
	Config   interface{}
	Filename string
	Project  string
	Run      string
}

type RunConfig struct {
	Command    []string
	Volumes    []string
	Privileged bool
}

type StepConfig struct {
	Build   *BuildConfig
	Compose *ComposeConfig
	Image   string
	Name    string
	Pull    bool
	Run     *RunConfig
}

type Steps []*StepConfig

type Config struct {
	Steps Steps
}

func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	if err = validate(config); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{"filename": filename}).Info("Configuration loaded")
	return config, nil
}

func validate(config *Config) error {
	// TODO: if pull is true, image must be set
	// TODO: run and compose actions are mutually exclusive
	// TODO: compose config and filename are mutually exclusive
	return nil
}
