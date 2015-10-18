package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type BuildConfig struct {
	Dockerfile string
	Context    string
	Args       map[string]string
}

type ComposeConfig struct {
	Project  string
	Run      string
	Config   interface{}
	Filename string
}

type RunConfig struct {
	Volumes []string
	Command string
}

type Step struct {
	Name    string
	Image   string
	Build   BuildConfig
	Run     RunConfig
	Compose ComposeConfig
}

type Steps []*Step

type Config struct {
	Steps Steps
}

func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
