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

{{ range .Fields }}

**{{ .Name }}** {{ if .IsRequired }} (required) {{ end }}

    :type: {{ .Type }}
    {{ with .Format }}:format: {{ . }}) {{ end }}
    {{ with .Default }}:default: {{ . }} {{ end }}
    :description: {{ .Description }}
    {{ with .Example }}:example: {{ . }} {{ end }}

{{ end }}
`
