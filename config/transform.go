package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
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
	for i := 0; i < target.Type().NumField(); i++ {
		structField := target.Type().Field(i)
		tags, err := newFieldTags(structField.Name, structField.Tag.Get("config"))
		if err != nil {
			return err
		}

		value, ok := raw[tags.name]
		if !ok {
			continue
		}
		delete(raw, tags.name)

		localPath := path.add(tags.name)
		rawValue := reflect.ValueOf(value)
		if err := transformField(localPath, rawValue, target.Field(i)); err != nil {
			return err
		}
	}

	for key := range raw {
		return PathErrorf(path.add(key), "unexpected key")
	}
	return nil
}

type fieldTags struct {
	isRequired bool
	doValidate bool
	name       string
}

func newFieldTags(name, tags string) (fieldTags, error) {
	field := fieldTags{}
	for index, item := range strings.Split(tags, ",") {
		switch {
		case item == "required":
			field.isRequired = true
		case item == "validate":
			field.doValidate = true
		case index == 0:
			field.name = item
		default:
			return field, fmt.Errorf("invalid field tag %q in %q", item, tags)
		}
	}
	if field.name == "" {
		field.name = titleCaseToDash(name)
	}
	return field, nil
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
		return PathErrorf(path, "cant set target")
	}
	if target.Kind() != raw.Kind() {
		return PathErrorf(path, "expected type %q not %q", target.Kind(), raw.Kind())
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
