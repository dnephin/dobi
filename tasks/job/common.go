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

// RemoveContainer removes a container by ID, and logs a warning if the remove
// fails.
func RemoveContainer(logger *log.Entry, client client.DockerClient, containerID string) (bool, error) {
	logger.Debug("Removing container")
	err := client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerID,
		RemoveVolumes: true,
		Force:         true,
	})
	switch err.(type) {
	case *docker.NoSuchContainer:
		return false, nil
	case nil:
		return true, nil
	}
	logger.WithFields(log.Fields{"container": containerID}).Warnf(
		"Failed to remove container: %s", err)
	return false, err
}
