package mapstruct

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTo(t *testing.T) {
	t.Run("TestTo", func(t *testing.T) {
		type Data struct {
			Name string `json:"name"`
		}

		type TestStruct struct {
			Field1 string
			Field2 int
			Field3 []*Data `json:"field3"`
		}

		m := map[string]any{
			"Field1": "value1",
			"Field2": 42,
			"field3": []interface{}{
				map[string]any{"name": "test1"},
				map[string]any{"name": "test2"},
			},
		}

		expect := &TestStruct{
			Field1: "value1",
			Field2: 42,
			Field3: []*Data{
				{Name: "test1"},
				{Name: "test2"},
			},
		}

		var s TestStruct
		err := To(m, &s)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if diff := cmp.Diff(s, *expect); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}
