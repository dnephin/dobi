package job

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
)

// ContainerName returns the name of the container
func ContainerName(ctx *context.ExecuteContext, name string) string {
	return fmt.Sprintf("%s-%s", ctx.Env.Unique(), name)
}

// RemoveContainer removes a container
func RemoveContainer(logger *log.Entry, client client.DockerClient, containerID string, expectContainer bool) {
	logger.Debug("Removing container")
	err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerID,
		RemoveVolumes: true,
	})
	switch err.(type) {
	case *docker.NoSuchContainer:
		if !expectContainer {
			return
		}
	case nil:
		return
	}
	logger.WithFields(log.Fields{"container": containerID}).Warnf(
		"Failed to remove container: %s", err)
}
