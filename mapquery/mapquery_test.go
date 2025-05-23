package mapquery_test

import (
	"reflect"
	"testing"

	"github.com/dwethmar/beherit/mapquery"
	"github.com/google/go-cmp/cmp"
)

func TestSet(t *testing.T) {
	t.Run("set value in map", func(t *testing.T) {
		m := map[string]any{"key": "value"}
		err := mapquery.Set("key", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"key": "newValue"}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in nested map", func(t *testing.T) {
		m := map[string]any{"outer": map[string]any{"inner": "value"}}
		err := mapquery.Set("outer.inner", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"outer": map[string]any{"inner": "newValue"}}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in slice", func(t *testing.T) {
		m := map[string]any{"outer": map[string]any{"inner": []any{"value1", "value2"}}}
		err := mapquery.Set("outer.inner[1]", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"outer": map[string]any{"inner": []any{"value1", "newValue"}}}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in slice with wildcard", func(t *testing.T) {
		m := map[string]any{"outer": map[string]any{"inner": []any{"value1", "value2"}}}
		err := mapquery.Set("outer.inner[*]", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"outer": map[string]any{"inner": []any{"newValue", "newValue"}}}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in non-existent path", func(t *testing.T) {
		m := map[string]any{}
		err := mapquery.Set("outer.inner", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"outer": map[string]any{"inner": "newValue"}}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in slice with nested map", func(t *testing.T) {
		m := map[string]any{"outer": map[string]any{"inner": []any{
			map[string]any{"key": "value1"}, map[string]any{"key": "value2"},
		}}}
		err := mapquery.Set("outer.inner[1].key", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"outer": map[string]any{"inner": []any{
			map[string]any{"key": "value1"}, map[string]any{"key": "newValue"},
		}}}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in slice with wildcard and nested map", func(t *testing.T) {
		m := map[string]any{"outer": map[string]any{"inner": []any{
			map[string]any{"key": "value1"}, map[string]any{"key": "value2"},
		}}}
		err := mapquery.Set("outer.inner[*].key", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"outer": map[string]any{"inner": []any{
			map[string]any{"key": "newValue"}, map[string]any{"key": "newValue"},
		}}}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in multiple nested slices", func(t *testing.T) {
		m := map[string]any{"outer": map[string]any{"inner": []any{
			map[string]any{"key": []string{"value1"}},
			map[string]any{"key": []string{"value2"}},
		}}}
		err := mapquery.Set("outer.inner[*].key[*]", "newValue", &m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := map[string]any{"outer": map[string]any{"inner": []any{
			map[string]any{"key": []any{"newValue"}},
			map[string]any{"key": []any{"newValue"}},
		}}}
		if diff := cmp.Diff(expected, m); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in raw slice", func(t *testing.T) {
		in := []string{"value1", "value2"}
		// Use interface wrapper to avoid reflect.Set panic
		var anySlice any = in
		err := mapquery.Set("[1]", "newValue", &anySlice)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		converted := anySlice.([]interface{})
		result := make([]string, len(converted))
		for i, v := range converted {
			result[i] = v.(string)
		}
		expected := []string{"value1", "newValue"}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in raw slice with wildcard", func(t *testing.T) {
		in := []string{"value1", "value2"}
		var anySlice any = in
		err := mapquery.Set("[*]", "newValue", &anySlice)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		converted := anySlice.([]interface{})
		result := make([]string, len(converted))
		for i, v := range converted {
			result[i] = v.(string)
		}
		expected := []string{"newValue", "newValue"}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in nested slice", func(t *testing.T) {
		in := []any{[]any{"value1", []any{"value2"}}}
		err := mapquery.Set("[0][1][0]", "newValue", &in)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := []any{[]any{"value1", []any{"newValue"}}}
		if diff := cmp.Diff(expected, in); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})

	t.Run("set value in nested slice with wildcard", func(t *testing.T) {
		in := []any{
			[]any{
				"value1",
				[]any{"value2", "value3"},
			},
		}
		err := mapquery.Set("[*][*][*]", "newValue", &in)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := []any{[]any{"value1", []any{"newValue", "newValue"}}}
		if diff := cmp.Diff(expected, in); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("get value from map", func(t *testing.T) {
		m := map[string]any{"key": "value"}
		val, err := mapquery.Get("key", m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if val != "value" {
			t.Fatalf("expected value to be 'value', got %v", val)
		}
	})

	t.Run("get value from nested map", func(t *testing.T) {
		m := map[string]any{"outer": map[string]any{"inner": "value"}}
		val, err := mapquery.Get("outer.inner", m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if val != "value" {
			t.Fatalf("expected value to be 'value', got %v", val)
		}
	})

	t.Run("get value from slice", func(t *testing.T) {
		m := map[string]any{"outer": []any{"value1", "value2"}}
		val, err := mapquery.Get("outer[1]", m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if val != "value2" {
			t.Fatalf("expected value to be 'value2', got %v", val)
		}
	})

	t.Run("get value from slice with wildcard", func(t *testing.T) {
		m := map[string]any{"outer": []any{"value1", "value2"}}
		val, err := mapquery.Get("outer[*]", m)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := []any{"value1", "value2"}
		if !reflect.DeepEqual(val, expected) {
			t.Fatalf("expected value to be %v, got %v", expected, val)
		}
	})
}
