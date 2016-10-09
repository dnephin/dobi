/*
Package docs generates user documentation from config structs.

The config struct is read using `reflect` and `go/ast` and translated into an
intermediate format.

*/
package docs

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dnephin/dobi/config"
)

// ConfigType holds the details about a configuration type
type ConfigType struct {
	Name        string
	Description string
	Example     string
	Fields      []ConfigField
}

// ConfigField holds the details for each field in the ConfigType
type ConfigField struct {
	Name        string
	IsRequired  bool
	Type        string
	Format      string
	Example     string
	Default     string
	Description string
}

// OutputFormat is an enumeration of supported output formats
type OutputFormat string

const (
	// ReStructuredText format is .rst
	ReStructuredText OutputFormat = "rst"
)

// Generate returns documentation for a config type in the requested format, or
// an error if the documentation could not be generated.
func Generate(config interface{}, format OutputFormat) (string, error) {
	configType, err := Parse(config)
	if err != nil {
		return "", err
	}

	switch format {
	case ReStructuredText:
		return FormatRst(configType)
	default:
		return "", fmt.Errorf("Unsupported format %q", format)
	}
}

// Parse parses the definition of the struct and returns the metadata about the
// config that it accepts.
func Parse(source interface{}) (ConfigType, error) {
	config := ConfigType{}
	structType := reflect.TypeOf(source)
	comments, err := getStructComments(structType.Name(), structType.PkgPath())
	if err != nil {
		return config, err
	}

	config.Name = getTypeName(structType.Name(), comments)
	config.Description = comments.comment.description
	config.Example = comments.comment.Get("example", "")
	config.Fields, err = buildConfigFields(structType, comments)
	return config, err
}

type structComments struct {
	comment parsedComment
	fields  map[string]parsedComment
}

type parsedComment struct {
	description string
	values      map[string]string
}

// Get returns the value at key, or the default value if key is not set
func (c parsedComment) Get(key, def string) string {
	if value, exists := c.values[key]; exists {
		return value
	}
	return def
}

// TODO: support multi-line examples by keeping track of the last field that was
// added to values, and adding the line to it.
func parseComment(name, comment string) parsedComment {
	lines := strings.Split(strings.TrimSpace(comment), "\n")

	parsed := parsedComment{values: make(map[string]string)}
	parsed.description = strings.TrimPrefix(lines[0], name+" ")

	var lastValue string
	for _, line := range lines[1:] {
		parts := strings.SplitN(line, ":", 2)
		switch {
		case len(parts) == 2 && isAnnotation(parts[0]):
			lastValue = strings.ToLower(parts[0])
			parsed.values[lastValue] = strings.TrimPrefix(parts[1], " ")
		case len(parsed.values) == 0:
			parsed.description += "\n" + line
		case lastValue != "":
			parsed.values[lastValue] += "\n" + line
		}
	}

	return parsed
}

func isAnnotation(key string) bool {
	switch strings.ToLower(key) {
	case "name", "type", "format", "example", "default":
		return true
	default:
		return false
	}
}

func getStructComments(name string, path string) (*structComments, error) {
	// TODO: better way to go from pkgPath to node?
	var err error
	var fullPath string
	pathSep := fmt.Sprintf("%c", os.PathListSeparator)
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), pathSep) {
		fullPath = filepath.Join(gopath, "src", path)
		if _, err = os.Stat(fullPath); err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	pkg, exists := pkgs[filepath.Base(path)]
	if !exists {
		return nil, fmt.Errorf("%q not found", path)
	}

	var typeSpec *ast.TypeSpec
	var comment *ast.CommentGroup
	for _, pkgFile := range pkg.Files {
		for _, decl := range pkgFile.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						if spec.Name.Name == name {
							comment = decl.Doc
							typeSpec = spec
							break
						}
					}
				}
			}
		}
	}

	// TODO: comment can be nil for embded types
	if comment == nil || typeSpec == nil {
		return nil, fmt.Errorf("%q not found in declarations", name)
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T for %q", typeSpec.Type, name)
	}

	comments := structComments{
		comment: parseComment(name, comment.Text()),
		fields:  make(map[string]parsedComment),
	}
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		fieldName := field.Names[0].Name
		comments.fields[fieldName] = parseComment(fieldName, field.Doc.Text())
	}
	return &comments, nil
}

func getTypeName(name string, comments *structComments) string {
	return comments.comment.Get("name", config.TitleCaseToDash(name))
}

func embededFields(structType reflect.Type) ([]ConfigField, error) {
	fields := []ConfigField{}
	comments, err := getStructComments(structType.Name(), structType.PkgPath())
	if err != nil {
		return fields, err
	}
	return buildConfigFields(structType, comments)
}

func buildConfigFields(
	structType reflect.Type,
	comments *structComments,
) ([]ConfigField, error) {
	fields := []ConfigField{}
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		if structField.Anonymous {
			embeded, err := embededFields(structField.Type)
			if err != nil {
				return fields, err
			}
			fields = append(fields, embeded...)
			continue
		}

		comment := comments.fields[structField.Name]
		fieldTags := config.NewFieldTags(
			structField.Name, structField.Tag.Get(config.StructTagKey))
		field := ConfigField{
			Name:        fieldTags.Name,
			IsRequired:  fieldTags.IsRequired,
			Type:        comment.Get("type", structField.Type.Name()),
			Format:      comment.Get("format", ""),
			Example:     comment.Get("example", ""),
			Default:     comment.Get("default", ""),
			Description: comment.description,
		}
		fields = append(fields, field)
	}
	return fields, nil
}
