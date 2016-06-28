package tasks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dnephin/dobi/config"
	shlex "github.com/kballard/go-shellquote"
)

// ExecEnv is a data object which contains variables for an ExecuteContext
type ExecEnv struct {
	ExecID  string
	Project string
}

// Unique returns a unique id for this execution
func (e *ExecEnv) Unique() string {
	return e.Project + "-" + e.ExecID
}

// NewExecEnv returns a new ExecEnv from a Config
func NewExecEnv(cfg *config.Config) (*ExecEnv, error) {
	execID, err := getExecID(cfg.Meta.UniqueExecID)
	if err != nil {
		return nil, fmt.Errorf("Failed to generated unique execution id: %s", err)
	}
	project := getProjectName(cfg.Meta.Project, cfg.WorkingDir)
	return &ExecEnv{ExecID: execID, Project: project}, nil
}

func getProjectName(project, workingDir string) string {
	if project != "" {
		return project
	}
	return filepath.Base(workingDir)
}

func getExecID(cmd string) (string, error) {
	if cmd == "" {
		return defaultExecID(), nil
	}

	stdout, err := runCommand(cmd)
	if err != nil {
		return "", err
	}

	return validateExecID(stdout)
}

func runCommand(cmdString string) (string, error) {
	cmdSlice, err := shlex.Split(cmdString)
	if err != nil {
		return "", fmt.Errorf("Failed to parse command: %s", err)
	}

	var stdout bytes.Buffer
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return stdout.String(), nil
}

func validateExecID(output string) (string, error) {
	output = strings.TrimSpace(output)

	if output == "" {
		return "", fmt.Errorf("Exec id command returned no output.")
	}
	lines := len(strings.Split(output, "\n"))
	if lines > 1 {
		return "", fmt.Errorf(
			"Exec id command returned %v lines, expected only one.", lines)
	}

	return output, nil
}

func defaultExecID() string {
	// TODO: cross-platform user name
	return os.Getenv("USER")
}
