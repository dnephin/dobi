package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"text/template"
	"path/filepath"
	"strings"
	"path"
	"net/http"
	"io/ioutil"
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
			goProject := TemplatingInfo{Name:path.Base(AbsolutePath), Path:str}

			res, err := http.Get("https://gist.githubusercontent.com/cescoferraro/f3a08f7e4c20b8243589e9d7a048877f/raw/8a87246bf8cbdcfc30e0ff53f73cdbdfaed5f74a/Dockerfile")
			if err != nil {
				return err
			}
			defer res.Body.Close()


			f, err := os.Create(path.Join(AbsolutePath, "Dockerfile"))
			if err != nil {
				return err
			}

			dockerfileDAta, err := ioutil.ReadAll(res.Body)

			if err != nil {
				return err
			}

			_, err = f.Write(dockerfileDAta)
			if err != nil {
				return err
			}

			f.Sync()

			response, err := http.Get("https://gist.githubusercontent.com/cescoferraro/f3a08f7e4c20b8243589e9d7a048877f/raw/d804ab9fcdbd9e806c4d8d37778f1af0de09e317/dobi.yaml")
			if err != nil {
				return err
			}
			defer response.Body.Close()


			htmlData, err := ioutil.ReadAll(response.Body)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			tmpl, err := template.New("").Parse(string(htmlData))

			if err != nil {
				panic(err)
			}

			f, err = os.Create(path.Join(AbsolutePath, "dobi.yaml"))
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
