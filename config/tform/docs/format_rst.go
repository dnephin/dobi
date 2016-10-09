package docs

import (
	"bytes"
	"strings"
	"text/template"
)

// FormatRst returns the configType formatted as ReStructuredType
func FormatRst(configType ConfigType) (string, error) {
	buf := new(bytes.Buffer)
	tmpl := template.New("config-output").Funcs(map[string]interface{}{
		"indent": indent,
		"repeat": repeat,
	})
	var err error

	if tmpl, err = tmpl.Parse(rawTmpl); err != nil {
		return "", err
	}
	err = tmpl.Execute(buf, configType)
	return buf.String(), err
}

func indent(spaces int, text string) string {
	lines := strings.Split(text, "\n")
	indent := strings.Repeat(" ", spaces)
	return indent + strings.Join(lines, indent)
}

func repeat(num int, text string) string {
	return strings.Repeat(text, num)
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

{{ .Description | indent 4 }}

{{ with .Example }}{{ . | indent 4 }} {{ end }}

{{ end }}
`
