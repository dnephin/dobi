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
func RemoveContainer(logger *log.Entry, client client.DockerClient, containerID string) {
	logger.Debug("Removing container")
	if err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerID,
		RemoveVolumes: true,
	}); err != nil {
		logger.WithFields(log.Fields{"container": containerID}).Warnf(
			"Failed to remove container: %s", err)
	}
}
