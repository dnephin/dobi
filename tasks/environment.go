package tasks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dnephin/dobi/config"
	"github.com/kballard/go-shellquote"
)

// ExecEnv is a data object which contains variables for an ExecuteContext
type ExecEnv struct {
	ExecID string
}

// NewExecEnv returns a new ExecEnv from a Config
func NewExecEnv(cfg *config.Config) (*ExecEnv, error) {
	execID, err := getExecID(cfg.Meta.UniqueExecId, cfg.WorkingDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to generated unique execution id: %s", err)
	}
	return &ExecEnv{ExecID: execID}, nil
}

func getExecID(cmd string, workingDir string) (string, error) {
	if cmd == "" {
		return defaultExecID(workingDir), nil
	}

	stdout, err := runCommand(cmd)
	if err != nil {
		return "", err
	}

	return validateExecID(stdout)
}

func runCommand(cmdString string) (string, error) {
	cmdSlice, err := shellquote.Split(cmdString)
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

func defaultExecID(workingDir string) string {
	// TODO: cross-platform user name
	return fmt.Sprintf("%s-%s", filepath.Base(workingDir), os.Getenv("USER"))
}
