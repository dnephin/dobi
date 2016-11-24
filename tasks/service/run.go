package service

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/logging"
	"github.com/dnephin/dobi/tasks/context"
	"github.com/dnephin/dobi/tasks/image"
	"github.com/dnephin/dobi/tasks/task"
	"github.com/dnephin/dobi/tasks/types"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"io"
	"strconv"
	"strings"
)

func newServeTask(name task.Name, conf config.Resource) types.Task {
	return &ServeTask{name: name, config: conf.(*config.ServiceConfig)}
}

// ServeTask is a task which runs a docker service
type ServeTask struct {
	types.NoStop
	name      task.Name
	config    *config.ServiceConfig
	outStream io.Writer
}

// Name returns the name of the task
func (t *ServeTask) Name() task.Name {
	return t.name
}

func (t *ServeTask) logger() *log.Entry {
	return logging.ForTask(t)
}

// Repr formats the task for logging
func (t *ServeTask) Repr() string {
	buff := &bytes.Buffer{}
	return fmt.Sprintf("%s%v", t.name.Format("job"), buff.String())
}

// Run the job command in a container
func (t *ServeTask) serviceOpts(ctx *context.ExecuteContext) docker.CreateServiceOptions {
	return docker.CreateServiceOptions{ServiceSpec: t.serviceSpec(ctx)}
}

// Run the job command in a container
func (t *ServeTask) endPoints() *swarm.EndpointSpec {
	var portConfigs []swarm.PortConfig
	for _, port := range t.config.Ports {
		split := strings.Split(port, ":")
		first, _ := strconv.Atoi(split[0])
		second, _ := strconv.Atoi(split[1])
		portConfigs = append(portConfigs, swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			TargetPort:    uint32(first),
			PublishedPort: uint32(second),
		})
	}
	return &swarm.EndpointSpec{
		Mode:  "vip",
		Ports: portConfigs,
	}
}

// Run the job command in a container
func (t *ServeTask) serviceSpec(ctx *context.ExecuteContext) swarm.ServiceSpec {
	name := ContainerName(ctx, t.name.Resource())
	imageName := image.GetImageName(ctx, ctx.Resources.Image(t.config.Use))
	replicas := uint64(t.config.Replicas)
	return swarm.ServiceSpec{
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicas,
			}},
		Annotations: swarm.Annotations{
			Name: name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: imageName,
			},
		},
		EndpointSpec: t.endPoints(),
	}
}

// Run the job command in a container
func (t *ServeTask) updateOpts(ctx *context.ExecuteContext, service swarm.Service) docker.UpdateServiceOptions {
	return docker.UpdateServiceOptions{Version: uint64(service.Version.Index), ServiceSpec: t.serviceSpec(ctx)}
}

func isServiceRunning(ctx *context.ExecuteContext, name string) (swarm.Service, bool, error) {
	services, err := ctx.Client.ListServices(docker.ListServicesOptions{})
	if err != nil {
		return swarm.Service{}, false, err
	}
	for _, service := range services {
		if service.Spec.Name == ContainerName(ctx, name) {
			return service, true, nil
		}
	}
	return swarm.Service{}, false, nil
}
func (t *ServeTask) hasChanged(ctx *context.ExecuteContext, service swarm.Service) bool {
	if int(*service.Spec.Mode.Replicated.Replicas) != t.config.Replicas {
		return true
	}
	return false
}

// Run containers as a docker service
func (t *ServeTask) Run(ctx *context.ExecuteContext, depsModified bool) (bool, error) {
	service, running, err := isServiceRunning(ctx, t.name.Resource())
	if err != nil {
		return false, err
	}
	if !running {
		_, err := ctx.Client.CreateService(t.serviceOpts(ctx))
		if err != nil {
			return false, err
		}
		t.logger().Info("created")
		return true, nil
	}
	if t.hasChanged(ctx, service) {
		err := ctx.Client.UpdateService(service.ID, t.updateOpts(ctx, service))
		if err != nil {
			return false, err
		}
		t.logger().Info("has changed")
		return true, nil
	}

	t.logger().Info("already running")
	return true, nil
}

// ContainerName returns the name of the container
func ContainerName(ctx *context.ExecuteContext, name string) string {
	return fmt.Sprintf("%s-%s", ctx.Env.Unique(), name)
}
