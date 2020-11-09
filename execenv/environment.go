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
	"github.com/pkg/errors"
	fasttmpl "github.com/valyala/fasttemplate"
)

const (
	startTag     = "{"
	endTag       = "}"
	execIDEnvVar = "DOBI_EXEC_ID"
)

// ExecEnv is a data object which contains variables for an ExecuteContext
type ExecEnv struct {
	ExecID     string
	Project    string
	tmplCache  map[string]string
	workingDir string
	startTime  time.Time
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
	for _, tmpl := range tmpls {
		item, err := e.Resolve(tmpl)
		if err != nil {
			return tmpls, err
		}
		resolved = append(resolved, item)
	}
	return resolved, nil
}

// nolint: gocyclo
func (e *ExecEnv) templateContext(out io.Writer, tag string) (int, error) {
	tag, defValue, hasDefault := splitDefault(tag)

	write := func(val string, err error) (int, error) {
		if err != nil {
			return 0, err
		}
		if val == "" {
			if !hasDefault {
				return 0, fmt.Errorf("a value is required for variable %q", tag)
			}
			val = defValue
		}
		return out.Write(bytes.NewBufferString(val).Bytes())
	}

	prefix, suffix := splitPrefix(tag)
	switch prefix {
	case "env":
		return write(os.Getenv(suffix), nil)
	case "git":
		return valueFromGit(out, e.workingDir, suffix, defValue)
	case "time":
		return write(fmtdate.Format(suffix, e.startTime), nil)
	case "fs":
		val, err := ValueFromFilesystem(suffix, e.workingDir)
		return write(val, err)
	case "user":
		val, err := valueFromUser(suffix)
		return write(val, err)
	}

	switch tag {
	case "unique":
		return write(e.Unique(), nil)
	case "project":
		return write(e.Project, nil)
	case "exec-id":
		return write(e.ExecID, nil)
	default:
		return 0, errors.Errorf("unknown variable %q", tag)
	}
}

// ValueFromFilesystem can return either `cwd` or `projectdir`
func ValueFromFilesystem(name string, workingdir string) (string, error) {
	switch name {
	case "cwd":
		return os.Getwd()
	case "projectdir":
		return workingdir, nil
	default:
		return "", errors.Errorf("unknown variable \"fs.%s\"", name)
	}
}

// nolint: gocyclo
func valueFromGit(out io.Writer, cwd string, tag, defValue string) (int, error) {
	writeValue := func(value string) (int, error) {
		return out.Write(bytes.NewBufferString(value).Bytes())
	}

	writeError := func(err error) (int, error) {
		if defValue == "" {
			return 0, fmt.Errorf("failed resolving variable {git.%s}: %s", tag, err)
		}

		logging.Log.Warnf("Failed to get variable \"git.%s\", using default", tag)
		return writeValue(defValue)
	}

	repo, err := git.OpenRepository(cwd)
	if err != nil {
		return writeError(err)
	}

	switch tag {
	case "branch":
		branch, err := repo.GetHEADBranch()
		if err != nil {
			return writeError(err)
		}
		return writeValue(branch.Name)
	case "sha":
		commit, err := repo.GetCommit("HEAD")
		if err != nil {
			return writeError(err)
		}
		return writeValue(commit.ID.String())
	case "short-sha":
		commit, err := repo.GetCommit("HEAD")
		if err != nil {
			return writeError(err)
		}
		return writeValue(commit.ID.String()[:10])
	default:
		return 0, errors.Errorf("unknown variable \"git.%s\"", tag)
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
	index := strings.Index(tag, ".")
	switch index {
	case -1, 0, len(tag) - 1:
		return "", tag
	default:
		return tag[:index], tag[index+1:]
	}
}

// NewExecEnvFromConfig returns a new ExecEnv from a Config
func NewExecEnvFromConfig(execID, project, workingDir string) (*ExecEnv, error) {
	env := NewExecEnv(defaultExecID(), getProjectName(project, workingDir), workingDir)
	var err error
	env.ExecID, err = getExecID(execID, env)
	return env, err
}

// NewExecEnv returns a new ExecEnv from values
func NewExecEnv(execID, project, workingDir string) *ExecEnv {
	return &ExecEnv{
		ExecID:     execID,
		Project:    project,
		tmplCache:  make(map[string]string),
		startTime:  time.Now(),
		workingDir: workingDir,
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
	username, err := getUserName()
	if err == nil {
		return username
	}
	return os.Getenv("USER")
}
