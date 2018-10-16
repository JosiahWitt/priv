package priv

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/joeycumines/go-dotnotation/dotnotation"
)

// ToMap converts the item (struct or array of structs) to a map.
// Only the fields in the fieldsToInclude list are included.
// This list supports the dot operator to access nested fields.
// This function panics if the field cannot be found or set.
// If you are calling this function with a dynamically built list of fieldsToInclude, you may want to use the priv.ToMapErr alternative.
//
// eg.
//  users := [{ID: "123", Some: {Nested: {Field: "abc"}}, SomethingElse: true}]
//  ToMap(users, "ID", "Some.Nested.Field")
// would result in:
//  [{ID: "123", Some: {Nested: {Field: "abc"}}]
func ToMap(item interface{}, fieldsToInclude ...string) interface{} {
	value, err := ToMapErr(item, fieldsToInclude...)
	if err != nil {
		panic("priv.ToMap: " + err.Error())
	}

	return value
}

// ToMapErr does the same thing as priv.ToMap, but returns an error if the field cannot be found or set.
// Most of the time panicing (as priv.ToMap does) is fine, since the fields are set before compiling, so a panic to warn of a typo is not a big deal.
func ToMapErr(item interface{}, fieldsToInclude ...string) (interface{}, error) {
	reflectItem := reflect.ValueOf(item)

	// If it is an array or slice
	if reflectItem.Kind() == reflect.Slice || reflectItem.Kind() == reflect.Array {
		result := []interface{}{}

		// Iterate over it, to convert each item to a map
		for i := 0; i < reflectItem.Len(); i++ {
			value, err := toMap(reflectItem.Index(i), fieldsToInclude...)
			if err != nil {
				return nil, err
			}

			result = append(result, value)
		}

		return result, nil
	}

	// Otherwise, convert the item to a map
	return toMap(reflectItem, fieldsToInclude...)
}

func toMap(item reflect.Value, fieldsToInclude ...string) (interface{}, error) {
	// Convert the struct to a map
	sMap := structs.Map(item.Interface())

	sMapAccessor := dotnotation.Accessor{}
	resultAccessor := dotnotation.Accessor{
		Getter: func(target interface{}, property string) (interface{}, error) {
			// Getter for result automatically adds another layer to the nested maps if required
			rMap, _ := target.(map[string]interface{})
			if rMap[property] == nil {
				rMap[property] = make(map[string]interface{})
			}

			return rMap[property], nil
		},
	}

	// Build the result
	result := make(map[string]interface{})
	for _, field := range fieldsToInclude {
		// Grab the value
		value, err := sMapAccessor.Get(sMap, field)
		if err != nil {
			return nil, err
		}

		// Set the value in the results
		if err := resultAccessor.Set(result, field, value); err != nil {
			return nil, err
		}
	}

	return result, nil
}
