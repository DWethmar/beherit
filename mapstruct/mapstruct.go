package mapstruct

import (
	"errors"
	"fmt"
	"reflect"
)

func To(m map[string]any, s any) error {
	if m == nil {
		return errors.New("map cannot be nil: provided map is nil")
	}

	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("invalid argument: expected pointer to a struct but got %v", val.Kind())
	}

	return mapToStruct(m, val.Elem())
}

func mapToStruct(m map[string]any, valElem reflect.Value) error {
	structType := valElem.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		mapVal, ok := m[key]
		if !ok || mapVal == nil {
			continue
		}

		structField := valElem.Field(i)
		if !structField.CanSet() {
			continue
		}

		if err := setValue(structField, mapVal); err != nil {
			return fmt.Errorf("error setting field '%s': %w", key, err)
		}
	}
	return nil
}

func setValue(field reflect.Value, val any) error {
	fieldType := field.Type()

	switch fieldType.Kind() {
	case reflect.Ptr:
		ptrValue := reflect.New(fieldType.Elem())
		if reflect.TypeOf(val).Kind() == reflect.Ptr && reflect.TypeOf(val).AssignableTo(fieldType) {
			field.Set(reflect.ValueOf(val))
			return nil
		}
		if reflect.TypeOf(val).Kind() == reflect.Map {
			if err := mapToStruct(val.(map[string]any), ptrValue.Elem()); err != nil {
				return fmt.Errorf("failed to set pointer struct: %w", err)
			}
			field.Set(ptrValue)
		}
	case reflect.Struct:
		if reflect.TypeOf(val).Kind() == reflect.Map {
			return mapToStruct(val.(map[string]any), field)
		} else if reflect.TypeOf(val).AssignableTo(fieldType) {
			field.Set(reflect.ValueOf(val))
		}
	case reflect.Slice:
		if reflect.TypeOf(val).AssignableTo(fieldType) {
			field.Set(reflect.ValueOf(val))
			return nil
		}
		sliceVal, ok := val.([]interface{})
		if !ok {
			return fmt.Errorf("expected slice but got %T", val)
		}
		slice := reflect.MakeSlice(fieldType, len(sliceVal), len(sliceVal))
		for i, item := range sliceVal {
			if err := setValue(slice.Index(i), item); err != nil {
				return fmt.Errorf("error setting slice index %d: %w", i, err)
			}
		}
		field.Set(slice)
	case reflect.String:
		if strVal, ok := val.(string); ok {
			field.SetString(strVal)
		} else {
			return fmt.Errorf("expected string but got %T", val)
		}
	default:
		if reflect.TypeOf(val).AssignableTo(fieldType) {
			field.Set(reflect.ValueOf(val))
		} else {
			return fmt.Errorf("cannot set field '%s' of type '%s' with value of type '%T'", field.Type(), fieldType, val)
		}
	}
	return nil
}
