package mapcopier_test

import (
	"testing"

	"github.com/dwethmar/beherit/mapcopier"
	"github.com/google/go-cmp/cmp"
)

func TestDeepCopy(t *testing.T) {
	t.Run("TestDeepCopy", func(t *testing.T) {
		// Test cases for deep copy
		tests := []struct {
			name     string
			input    map[string]any
			expected map[string]any
		}{
			{
				name: "Simple map",
				input: map[string]any{
					"key1": "value1",
					"key2": 42,
				},
				expected: map[string]any{
					"key1": "value1",
					"key2": 42,
				},
			},
			{
				name: "Nested map",
				input: map[string]any{
					"key1": map[string]any{
						"nestedKey1": "nestedValue1",
						"nestedKey2": 100,
					},
					"key2": 42,
				},
				expected: map[string]any{
					"key1": map[string]any{
						"nestedKey1": "nestedValue1",
						"nestedKey2": 100,
					},
					"key2": 42,
				},
			},
			{
				name: "Slice in map",
				input: map[string]any{
					"key1": []any{"value1", "value2"},
					"key2": 42,
				},
				expected: map[string]any{
					"key1": []any{"value1", "value2"},
					"key2": 42,
				},
			},
			{
				name: "Complex nested structure",
				input: map[string]any{
					"key1": []any{
						map[string]any{"subKey1": "subValue1"},
						map[string]any{"subKey2": 200},
					},
					"key2": 42,
				},
				expected: map[string]any{
					"key1": []any{
						map[string]any{"subKey1": "subValue1"},
						map[string]any{"subKey2": 200},
					},
					"key2": 42,
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := mapcopier.Copy(tt.input)
				if diff := cmp.Diff(tt.expected, got); diff != "" {
					t.Errorf("DeepCopy() mismatch (-want +got):\n%s", diff)
				}
			})
		}
	})
}
