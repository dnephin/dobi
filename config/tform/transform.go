package tform

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	pth "github.com/dnephin/dobi/config/tform/path"
)

const (
	// StructTagKey is the key used to find struct tags on a field
	StructTagKey = "config"
)

// Transform recursively copies values from a raw map of values into the
// target structure. A PathError is returned if the raw type has a field
// with an incorrect type, or has extra fields.
func Transform(root string, raw map[string]interface{}, target interface{}) error {
	targetValue := reflect.ValueOf(target).Elem()

	if kind := targetValue.Kind(); kind != reflect.Struct {
		return fmt.Errorf("invalid target type %s, must be a Struct", kind)
	}

	return transformAtPath(pth.NewPath(root), raw, targetValue)
}

func transformAtPath(path pth.Path, raw map[string]interface{}, target reflect.Value) error {
	for _, field := range structFields(target) {
		value, ok := raw[field.tags.Name]
		if !ok {
			continue
		}
		delete(raw, field.tags.Name)

		localPath := path.Add(field.tags.Name)
		rawValue := reflect.ValueOf(value)
		if err := transformField(localPath, rawValue, field.value); err != nil {
			return err
		}
	}

	for key := range raw {
		return pth.Errorf(path.Add(key), "unexpected key")
	}
	return nil
}

// structFields iterates over a struct and returns a list of all the fields
// included fields from embded types.
func structFields(target reflect.Value) []field {
	fields := []field{}

	structType := target.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		if structField.Anonymous {
			fields = append(fields, structFields(target.Field(i))...)
			continue
		}

		tags := NewFieldTags(structField.Name, structField.Tag.Get(StructTagKey))
		fields = append(fields, field{
			value:       target.Field(i),
			structField: structField,
			tags:        tags,
		})
	}
	return fields
}

type field struct {
	value       reflect.Value
	structField reflect.StructField
	tags        FieldTags
}

// FieldTags are annotations that specify properties of the config field
type FieldTags struct {
	IsRequired bool
	DoValidate bool
	Name       string
}

// NewFieldTags creates a FieldTags struct from a StructField.Tag string
func NewFieldTags(name, tags string) FieldTags {
	field := FieldTags{}
	for index, item := range strings.Split(tags, ",") {
		switch {
		case item == "required":
			field.IsRequired = true
		case item == "validate":
			field.DoValidate = true
		case index == 0:
			field.Name = item
		default:
			panic(fmt.Errorf("invalid field tag %q in %q", item, tags))
		}
	}
	if field.Name == "" {
		field.Name = TitleCaseToDash(name)
	}
	return field
}

// TitleCaseToDash converts a CamelCased name into a dashed config field name
func TitleCaseToDash(source string) string {
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

func transformField(path pth.Path, raw reflect.Value, target reflect.Value) error {
	// TODO: handle pointer types
	if !target.CanSet() {
		return pth.Errorf(path, "cant set target")
	}
	// Structs can be other types because they can use a TransformConfig method
	// to convert the raw type into their type.
	if target.Kind() != reflect.Struct && target.Kind() != raw.Kind() {
		return pth.Errorf(path, "expected type %q not %q", target.Kind(), raw.Kind())
	}
	switch target.Kind() {
	case reflect.Slice:
		return transformSlice(path, raw, target)
	case reflect.Map:
		return transformMap(path, raw, target)
	case reflect.Struct:
		return transformStruct(path, raw, target)
	default:
		target.Set(raw)
	}
	return nil
}

func transformStruct(path pth.Path, raw reflect.Value, target reflect.Value) error {
	methodName := "TransformConfig"

	// Defer to the structs TransformConfig() method if it exists
	ptrTarget := target.Addr()
	method := ptrTarget.MethodByName(methodName)
	if !method.IsValid() {
		mapping := make(map[string]interface{})
		err := transformMap(path, raw, reflect.ValueOf(mapping))
		if err != nil {
			return err
		}
		return transformAtPath(path, mapping, target)
	}

	switch transformFunc := method.Interface().(type) {
	case func(reflect.Value) error:
		if err := transformFunc(raw); err != nil {
			return pth.Errorf(path, err.Error())
		}
	default:
		return fmt.Errorf(
			"%s.%s must be of type \"func(reflect.Value) error\" not %T",
			target.Type(), methodName, transformFunc)
	}
	return nil
}

func transformSlice(path pth.Path, raw reflect.Value, target reflect.Value) error {
	target.Set(reflect.MakeSlice(target.Type(), raw.Len(), raw.Len()))
	for i := 0; i < raw.Len(); i++ {
		err := transformField(
			path.Add(strconv.FormatInt(int64(i), 10)),
			raw.Index(i).Elem(),
			target.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

// https://github.com/go-yaml/yaml/blob/v2/decode.go#L539
func transformMap(path pth.Path, raw reflect.Value, target reflect.Value) error {
	targetType := target.Type()
	elementType := targetType.Elem()
	keyType := targetType.Key()

	if target.IsNil() {
		target.Set(reflect.MakeMap(targetType))
	}

	for _, key := range raw.MapKeys() {
		localPath := path.Add(key.String())
		item := raw.MapIndex(key)

		keyValue, err := castScalar(key, keyType)
		if err != nil {
			return pth.Errorf(localPath, err.Error())
		}

		itemValue, err := castScalar(item, elementType)
		if err != nil {
			return pth.Errorf(localPath, err.Error())
		}

		target.SetMapIndex(keyValue, itemValue)
	}
	return nil
}

func castScalar(raw reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	target := reflect.New(targetType).Elem()
	rawValue := raw.Interface()
	switch targetType.Kind() {
	case reflect.Interface:
		target.Set(raw)
	case reflect.String:
		value, ok := rawValue.(string)
		if !ok {
			return target, fmt.Errorf("expected string but got %T", rawValue)
		}
		target.SetString(value)
	default:
		// TODO: more scalar types
		panic(fmt.Sprintf("Not implemeneted yet: Kind %s", targetType.Kind()))
	}
	return target, nil
}
