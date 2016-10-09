package config

import (
	"fmt"
	"reflect"
)

// ValidateFields runs validations as defined by struct tags
func ValidateFields(path Path, resource interface{}) error {
	structValue := reflect.ValueOf(resource)
	value := structValue.Elem()

	if kind := value.Kind(); kind != reflect.Struct {
		return fmt.Errorf("invalid target type %s, must be a Struct", kind)
	}

	for _, field := range structFields(value) {
		if err := validateField(path, structValue, field); err != nil {
			return err
		}
	}
	return nil
}

func validateField(path Path, structValue reflect.Value, field field) error {
	path = path.add(field.tags.Name)

	if field.tags.IsRequired {
		zero := reflect.Zero(field.value.Type()).Interface()
		if reflect.DeepEqual(field.value.Interface(), zero) {
			return PathErrorf(path, "a value is required")
		}
	}
	if field.tags.DoValidate {
		if err := runValidationFunc(path, structValue, field.structField.Name); err != nil {
			return err
		}
	}
	return nil
}

func runValidationFunc(path Path, structValue reflect.Value, field string) error {
	methodName := "Validate" + field
	methodValue, err := getMethodFromStruct(structValue, methodName)
	if err != nil {
		return err
	}

	switch validationFunc := methodValue.Interface().(type) {
	case func() error:
		if err := validationFunc(); err != nil {
			return PathErrorf(path, "failed validation: %s", err)
		}
		return nil
	default:
		return fmt.Errorf("%s.%s must be of type \"func() error\" not %T",
			structValue.Elem().Type(), methodName, validationFunc)
	}
}

func getMethodFromStruct(structValue reflect.Value, methodName string) (reflect.Value, error) {
	// First look for method with non-pointer receiver
	methodValue := structValue.Elem().MethodByName(methodName)
	if methodValue.IsValid() {
		return methodValue, nil
	}

	// Second look for method with pointer receiver
	methodValue = structValue.MethodByName(methodName)
	if methodValue.IsValid() {
		return methodValue, nil
	}

	return reflect.Value{}, fmt.Errorf("%s is missing validation function %q",
		structValue.Type(), methodName)
}
