package priv_test

import (
	"testing"

	"github.com/JosiahWitt/priv"
	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	r := require.New(t)

	type DoubleNestedStruct1 struct {
		Field1 int
		Field2 float64
	}

	type InternalStruct1 struct {
		Field1 bool
		Field2 *DoubleNestedStruct1
	}

	type TestStruct1 struct {
		Field1 string
		Field2 InternalStruct1
	}

	testStruct1 := TestStruct1{
		Field1: "Value 1",
		Field2: InternalStruct1{
			Field1: true,
			Field2: &DoubleNestedStruct1{
				Field1: 42,
				Field2: 150.25,
			},
		},
	}

	// When no fields are included
	r.Equal(
		map[string]interface{}{},
		priv.ToMap(testStruct1),
	)

	// When one field is included
	r.Equal(
		map[string]interface{}{"Field1": testStruct1.Field1},
		priv.ToMap(testStruct1, "Field1"),
	)

	// When a doubly nested field is included
	r.Equal(
		map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field1": testStruct1.Field2.Field1,
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
		priv.ToMap(testStruct1, "Field2"),
	)

	// When both fields are included
	r.Equal(
		map[string]interface{}{
			"Field1": testStruct1.Field1,
			"Field2": map[string]interface{}{
				"Field1": testStruct1.Field2.Field1,
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
		priv.ToMap(testStruct1, "Field2", "Field1"),
	)

	// When part of a doubly nested field is included
	r.Equal(
		map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field1": testStruct1.Field2.Field1,
			},
		},
		priv.ToMap(testStruct1, "Field2.Field1"),
	)

	// When part of a doubly nested field is included
	r.Equal(
		map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
		priv.ToMap(testStruct1, "Field2.Field2"),
	)

	// When part of a doubly nested field is included
	r.Equal(
		map[string]interface{}{
			"Field1": testStruct1.Field1,
			"Field2": map[string]interface{}{
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
		priv.ToMap(testStruct1, "Field1", "Field2.Field2"),
	)

	// When part of a doubly nested field is included
	r.Equal(
		map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
				},
			},
		},
		priv.ToMap(testStruct1, "Field2.Field2.Field1"),
	)
}
