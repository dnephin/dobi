package docs

import (
	"bytes"
	"text/template"
)

// FormatRst returns the configType formatted as ReStructuredType
func FormatRst(configType ConfigType) (string, error) {
	buf := new(bytes.Buffer)
	tmpl := template.New("config-output")
	var err error

	if tmpl, err = tmpl.Parse(rawTmpl); err != nil {
		return "", err
	}
	err = tmpl.Execute(buf, configType)
	return buf.String(), err
}

var rawTmpl = `

{{ .Name }}
-----------

{{ .Description }}

{{ with .Example }}
Example
~~~~~~~

{{ . }}
{{ end }}

Fields
~~~~~~

{{ range .Fields }}

**{{ .Name }}** {{ if .IsRequired }} (required) {{ end }}

    Type: {{ .Type }} {{ with .Format }} (format: {{ . }}) {{ end }}

	{{ .Description }}

	{{ with .Example }} Example: {{ . }} {{ end }}

{{ end }}
`
