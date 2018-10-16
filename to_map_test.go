package priv_test

import (
	"testing"

	"github.com/JosiahWitt/priv"
	"github.com/stretchr/testify/require"
)

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

var testStruct1 = TestStruct1{
	Field1: "Value 1",
	Field2: InternalStruct1{
		Field1: true,
		Field2: &DoubleNestedStruct1{
			Field1: 42,
			Field2: 150.25,
		},
	},
}

var testStruct2 = TestStruct1{
	Field1: "Value 2",
	Field2: InternalStruct1{
		Field1: true,
		Field2: &DoubleNestedStruct1{
			Field1: 45,
			Field2: 243.42,
		},
	},
}

var testStructSlice = []TestStruct1{testStruct1, testStruct2}
var testStructArray = [2]TestStruct1{testStruct1, testStruct2}

var cases = []struct {
	when             string
	item             interface{}
	expected         interface{}
	fieldList        []string
	shouldPanicOrErr bool
	isArrayOrSlice   bool
}{
	{
		when:      "no fields are included",
		fieldList: []string{},
		item:      testStruct1,
		expected:  map[string]interface{}{},
	},
	{
		when:      "one field is included",
		fieldList: []string{"Field1"},
		item:      testStruct1,
		expected:  map[string]interface{}{"Field1": testStruct1.Field1},
	},
	{
		when:      "a complete doubly nested field is included",
		fieldList: []string{"Field2"},
		item:      testStruct1,
		expected: map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field1": testStruct1.Field2.Field1,
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
	},
	{
		when:      "all fields are included",
		fieldList: []string{"Field2", "Field1"},
		item:      testStruct1,
		expected: map[string]interface{}{
			"Field1": testStruct1.Field1,
			"Field2": map[string]interface{}{
				"Field1": testStruct1.Field2.Field1,
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
	},
	{
		when:      "a nested field is included",
		fieldList: []string{"Field2.Field1"},
		item:      testStruct1,
		expected: map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field1": testStruct1.Field2.Field1,
			},
		},
	},
	{
		when:      "all of a doubly nested field is included",
		fieldList: []string{"Field2.Field2"},
		item:      testStruct1,
		expected: map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
	},
	{
		when:      "the first field is included with all of a doubly nested field is included",
		fieldList: []string{"Field1", "Field2.Field2"},
		item:      testStruct1,
		expected: map[string]interface{}{
			"Field1": testStruct1.Field1,
			"Field2": map[string]interface{}{
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
					"Field2": testStruct1.Field2.Field2.Field2,
				},
			},
		},
	},
	{
		when:      "one field of a doubly nested field is included",
		fieldList: []string{"Field2.Field2.Field1"},
		item:      testStruct1,
		expected: map[string]interface{}{
			"Field2": map[string]interface{}{
				"Field2": map[string]interface{}{
					"Field1": testStruct1.Field2.Field2.Field1,
				},
			},
		},
	},
	{
		when:             "a non-nested field does not exist",
		fieldList:        []string{"DNE"},
		item:             testStruct1,
		expected:         nil,
		shouldPanicOrErr: true,
	},
	{
		when:             "a nested field does not exist",
		fieldList:        []string{"Field2.Field2z.Field1"},
		item:             testStruct1,
		expected:         nil,
		shouldPanicOrErr: true,
	},
	{
		when:           "given a slice, and the first field is included with all of a doubly nested field is included",
		fieldList:      []string{"Field1", "Field2.Field2"},
		item:           testStructSlice,
		isArrayOrSlice: true,
		expected: []map[string]interface{}{
			{
				"Field1": testStruct1.Field1,
				"Field2": map[string]interface{}{
					"Field2": map[string]interface{}{
						"Field1": testStruct1.Field2.Field2.Field1,
						"Field2": testStruct1.Field2.Field2.Field2,
					},
				},
			},
			{
				"Field1": testStruct2.Field1,
				"Field2": map[string]interface{}{
					"Field2": map[string]interface{}{
						"Field1": testStruct2.Field2.Field2.Field1,
						"Field2": testStruct2.Field2.Field2.Field2,
					},
				},
			},
		},
	},
	{
		when:           "given an array, and the first field is included with all of a doubly nested field is included",
		fieldList:      []string{"Field1", "Field2.Field2"},
		item:           testStructArray,
		isArrayOrSlice: true,
		expected: []map[string]interface{}{
			{
				"Field1": testStruct1.Field1,
				"Field2": map[string]interface{}{
					"Field2": map[string]interface{}{
						"Field1": testStruct1.Field2.Field2.Field1,
						"Field2": testStruct1.Field2.Field2.Field2,
					},
				},
			},
			{
				"Field1": testStruct2.Field1,
				"Field2": map[string]interface{}{
					"Field2": map[string]interface{}{
						"Field1": testStruct2.Field2.Field2.Field1,
						"Field2": testStruct2.Field2.Field2.Field2,
					},
				},
			},
		},
	},
	{
		when:             "given a slice, and a nested field does not exist",
		fieldList:        []string{"Field2.Field2z.Field1"},
		item:             testStructSlice,
		expected:         nil,
		shouldPanicOrErr: true,
	},
	{
		when:             "given an array, and a nested field does not exist",
		fieldList:        []string{"Field2.Field2z.Field1"},
		item:             testStructArray,
		expected:         nil,
		shouldPanicOrErr: true,
	},
}

func TestToMap(t *testing.T) {
	r := require.New(t)

	for _, c := range cases {
		if c.shouldPanicOrErr {
			defer func() {
				r.NotNil(recover(), "When "+c.when+", it should panic")
			}()
		}

		if c.isArrayOrSlice {
			r.ElementsMatch(c.expected, priv.ToMap(c.item, c.fieldList...), "When "+c.when)
		} else {
			r.Equal(c.expected, priv.ToMap(c.item, c.fieldList...), "When "+c.when)
		}
	}
}

func TestToMapErr(t *testing.T) {
	r := require.New(t)

	for _, c := range cases {
		v, err := priv.ToMapErr(c.item, c.fieldList...)

		if c.shouldPanicOrErr {
			r.Error(err, "When "+c.when+", it should error")
		} else {
			r.NoError(err, "When "+c.when+", it should not error")
		}

		if c.isArrayOrSlice {
			r.ElementsMatch(c.expected, v, "When "+c.when)
		} else {
			r.Equal(c.expected, v, "When "+c.when)
		}
	}
}
