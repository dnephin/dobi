package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/config/tform/docs"
)

var (
	basePath = "docs/gen/config/"
)

func write(filepath string, source interface{}) error {
	content, err := docs.Generate(source, docs.ReStructuredText)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, bytes.NewBufferString(content).Bytes(), 0644)
}

func writeDocs() error {
	for _, item := range []struct {
		filename string
		source   interface{}
	}{
		{"meta.rst", config.MetaConfig{}},
		{"alias.rst", config.AliasConfig{}},
		{"compose.rst", config.ComposeConfig{}},
		{"image.rst", config.ImageConfig{}},
		{"mount.rst", config.MountConfig{}},
		{"job.rst", config.JobConfig{}},
	} {
		fmt.Printf("Generating doc %q\n", basePath+item.filename)
		if err := write(basePath+item.filename, item.source); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := writeDocs(); err != nil {
		fmt.Printf("Failed to generate docs: %s\n", err)
		os.Exit(1)
	}
}
