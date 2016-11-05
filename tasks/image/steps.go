package image

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

var customSteṕs = map[string]string{
	"GLIDE":  "curl -sL https://github.com/Masterminds/glide/releases/download/v{{ .GLIDE }}/glide-v{{ .GLIDE}}-linux-amd64.tar.gz | tar -xz linux-amd64/glide && mv linux-amd64/glide /usr/bin/glide && chmod +x /usr/bin/glide",
	"GO-GET": "go get  {{ . }}",
	"GO-BIN": "go get -u {{ . }} && cp /go/bin/{{ .GOBIN }} /usr/bin/ && rm -rf /go/src/* /go/pkg/* /go/bin/*",
}

type stepInjector string

func (step stepInjector) GLIDE() string {
	if string(step) == " install" {

	}
	return string(step)
}
func (step stepInjector) GOBIN() string {
	strps := strings.Split(string(step), "/")
	return strps[len(strps)-1]
}

func (t *Task) templateStep(STEP map[string]string, key, value string) error {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(customSteṕs[key])
	if err != nil {
		return err
	}

	err = tmpl.Execute(buf, stepInjector(value))
	if err != nil {
		return err
	}
	delete(STEP, key)
	key = "RUN"
	STEP[key] = buf.String()
	return nil
}

func (t *Task) replaceCustomSteps() error {
	for _, STEP := range t.config.Steps {
		for key, value := range STEP {
			if _, ok := customSteṕs[key]; ok {
				err := t.templateStep(STEP, key, value)
				if err != nil {
					return fmt.Errorf("template step"+key+"s failed with error", err)
				}
			}
		}
	}
	return nil
}
