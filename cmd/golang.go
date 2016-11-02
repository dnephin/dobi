package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"text/template"
	"path/filepath"
	"strings"
	"path"
)

type TemplatingInfo struct {
	Name string
	Path string
}

func newGolangCommand(opts *dobiOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "golang",
		Short: "Run the remove action for all resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				return fmt.Errorf("GOPATH NOT SET!")
			}

			AbsolutePath, err := os.Getwd()
			if err != nil {
				return err
			}

			j, err := filepath.Rel(gopath, AbsolutePath)
			if err != nil {
				return fmt.Errorf("You are outside the gopath baby")
			}

			str := strings.Replace(j, "src/", "", -1)
			goProject := TemplatingInfo{Name:path.Base(AbsolutePath) ,Path:str}
			tmpl, err := template.New("").Parse(gotemplate)

			if err != nil {
				panic(err)
			}

			f, err := os.Create(path.Join(AbsolutePath,"dobi-golang.yaml"))
			if err != nil {
				return err
			}


			err = tmpl.Execute(f, goProject)
			if err != nil {
				panic(err)
			}

			return nil
		},
	}
	return cmd
}


var gotemplate = `meta:
    project: {{.Name}}
    default: binary

env=vars:
    variables: [VERSION=0.0.1]


mount=source:
    bind: .
    path: /go/src/{{.Path}}

mount=dist:
    bind: ./dist/bin/
    path: /go/bin/


image=binary-builder:
    image: golang:1.7
    pull: "once"

job=binary:
    use: binary-builder
    working-dir: /go/src/{{.Path}}
    artifact: ./dist/bin/{{.Name}}
    mounts: [source, dist]
    depends: [vars]
    command: go build -o dist/bin/{{.Name}}
    description: "Build a binary
`
