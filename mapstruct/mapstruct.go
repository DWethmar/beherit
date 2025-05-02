package mapstruct

import (
	"errors"
	"fmt"
	"reflect"
)

// To converts a map to a struct. The map keys should match the struct field names or their JSON tags.
func To(m map[string]any, s any) error {
	if m == nil {
		return errors.New("map cannot be nil")
	}

	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	return mapToStruct(m, val.Elem())
}

func mapToStruct(m map[string]any, valElem reflect.Value) error {
	structType := valElem.Type()

	for i := range structType.NumField() {
		field := structType.Field(i)

		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		mapVal, ok := m[key]
		if !ok || mapVal == nil {
			continue // Skip if the key is not present in the map or the value is nil
		}

		structField := valElem.Field(i)
		if !structField.CanSet() {
			return fmt.Errorf("cannot set field '%s' of type '%s'", field.Name, structField.Type())
		}

		if err := setValue(structField, mapVal); err != nil {
			return fmt.Errorf("error setting field '%s' of type '%s': %w", field.Name, structField.Type(), err)
		}
	}

	return nil
}

func setValue(field reflect.Value, value any) error {
	fieldType := field.Type()
	valueReflect := reflect.ValueOf(value)

	switch fieldType.Kind() {
	case reflect.Ptr:
		if valueReflect.Kind() == reflect.Map {
			ptr := reflect.New(fieldType.Elem())
			if err := mapToStruct(value.(map[string]any), ptr.Elem()); err != nil {
				return err
			}
			field.Set(ptr)
			return nil
		} else if valueReflect.Type().AssignableTo(fieldType) {
			field.Set(valueReflect)
			return nil
		} else if valueReflect.Type().ConvertibleTo(fieldType) {
			field.Set(valueReflect.Convert(fieldType))
			return nil
		}
		return fmt.Errorf("cannot assign to pointer type '%s' from type '%s'", fieldType, valueReflect.Type())

	case reflect.Struct:
		if valueReflect.Kind() == reflect.Map {
			return mapToStruct(value.(map[string]any), field)
		}
		return fmt.Errorf("expected map for struct but got '%s'", valueReflect.Kind())

	case reflect.Slice:
		if valueReflect.Kind() != reflect.Slice {
			return fmt.Errorf("expected slice but got '%s'", valueReflect.Kind())
		}

		slice := reflect.MakeSlice(fieldType, valueReflect.Len(), valueReflect.Len())

		for i := 0; i < valueReflect.Len(); i++ {
			elemVal := valueReflect.Index(i).Interface()
			elemField := slice.Index(i)

			if err := setValue(elemField, elemVal); err != nil {
				return fmt.Errorf("slice index %d: %w", i, err)
			}
		}

		field.Set(slice)
		return nil

	default:
		if valueReflect.Type().AssignableTo(fieldType) {
			field.Set(valueReflect)
			return nil
		} else if valueReflect.Type().ConvertibleTo(fieldType) {
			field.Set(valueReflect.Convert(fieldType))
			return nil
		}
		return fmt.Errorf("cannot assign value of type '%s' to field of type '%s'", valueReflect.Type(), fieldType)
	}
}
