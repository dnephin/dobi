package image

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dnephin/dobi/config"
	"github.com/dnephin/dobi/tasks/context"
	yaml "gopkg.in/yaml.v2"
)

const (
	imageRecordDir = ".dobi/images"
	all            = -1
)

type imageModifiedRecord struct {
	ImageID      string
	LastPull     *time.Time  `yaml:",omitempty"`
	Info         os.FileInfo `yaml:",omitempty"`
	LastModified *time.Time  `yaml:",omitempty"`
}

func (r *imageModifiedRecord) WasModifiedBefore(ctx *context.ExecuteContext, mtime time.Time) bool { // nolint: lll
	if ctx.Env.Hosted && r.LastModified != nil {
		return r.LastModified.Before(mtime)
	}
	return r.Info.ModTime().Before(mtime)
}

func updateImageRecord(ctx *context.ExecuteContext, path string, record imageModifiedRecord) error { // nolint: lll
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	if ctx.Env.Hosted {
		record.LastModified = now()
	}

	bytes, err := yaml.Marshal(record)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

// TODO: verify error message are sufficient
func getImageRecord(filepath string) (imageModifiedRecord, error) {
	record := imageModifiedRecord{}
	var err error

	record.Info, err = os.Stat(filepath)
	if err != nil {
		return record, err
	}

	recordBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return record, err
	}

	return record, yaml.Unmarshal(recordBytes, &record)
}

func recordPath(ctx *context.ExecuteContext, conf *config.ImageConfig) string {
	return recordPathForTag(ctx.WorkingDir, GetImageName(ctx, conf))
}

func recordPathForTag(workdir string, tag string) string {
	imageName := strings.Replace(tag, "/", " ", all)
	imageName = strings.Replace(imageName, ":", " ", all)
	return filepath.Join(workdir, imageRecordDir, imageName)
}
