package execenv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dnephin/dobi/logging"
	git "github.com/gogits/git-module"
	"github.com/metakeule/fmtdate"
	fasttmpl "github.com/valyala/fasttemplate"
)

const (
	startTag     = "{"
	endTag       = "}"
	execIDEnvVar = "DOBI_EXEC_ID"
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

// ResolveSlice resolves all strings in the slice
func (e *ExecEnv) ResolveSlice(tmpls []string) ([]string, error) {
	resolved := []string{}
	for _, item := range tmpls {
		item, err := e.Resolve(item)
		if err != nil {
			return tmpls, err
		}
		resolved = append(resolved, item)
	}
	return resolved, nil
}

func (e *ExecEnv) templateContext(out io.Writer, tag string) (int, error) {
	tag, defValue, hasDefault := splitDefault(tag)

	write := func(val string) (int, error) {
		if val == "" {
			if !hasDefault {
				return 0, fmt.Errorf("A value is required for variable %q", tag)
			}
			val = defValue
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

		logging.Log.Warnf("Failed to get variable \"git.%s\", using default", tag)
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

func splitDefault(tag string) (string, string, bool) {
	parts := strings.Split(tag, ":")
	if len(parts) == 1 {
		return tag, "", false
	}
	last := len(parts) - 1
	return strings.Join(parts[:last], ":"), parts[last], true
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
func NewExecEnvFromConfig(execID, project, workingDir string) (*ExecEnv, error) {
	env := NewExecEnv(defaultExecID(), getProjectName(project, workingDir))
	var err error
	env.ExecID, err = getExecID(execID, env)
	return env, err
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
	project = filepath.Base(workingDir)
	logging.Log.Warnf("meta.project is not set. Using default %q.", project)
	return project
}

func getExecID(execID string, env *ExecEnv) (string, error) {
	var err error

	if value, exists := os.LookupEnv(execIDEnvVar); exists {
		return validateExecID(value)
	}
	if execID == "" {
		return env.ExecID, nil
	}

	execID, err = env.Resolve(execID)
	if err != nil {
		return "", err
	}
	return validateExecID(execID)
}

func validateExecID(output string) (string, error) {
	output = strings.TrimSpace(output)

	if output == "" {
		return "", fmt.Errorf("exec-id template was empty after rendering")
	}
	lines := len(strings.Split(output, "\n"))
	if lines > 1 {
		return "", fmt.Errorf(
			"exec-id template rendered to %v lines, expected only one", lines)
	}

	return output, nil
}

func defaultExecID() string {
	// TODO: cross-platform user name
	return os.Getenv("USER")
}
