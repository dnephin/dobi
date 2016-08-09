package client

import (
	docker "github.com/fsouza/go-dockerclient"
)

// DockerClient is the Docker API Client interface used by tasks
type DockerClient interface {
	BuildImage(docker.BuildImageOptions) error
	InspectImage(string) (*docker.Image, error)
	PushImage(docker.PushImageOptions, docker.AuthConfiguration) error
	RemoveImage(string) error
	TagImage(string, docker.TagImageOptions) error

	AttachToContainerNonBlocking(docker.AttachToContainerOptions) (docker.CloseWaiter, error)
	CreateContainer(docker.CreateContainerOptions) (*docker.Container, error)
	KillContainer(docker.KillContainerOptions) error
	RemoveContainer(docker.RemoveContainerOptions) error
	StartContainer(string, *docker.HostConfig) error
	WaitContainer(string) (int, error)
}
