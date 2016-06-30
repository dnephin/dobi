package logging

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 90
)

type LogRepresenter interface {
	Repr() string
}

type Formatter struct{}

func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	buff := &bytes.Buffer{}
	buff.WriteString(writeLevel(entry.Level))
	buff.WriteString(writeData(entry.Data))
	buff.WriteString(entry.Message)
	buff.WriteString("\n")
	return buff.Bytes(), nil
}

func withColor(color int, msg string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, msg)
}

func writeLevel(level log.Level) string {
	switch level {
	case log.DebugLevel:
		return fmt.Sprintf("[%s] ", withColor(gray, "DEBUG"))
	case log.WarnLevel:
		return fmt.Sprintf("[%s] ", withColor(yellow, "WARN"))
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		return fmt.Sprintf("[%s] ", withColor(red, "ERROR"))
	default:
		return ""
	}
}

func writeData(fields log.Fields) string {
	buff := []string{}
	for key, value := range fields {
		switch value := value.(type) {
		case LogRepresenter:
			buff = append(buff, value.Repr())
		default:
			buff = append(buff, fmt.Sprintf("%v=%v", key, value))
		}
	}
	if len(buff) > 0 {
		buff = append(buff, "")
	}
	return strings.Join(buff, " ")
}
