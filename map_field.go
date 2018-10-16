package priv

import "strings"

type mapField struct {
	To   string // Where the field should be mapped to, if empty, use value in From
	From string // Where the field should be mapped from
}

// ToLocation fetches the location the field should be mapped to
func (f *mapField) ToLocation() string {
	if f.To == "" {
		return f.From
	}
	return f.To
}

// FromLocation fetches the location the field should be mapped from
func (f *mapField) FromLocation() string {
	return f.From
}

func mapFieldFromFieldString(fieldString string) (f mapField) {
	// Separate the key parts by the rename operator ->
	fieldParts := strings.Split(fieldString, "->")
	f.From = fieldParts[0]

	// Add the To location if present
	if len(fieldParts) > 1 {
		f.To = fieldParts[1]
	}

	return
}
