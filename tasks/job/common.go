package job

import (
	"fmt"

	"github.com/dnephin/dobi/tasks/client"
	"github.com/dnephin/dobi/tasks/context"
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

// containerName returns the name of the container
func containerName(ctx *context.ExecuteContext, name string) string {
	return fmt.Sprintf("%s-%s", ctx.Env.Unique(), name)
}

// removeContainer removes a container by ID, and logs a warning if the remove
// fails.
func removeContainer(
	logger *log.Entry,
	client client.DockerClient,
	containerID string,
) (bool, error) {
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
