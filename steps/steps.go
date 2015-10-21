package steps

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	cf "github.com/dnephin/buildpipe/config"
	"github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/errwrap"
	"gopkg.in/yaml.v2"
)

type Action func() error

type Step struct {
	client *docker.Client
	config *cf.StepConfig
}

func (step *Step) run() error {
	log.WithFields(log.Fields{"name": step.config.Name}).Info("Running step")

	actions := map[string]Action{
		"Pull image":    step.pullImage,
		"Build image":   step.buildImage,
		"Run container": step.runContainer,
		"Run compose":   step.runCompose,
	}
	for name, action := range actions {
		if err := action(); err != nil {
			return errwrap.Wrapf(name+" step failed: {{err}}", err)
		}
	}
	return nil
}

func (step *Step) pullImage() error {
	if !step.config.Pull {
		return nil
	}

	log.WithFields(log.Fields{"image": step.config.Image}).Info("Pulling image")
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

	log.WithFields(log.Fields{"image": step.config.Image}).Info("Building image")
	// TODO: add first tag if there are build tags
	// TODO: support other build options
	err := step.client.BuildImage(docker.BuildImageOptions{
		Name:       step.config.Image,
		Dockerfile: step.config.Build.Dockerfile,
		Pull:       step.config.Build.Pull,
		ContextDir: step.config.Build.Context,
		// TODO: where should this output go?
		OutputStream: os.Stdout,
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
	log.WithFields(log.Fields{
		"cmd": step.config.Run.Command,
	}).Info("Running command")

	// TODO: support other run options
	container, err := step.client.CreateContainer(docker.CreateContainerOptions{
		// TODO: give the container a unique name based on UNIQUE_ID and step name?
		Config: &docker.Config{
			Cmd:   step.config.Run.Command,
			Image: step.config.Image,
		},
		HostConfig: &docker.HostConfig{
			// TODO: support relative paths or {{PWD}}
			Binds:      step.config.Run.Volumes,
			Privileged: step.config.Run.Privileged,
		},
	})
	if err != nil {
		return err
	}
	if err := step.client.StartContainer(container.ID, nil); err != nil {
		return err
	}

	if err := step.client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Logs:         true,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
	}); err != nil {
		return err
	}
	if err := step.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            container.ID,
		RemoveVolumes: true,
	}); err != nil {
		log.WithFields(log.Fields{
			"container": container.ID,
			"error":     err.Error(),
		}).Warn("Failed to remove container.")
	}
	return nil
}

func (step *Step) runCompose() error {
	if step.config.Compose == nil {
		return nil
	}
	log.WithFields(log.Fields{
		"run": step.config.Compose.Run,
	}).Info("Running docker-compose")

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
	log.Info("Docker client created")

	for _, stepConfig := range config.Steps {
		step := Step{client: client, config: stepConfig}
		if err = step.run(); err != nil {
			return err
		}
	}
	return nil
}
