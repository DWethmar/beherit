package mapstruct

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTo(t *testing.T) {
	t.Run("TestTo", func(t *testing.T) {
		type Data struct {
			Name    string  `json:"name"`
			SubData []*Data `json:"sub_data"`
		}

		type TestStruct struct {
			TypeInt       int
			TypeUint      uint
			TypeInt8      int8
			TypeInt16     int16
			TypeInt32     int32
			TypeInt64     int64
			TypeUint8     uint8
			TypeUint16    uint16
			TypeUint32    uint32
			TypeUint64    uint64
			TypeByte      byte
			TypeRune      rune
			TypeStr       string
			TypeBool      bool
			TypeFloat     float64
			TypeFloat32   float32
			TypeFloat64   float64
			TypeAny       any
			TypeSlice     []string
			TypeMap       map[string]any
			TypeStruct    Data
			TypeStructPtr *Data
		}

		m := map[string]any{
			"TypeInt":     1,
			"TypeUint":    uint(2),
			"TypeInt8":    int8(3),
			"TypeInt16":   int16(4),
			"TypeInt32":   int32(5),
			"TypeInt64":   int64(6),
			"TypeUint8":   uint8(7),
			"TypeUint16":  uint16(8),
			"TypeUint32":  uint32(9),
			"TypeUint64":  uint64(10),
			"TypeByte":    byte(11),
			"TypeRune":    rune(12),
			"TypeStr":     "test",
			"TypeBool":    true,
			"TypeFloat":   1.23,
			"TypeFloat32": float32(4.56),
			"TypeFloat64": float64(7.89),
			"TypeAny":     "any",
			"TypeSlice":   []string{"a", "b", "c"},
			"TypeMap": map[string]any{
				"key1": "value1",
				"key2": 42,
			},
			"TypeStruct": Data{
				Name: "test",
				SubData: []*Data{
					{Name: "sub_test1"},
					{Name: "sub_test2"},
				},
			},
			"TypeStructPtr": &Data{
				Name: "test_ptr",
				SubData: []*Data{
					{Name: "sub_test_ptr1"},
					{Name: "sub_test_ptr2"},
				},
			},
		}

		expect := &TestStruct{
			TypeInt:     1,
			TypeUint:    uint(2),
			TypeInt8:    int8(3),
			TypeInt16:   int16(4),
			TypeInt32:   int32(5),
			TypeInt64:   int64(6),
			TypeUint8:   uint8(7),
			TypeUint16:  uint16(8),
			TypeUint32:  uint32(9),
			TypeUint64:  uint64(10),
			TypeByte:    byte(11),
			TypeRune:    rune(12),
			TypeStr:     "test",
			TypeBool:    true,
			TypeFloat:   1.23,
			TypeFloat32: float32(4.56),
			TypeFloat64: float64(7.89),
			TypeAny:     "any",
			TypeSlice:   []string{"a", "b", "c"},
			TypeMap: map[string]any{
				"key1": "value1",
				"key2": 42,
			},
			TypeStruct: Data{
				Name: "test",
				SubData: []*Data{
					{Name: "sub_test1"},
					{Name: "sub_test2"},
				},
			},
			TypeStructPtr: &Data{
				Name: "test_ptr",
				SubData: []*Data{
					{Name: "sub_test_ptr1"},
					{Name: "sub_test_ptr2"},
				},
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
