package image

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type stepInjector string

var customSteṕs = map[string]string{
	"GLIDE":  "{{ .GLIDE }}",
	"GO-GET": "{{ .GOGET }}",
	"GO-BIN": "{{.GOBIN}}",
}

func (step stepInjector) GLIDE() string {
	return "curl -sL https://github.com/Masterminds/glide/releases/download/v" +
		string(step) + "/glide-v" + string(step) +
		"-linux-amd64.tar.gz | " +
		"tar -xz linux-amd64/glide && " +
		"mv linux-amd64/glide /usr/bin/glide && chmod +x /usr/bin/glide"
}

func (step stepInjector) GOBIN() string {
	strps := strings.Split(string(step), "/")
	return "go get -u " + string(step) + " && cp /go/bin/" +
		strps[len(strps)-1] + "/usr/bin/ && rm -rf /go/src/* /go/pkg/* /go/bin/*"
}

func (step stepInjector) GOGET() string {
	return "go get -u " + string(step)
}

func (t *Task) templateStep(key string, STEP map[string]string) error {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(customSteṕs[key])
	if err != nil {
		return err
	}

	err = tmpl.Execute(buf, stepInjector(STEP[key]))
	if err != nil {
		return err
	}
	delete(STEP, key)
	key = "RUN"
	// I dont know why but text/template replaces & and + from template function
	// temporary bug fix. There has to be a better way to handle this.
	// I dont know which other chars would incur to the same error.
	STEP[key] = strings.Replace(strings.Replace(buf.String(), "&amp;", "&", -1), "&#43;", "+", -1)
	return nil
}

func (t *Task) replaceCustomSteps() error {
	for _, STEP := range t.config.Steps {
		for key := range STEP {
			if _, ok := customSteṕs[key]; ok {
				err := t.templateStep(key, STEP)
				if err != nil {
					return fmt.Errorf("template step"+key+"s failed with error", err)
				}
			}
		}
	}
	return nil
}
