package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/docker/docker/pkg/term"
	docker "github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"log"
	"os"
)

func newInitCommand(opts *dobiOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Run the remove action for all resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.cookie, "cookie", "c", "https://github.com/cescoferraro/go-dobi-cutter", "Url to")
	return cmd
}

func runInit(opts *dobiOptions) error {
	client, err := buildClient()
	if err != nil {
		return fmt.Errorf("Failed to create client: %s", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dockeropts := docker.CreateContainerOptions{
		Name: "fdgdfg",
		Config: &docker.Config{
			Cmd:        []string{"gh:cescoferraro/go-dobi-cutter"},
			Image:      "cescoferraro/cookiecutter",
			Tty:        true,
			WorkingDir: "/srv/app",
			Mounts:     []docker.Mount{{Source: pwd, Destination: "/srv/app"}},
		},
		HostConfig: &docker.HostConfig{},
	}
	container, err := client.CreateContainer(dockeropts)
	if err != nil {
		return err
	}
	err = client.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		return err
	}

	_, err = client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stdout,
		InputStream:  ioutil.NopCloser(os.Stdin),
		Stream:       true,
		Stdin:        true,
		RawTerminal:  true,
		Stdout:       true,
		Stderr:       true,
	})
	if err != nil {
		return fmt.Errorf("Failed attaching to container %q: %s", container.Name, err)
	}

	status, err := client.WaitContainer(container.ID)
	if err != nil {
		return fmt.Errorf("Failed to wait on container exit: %s", err)
	}
	if status != 0 {
		return fmt.Errorf("Exited with non-zero status code %d", status)
	}

	inFd, _ := term.GetFdInfo(os.Stdin)
	state, err := term.SetRawTerminal(inFd)
	if err != nil {
		return err
	}
	defer func() {
		if err := term.RestoreTerminal(inFd, state); err != nil {
			log.Println(err.Error())
		}
	}()

	return nil
}
