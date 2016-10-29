package logging

import (
	log "github.com/Sirupsen/logrus"
)

var (
	// Log is the logger used by dobi
	Log = log.New()
)

// ForTask returns a logger for a task which implemented LogRepresenter. The
// logger has the task added as the `task` field.
func ForTask(repr LogRepresenter) *log.Entry {
	return Log.WithFields(log.Fields{"task": repr})
}
