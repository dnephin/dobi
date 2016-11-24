package client

import (
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

//go:generate mockgen -source iface.go -destination mock_iface.go -package client

// DockerClient is the Docker API Client interface used by tasks
type DockerClient interface {
	BuildImage(docker.BuildImageOptions) error
	InspectImage(string) (*docker.Image, error)
	PushImage(docker.PushImageOptions, docker.AuthConfiguration) error
	PullImage(docker.PullImageOptions, docker.AuthConfiguration) error
	RemoveImage(string) error
	TagImage(string, docker.TagImageOptions) error

	AttachToContainerNonBlocking(docker.AttachToContainerOptions) (docker.CloseWaiter, error)
	CreateContainer(docker.CreateContainerOptions) (*docker.Container, error)
	KillContainer(docker.KillContainerOptions) error
	RemoveContainer(docker.RemoveContainerOptions) error
	StartContainer(string, *docker.HostConfig) error
	WaitContainer(string) (int, error)
	CreateService(opts docker.CreateServiceOptions) (*swarm.Service, error)
	ListServices(opts docker.ListServicesOptions) ([]swarm.Service, error)
	UpdateService(id string, opts docker.UpdateServiceOptions) error
	RemoveService(opts docker.RemoveServiceOptions) error
}
