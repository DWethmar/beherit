package mapquery

// todo
// https://chatgpt.com/c/6826f258-00fc-800b-9625-4283f3bdf913

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	sliceIndexRegex = regexp.MustCompile(`^([a-zA-Z0-9_]+|\[\*|\d+\])$`)
)

// Set sets a value into the nested data structure at the given path.
// It supports maps, slices, and wildcards (e.g., [*]) and modifies the input in-place.
func Set(path string, value any, target any) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}

	segments := parsePath(path)

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || !rv.Elem().CanSet() {
		return errors.New("target must be a pointer to a settable value")
	}

	root := rv.Elem().Interface()
	newVal, err := setValue(segments, value, root)
	if err != nil {
		return err
	}

	newValVal := reflect.ValueOf(newVal)
	if newValVal.Type().AssignableTo(rv.Elem().Type()) {
		rv.Elem().Set(newValVal)
		return nil
	}

	if rv.Elem().Kind() == reflect.Slice && newValVal.Kind() == reflect.Slice {
		slice := rv.Elem()
		slice.SetLen(0)
		for i := 0; i < newValVal.Len(); i++ {
			slice = reflect.Append(slice, newValVal.Index(i))
		}
		rv.Elem().Set(slice)
		return nil
	}

	return fmt.Errorf("cannot assign value of type %T to target of type %T", newVal, rv.Elem().Interface())
}

// Get retrieves a value from a nested data structure at the given path.
func Get(path string, data any) (any, error) {
	segments := parsePath(path)
	var current any = data

	for _, seg := range segments {
		if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
			index := seg[1 : len(seg)-1]
			slice, ok := toAnySlice(current)
			if !ok {
				return nil, fmt.Errorf("expected slice at segment %q but got %T", seg, current)
			}
			if index == "*" {
				return slice, nil // simple wildcard handling: return full slice here
			}
			i, err := strconv.Atoi(index)
			if err != nil {
				return nil, fmt.Errorf("invalid slice index %q", seg)
			}
			if i < 0 || i >= len(slice) {
				return nil, fmt.Errorf("index %d out of bounds", i)
			}
			current = slice[i]
		} else {
			m, ok := current.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("expected map at segment %q but got %T", seg, current)
			}
			val, exists := m[seg]
			if !exists {
				return nil, fmt.Errorf("key %q not found", seg)
			}
			current = val
		}
	}
	return current, nil
}

func parsePath(path string) []string {
	parts := []string{}
	for _, segment := range strings.Split(path, ".") {
		for {
			if i := strings.Index(segment, "["); i != -1 {
				if i > 0 {
					parts = append(parts, segment[:i])
				}
				j := strings.Index(segment[i:], "]")
				if j == -1 {
					break
				}
				parts = append(parts, segment[i:i+j+1])
				segment = segment[i+j+1:]
				if segment == "" {
					break
				}
			} else {
				parts = append(parts, segment)
				break
			}
		}
	}
	return parts
}

func setValue(segments []string, value any, current any) (any, error) {
	if len(segments) == 0 {
		return value, nil
	}

	seg := segments[0]

	if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
		slice, ok := toAnySlice(current)
		if !ok {
			return nil, fmt.Errorf("expected slice at %q but got %T", seg, current)
		}

		index := seg[1 : len(seg)-1]
		if index == "*" {
			for i := range slice {
				if len(segments) == 1 {
					slice[i] = value
					continue
				}
				if newVal, err := setValue(segments[1:], value, slice[i]); err == nil {
					slice[i] = newVal
				} else if len(segments[1:]) == 0 {
					slice[i] = value
				}
			}
			return slice, nil
		}

		i, err := strconv.Atoi(index)
		if err != nil {
			return nil, fmt.Errorf("invalid slice index %q", seg)
		}

		for len(slice) <= i {
			slice = append(slice, nil)
		}

		newVal, err := setValue(segments[1:], value, slice[i])
		if err != nil {
			return nil, err
		}
		slice[i] = newVal
		return slice, nil
	}

	var m map[string]any
	if current == nil {
		m = make(map[string]any)
	} else {
		var ok bool
		m, ok = current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map at segment %q but got %T", seg, current)
		}
	}

	child, exists := m[seg]
	if !exists {
		child = nil
	}
	newVal, err := setValue(segments[1:], value, child)
	if err != nil {
		return nil, err
	}
	m[seg] = newVal
	return m, nil
}

func toAnySlice(input any) ([]any, bool) {
	rv := reflect.ValueOf(input)
	if rv.Kind() != reflect.Slice {
		return nil, false
	}

	out := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out[i] = rv.Index(i).Interface()
	}
	return out, true
}
