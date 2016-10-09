package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
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

	return transformAtPath(NewPath(root), raw, targetValue)
}

// TODO: better code sharing with ValidateFields
func transformAtPath(path Path, raw map[string]interface{}, target reflect.Value) error {
	fields, err := structFields(target)
	if err != nil {
		return err
	}

	for _, field := range fields {
		value, ok := raw[field.tags.Name]
		if !ok {
			continue
		}
		delete(raw, field.tags.Name)

		localPath := path.add(field.tags.Name)
		rawValue := reflect.ValueOf(value)
		if err := transformField(localPath, rawValue, field.value); err != nil {
			return err
		}
	}

	for key := range raw {
		return PathErrorf(path.add(key), "unexpected key")
	}
	return nil
}

// structFields iterates over a struct and returns a list of all the fields
// included fields from embded types.
func structFields(target reflect.Value) ([]field, error) {
	fields := []field{}

	structType := target.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		if structField.Anonymous {
			embededFields, err := structFields(target.Field(i))
			if err != nil {
				return fields, err
			}
			fields = append(fields, embededFields...)
			continue
		}

		tags, err := NewFieldTags(structField.Name, structField.Tag.Get(StructTagKey))
		if err != nil {
			return fields, err
		}
		fields = append(fields, field{value: target.Field(i), tags: tags})
	}
	return fields, nil
}

type field struct {
	value reflect.Value
	tags  FieldTags
}

// FieldTags are annotations that specify properties of the config field
type FieldTags struct {
	IsRequired bool
	DoValidate bool
	Name       string
}

// NewFieldTags creates a FieldTags struct from a StructField.Tag string
func NewFieldTags(name, tags string) (FieldTags, error) {
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
			return field, fmt.Errorf("invalid field tag %q in %q", item, tags)
		}
	}
	if field.Name == "" {
		field.Name = TitleCaseToDash(name)
	}
	return field, nil
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

func transformField(path Path, raw reflect.Value, target reflect.Value) error {
	// TODO: handle pointer types
	if !target.CanSet() {
		return PathErrorf(path, "cant set target")
	}
	// Structs can be other types because they can use a TransformConfig method
	// to convert the raw type into their type.
	if target.Kind() != reflect.Struct && target.Kind() != raw.Kind() {
		return PathErrorf(path, "expected type %q not %q", target.Kind(), raw.Kind())
	}
	// TODO: recursive call for struct type
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

// TODO: share code with runValidationFunc
func transformStruct(path Path, raw reflect.Value, target reflect.Value) error {
	methodName := "TransformConfig"

	// Defer to the structs TransformConfig() method if it exists
	ptrTarget := target.Addr()
	method := ptrTarget.MethodByName(methodName)
	if !method.IsValid() {
		// TODO: Otherwise use transformAtPath (needs some refactor)
		//return transformAtPath(path, raw, target)
		return nil
	}

	switch transformFunc := method.Interface().(type) {
	case func(reflect.Value) error:
		if err := transformFunc(raw); err != nil {
			return PathErrorf(path, err.Error())
		}
	default:
		return fmt.Errorf(
			"%s.%s must be of type \"func(reflect.Value) error\" not %T",
			target.Type(), methodName, transformFunc)
	}
	return nil
}

func transformSlice(path Path, raw reflect.Value, target reflect.Value) error {
	elementType := target.Type().Elem()

	target.Set(reflect.MakeSlice(target.Type(), raw.Len(), raw.Len()))
	for i := 0; i < raw.Len(); i++ {
		item := raw.Index(i).Elem()

		if item.Kind() != elementType.Kind() {
			return PathErrorf(path.add(strconv.FormatInt(int64(i), 10)),
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
			return PathErrorf(path.add(key.String()),
				"item in the map is of wrong type %q, expected %q",
				item.Kind(), elementType)
		}
		target.SetMapIndex(key, item)
	}
	return nil
}
