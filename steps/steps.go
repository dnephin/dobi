package steps

import (
	"os/exec"
	"strings"

	cf "github.com/dnephin/buildpipe/config"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
)

type Action func() error

type Step struct {
	client *docker.Client
	config *cf.StepConfig
}

func (step *Step) run() error {
	actions := []Action{
		step.pullImage,
		step.buildImage,
		step.runContainer,
		step.runCompose,
	}
	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}
	return nil
}

func (step *Step) pullImage() error {
	if !step.config.Pull {
		return nil
	}
	return step.client.PullImage(docker.PullImageOptions{
		Repository: step.config.Image,
		// TODO: OutputStream to stdout?
		// TODO: support Auth
	}, docker.AuthConfiguration{})
}

func (step *Step) buildImage() error {
	if step.config.Build == nil {
		return nil
	}

	// TODO: add first tag if there are build tags
	// TODO: support other build options
	err := step.client.BuildImage(docker.BuildImageOptions{
		Name:       step.config.Image,
		Dockerfile: step.config.Build.Dockerfile,
		Pull:       step.config.Build.Pull,
		ContextDir: step.config.Build.Context,
		// TODO: OutputStream
	})
	if err != nil {
		return err
	}
	return step.tagImage()
}

// TODO:
func (step *Step) tagImage() error {
	return nil
}

func (step *Step) runContainer() error {
	if step.config.Run == nil {
		return nil
	}

	// TODO: support other run options
	container, err := step.client.CreateContainer(docker.CreateContainerOptions{
		// TODO: give the container a unique name based on UNIQUE_ID and step name?
		Config: &docker.Config{
			Cmd:    step.config.Run.Command,
			Mounts: step.config.Run.VolumeMounts(),
		},
	})
	if err != nil {
		return err
	}

	return step.client.StartContainer(container.ID, nil)
}

func (step *Step) runCompose() error {
	if step.config.Compose == nil {
		return nil
	}

	// TODO: accept an alterate path to binary
	cmd := exec.Command(
		"docker-compose",
		"-p", step.config.Compose.Project,
		"-f", "-",
		"run",
		// TODO: tokenize command
		step.config.Compose.Run)

	// TODO: support filename instead of just inline config
	composeFile, err := yaml.Marshal(step.config.Compose.Config)
	if err != nil {
		return err
	}
	cmd.Stdin = strings.NewReader(string(composeFile))

	if err = cmd.Start(); err != nil {
		return err
	}

	// TODO: print stdout/stderr
	return cmd.Wait()
}

func Run(config *cf.Config) error {
	// TODO: args for client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	for _, stepConfig := range config.Steps {
		step := Step{client: client, config: stepConfig}
		if err = step.run(); err != nil {
			return err
		}
	}
	return nil
}
