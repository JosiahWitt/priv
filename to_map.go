package priv

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/joeycumines/go-dotnotation/dotnotation"
)

// ToMap converts the item (struct) to a map.
// Only the fields in the fieldsToInclude list are included.
// This list supports the dot operator to access nested fields.
//
// eg.
//  user := {ID: "123", Some: {Nested: {Field: "abc"}}, SomethingElse: true}
//  ToMap(user, "ID", "Some.Nested.Field->Renamed.Location")
// would result in:
//  {ID: "123", Renamed: {Location: "abc"}}
func ToMap(item interface{}, fieldsToInclude ...string) interface{} {
	reflectItem := reflect.ValueOf(item)
	return toMap(reflectItem, fieldsToInclude...)
}

func toMap(item reflect.Value, fieldsToInclude ...string) interface{} {
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
			panic(err)
		}

		// Set the value in the results
		err = resultAccessor.Set(result, field, value)
		if err != nil {
			panic(err)
		}
	}

	return result
}
