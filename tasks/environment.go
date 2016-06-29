package tasks

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dnephin/dobi/config"
	git "github.com/gogits/git-module"
	shlex "github.com/kballard/go-shellquote"
	fasttmpl "github.com/valyala/fasttemplate"
	"github.com/metakeule/fmtdate"
)

const (
	startTag = "{"
	endTag   = "}"
)

// ExecEnv is a data object which contains variables for an ExecuteContext
type ExecEnv struct {
	ExecID    string
	Project   string
	tmplCache map[string]string
	startTime time.Time
}

// Unique returns a unique id for this execution
func (e *ExecEnv) Unique() string {
	return e.Project + "-" + e.ExecID
}

// Resolve template variables to a string value and cache the value
func (e *ExecEnv) Resolve(tmpl string) (string, error) {
	if val, ok := e.tmplCache[tmpl]; ok {
		return val, nil
	}

	template, err := fasttmpl.NewTemplate(tmpl, startTag, endTag)
	if err != nil {
		return "", err
	}

	buff := &bytes.Buffer{}
	_, err = template.ExecuteFunc(buff, e.templateContext)
	if err == nil {
		e.tmplCache[tmpl] = buff.String()
	}
	return buff.String(), err
}

func (e *ExecEnv) templateContext(out io.Writer, tag string) (int, error) {
	tag, defValue := splitDefault(tag)

	write := func(val string) (int, error) {
		if val == "" {
			val = defValue
		}
		if val == "" {
			return 0, fmt.Errorf("A value is required for variable %q", tag)
		}
		return out.Write(bytes.NewBufferString(val).Bytes())
	}

	prefix, suffix := splitPrefix(tag)
	switch prefix {
	case "env":
		return write(os.Getenv(suffix))
	case "git":
		return valueFromGit(out, suffix, defValue)
	case "time":
		return write(fmtdate.Format(suffix, e.startTime))
	}

	switch tag {
	case "unique":
		return write(e.Unique())
	case "project":
		return write(e.Project)
	case "exec-id":
		return write(e.ExecID)
	default:
		return 0, fmt.Errorf("Unknown variable %q", tag)
	}
}

func valueFromGit(out io.Writer, tag, defValue string) (int, error) {
	write := func(value string) (int, error) {
		return out.Write(bytes.NewBufferString(value).Bytes())
	}

	writeWithError := func(err error) (int, error) {
		if defValue == "" {
			return 0, err
		}

		log.Warnf("Failed to get variable \"git.%s\", using default", tag)
		return write(defValue)
	}

	repo, err := git.OpenRepository(".")
	if err != nil {
		return writeWithError(err)
	}

	switch tag {
	case "branch":
		branch, err := repo.GetHEADBranch()
		if err != nil {
			return writeWithError(err)
		}
		return write(branch.Name)
	case "sha":
		commit, err := repo.GetCommit("HEAD")
		if err != nil {
			return writeWithError(err)
		}
		return write(commit.ID.String())
	case "short-sha":
		commit, err := repo.GetCommit("HEAD")
		if err != nil {
			return writeWithError(err)
		}
		return write(commit.ID.String()[:10])
	default:
		return 0, fmt.Errorf("Unknown variable \"git.%s\"", tag)
	}
}

// GetVar returns a variable from the cache, or panics if it doesn't exist
func (e *ExecEnv) GetVar(tmpl string) string {
	val, ok := e.tmplCache[tmpl]
	if !ok {
		panic(fmt.Sprintf("Variables was not prepared %q", tmpl))
	}
	return val
}

func splitDefault(tag string) (string, string) {
	parts := strings.Split(tag, ":")
	if len(parts) == 1 {
		return tag, ""
	}
	last := len(parts) - 1
	return strings.Join(parts[:last], ":"), parts[last]
}

func splitPrefix(tag string) (string, string) {
	for _, prefix := range []string{"env", "git", "time"} {
		if strings.HasPrefix(tag, prefix+".") {
			return prefix, tag[len(prefix)+1:]
		}
	}
	return "", tag
}

// NewExecEnvFromConfig returns a new ExecEnv from a Config
func NewExecEnvFromConfig(cfg *config.Config) (*ExecEnv, error) {
	execID, err := getExecID(cfg.Meta.UniqueExecID)
	if err != nil {
		return nil, fmt.Errorf("Failed to generated unique execution id: %s", err)
	}
	project := getProjectName(cfg.Meta.Project, cfg.WorkingDir)
	return NewExecEnv(execID, project), nil
}

// NewExecEnv returns a new ExecEnv from values
func NewExecEnv(execID, project string) *ExecEnv {
	return &ExecEnv{
		ExecID:    execID,
		Project:   project,
		tmplCache: make(map[string]string),
		startTime: time.Now(),
	}
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
