package cmd

import (
	"github.com/spf13/cobra"
	"fmt"

	"github.com/docker/docker/pkg/term"
	docker "github.com/fsouza/go-dockerclient"
	"io"
	"os"
	"io/ioutil"
	"log"
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
	dockeropts := docker.CreateContainerOptions{
		Name: "fdgdfg",
		Config: &docker.Config{
			//Cmd:          []string{"pip","install","cookiecutter"},
			Cmd:          []string{"ls"},
			Image:        "busybox:latest",
			OpenStdin:    true,
			Tty:          true,
			AttachStdin:  true,
			StdinOnce:    true,
			AttachStderr: true,
			AttachStdout: true,
			WorkingDir:   "/",
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
		OutputStream: io.MultiWriter(os.Stdout, os.Stdout),
		ErrorStream: os.Stdout,
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