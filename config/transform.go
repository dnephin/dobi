package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// TransformError is an error during  transforming a raw type to a structured
// type.
type TransformError struct {
	path Path
	msg  string
}

func (e *TransformError) Error() string {
	return fmt.Sprintf("Error at %s: %s", e.path.String(), e.msg)
}

// Path returns the config path where the error occurred
func (e *TransformError) Path() Path {
	return e.path
}

func tErrorf(path Path, msg string, args ...interface{}) *TransformError {
	return &TransformError{path: path, msg: fmt.Sprintf(msg, args...)}
}

// Path is a dotted path of key names to config values
type Path struct {
	path []string
}

func (p *Path) add(next string) Path {
	return Path{path: append(p.path, next)}
}

// Path returns the config keys in the path
func (p *Path) Path() []string {
	return p.path
}

func (p *Path) String() string {
	return strings.Join(p.path, ".")
}

// NewPath returns a new root Path
func NewPath(root string) Path {
	return Path{path: []string{root}}
}

// Transform recursively copies values from a raw map of values into the
// target structure. A TransformError is returned if the raw type has a field
// with an incorrect type, or has extra fields.
func Transform(root string, raw map[string]interface{}, target interface{}) error {
	targetValue := reflect.ValueOf(target).Elem()

	if kind := targetValue.Kind(); kind != reflect.Struct {
		return fmt.Errorf("invalid target type %s, must be a Struct", kind)
	}

	return transformAtPath(NewPath(root), raw, targetValue)
}

func transformAtPath(path Path, raw map[string]interface{}, target reflect.Value) error {
	for key, value := range raw {
		localPath := path.add(key)
		key := dashToTitleCase(key)
		field := target.FieldByName(key)

		if !field.IsValid() {
			return tErrorf(localPath, "unexpected key")
		}

		rawValue := reflect.ValueOf(value)
		if err := transformField(localPath, rawValue, field); err != nil {
			return err
		}
	}
	return nil
}

func dashToTitleCase(source string) string {
	return strings.Replace(strings.Title(source), "-", "", -1)
}

func titleCaseToDash(source string) string {
	var buff bytes.Buffer
	var prevCharIsLower bool

	for _, char := range source {
		var n rune
		switch {
		case unicode.IsLower(char):
			n = char
		case prevCharIsLower:
			buff.WriteRune('-')
			n = unicode.ToLower(char)
		default:
			n = unicode.ToLower(char)
		}
		buff.WriteRune(n)
		prevCharIsLower = unicode.IsLower(char)
	}
	return buff.String()
}

func transformField(path Path, raw reflect.Value, target reflect.Value) error {
	// TODO: handle pointer types
	if !target.CanSet() {
		return tErrorf(path, "cant set target")
	}
	if target.Kind() != raw.Kind() {
		return tErrorf(path, "expected type %q not %q", target.Kind(), raw.Kind())
	}
	// TODO: recursive call for struct type
	switch target.Kind() {
	case reflect.Slice:
		return transformSlice(path, raw, target)
	case reflect.Map:
		return transformMap(path, raw, target)
	default:
		target.Set(raw)
	}
	return nil
}

func transformSlice(path Path, raw reflect.Value, target reflect.Value) error {
	elementType := target.Type().Elem()

	target.Set(reflect.MakeSlice(target.Type(), raw.Len(), raw.Len()))
	for i := 0; i < raw.Len(); i++ {
		item := raw.Index(i).Elem()

		if item.Kind() != elementType.Kind() {
			return tErrorf(path.add(strconv.FormatInt(int64(i), 10)),
				"item in the list is of wrong type %q, expected %q",
				item.Kind(), elementType)
		}
		target.Index(i).Set(item)
	}
	return nil
}

func transformMap(path Path, raw reflect.Value, target reflect.Value) error {
	elementType := target.Type().Elem()

	target.Set(reflect.MakeMap(target.Type()))
	for _, key := range raw.MapKeys() {
		// TODO: how to I check keys Kind? against target key.Kind()
		key = key.Elem()
		item := raw.MapIndex(key).Elem()

		if item.Kind() != elementType.Kind() {
			return tErrorf(path.add(key.String()),
				"item in the map is of wrong type %q, expected %q",
				item.Kind(), elementType)
		}
		target.SetMapIndex(key, item)
	}
	return nil
}
