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
			return nil
		}
	case reflect.Struct:
		if reflect.TypeOf(val).Kind() == reflect.Map {
			return mapToStruct(val.(map[string]any), field)
		} else if reflect.TypeOf(val).AssignableTo(fieldType) {
			field.Set(reflect.ValueOf(val))
			return nil
		}
	case reflect.Slice:
		valValue := reflect.ValueOf(val)
		if valValue.Kind() != reflect.Slice {
			return fmt.Errorf("expected slice but got %T", val)
		}

		slice := reflect.MakeSlice(fieldType, valValue.Len(), valValue.Len())
		for i := 0; i < valValue.Len(); i++ {
			item := valValue.Index(i).Interface()
			if err := setValue(slice.Index(i), item); err != nil {
				return fmt.Errorf("error setting slice index %d: %w", i, err)
			}
		}
		field.Set(slice)
		return nil
	case reflect.String:
		strVal, ok := val.(string)
		if !ok {
			return fmt.Errorf("expected string but got %T", val)
		}
		field.SetString(strVal)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setNumericValue(field, val, true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setNumericValue(field, val, false)
	default:
		// Handle custom numeric types explicitly:
		underlyingKind := fieldType.Kind()
		if underlyingKind == reflect.Uint || underlyingKind == reflect.Int {
			return setNumericValue(field, val, underlyingKind == reflect.Uint)
		}
		if reflect.TypeOf(val).AssignableTo(fieldType) {
			field.Set(reflect.ValueOf(val))
			return nil
		}
		return fmt.Errorf("cannot set field '%s' of type '%s' with value of type '%T'", field.Type(), fieldType, val)
	}
	return nil
}

func setNumericValue(field reflect.Value, val any, isUnsigned bool) error {
	v := reflect.ValueOf(val)

	var num uint64
	var numSigned int64
	var isFloat bool

	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		num = uint64(v.Float())
		numSigned = int64(v.Float())
		isFloat = true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num = uint64(v.Int())
		numSigned = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num = v.Uint()
		numSigned = int64(v.Uint())
	default:
		return fmt.Errorf("expected numeric type but got %T", val)
	}

	if isUnsigned {
		field.SetUint(num)
	} else {
		field.SetInt(numSigned)
	}

	if isFloat && float64(numSigned) != v.Float() {
		return fmt.Errorf("lossy conversion from float to int")
	}

	return nil
}
