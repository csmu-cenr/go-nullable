package nullable

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func hasField(typ reflect.Type, fieldName string) bool {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			return true
		}
	}
	return false
}

func fieldNameIsInFields(fieldName string, fields []string) bool {
	for i := 0; i < len(fields); i++ {
		if fieldName == fields[i] {
			return true
		}
	}
	return false
}

func GetSelectedFieldsSlice(slice interface{}, fields []string) []map[string]interface{} {
	var results []map[string]interface{}

	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return results
	}

	results = make([]map[string]interface{}, 0)
	for i := 0; i < sliceValue.Len(); i++ {
		element := sliceValue.Index(i).Interface()
		selectedFields := GetSelectedFields(element, fields)
		results = append(results, selectedFields)
	}

	return results
}

func GetSelectedFields(v interface{}, fields []string) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name
		jsonTag := fieldType.Tag.Get("json")
		if strings.Contains(jsonTag, ",") {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}
		if jsonTag == "" {
			jsonTag = fieldName
		}
		fieldNameInFields := fieldNameIsInFields(fieldName, fields)
		if !fieldNameInFields {
			fieldNameInFields = fieldNameIsInFields(jsonTag, fields)
		}
		if field.Kind() == reflect.Struct && hasField(fieldType.Type, "Selected") {
			selectedField := field.FieldByName("Selected")
			if (selectedField.IsValid() && selectedField.Bool()) || fieldNameInFields {
				result[jsonTag] = field.Interface()
			}
		} else if fieldNameInFields {
			result[jsonTag] = field.Interface()
		}
	}

	return result
}

// Nullable represents data that also can be NULL
type Nullable[T any] struct {
	Data     T    `json:"data,omitempty"`
	Valid    bool `json:"-"`
	Selected bool `json:"-"`
}

// Value Create a Nullable from a value
func Value[T any](value T) Nullable[T] {
	if any(value) == nil {
		return Nullable[T]{Valid: false}
	}
	return Nullable[T]{Data: value, Valid: true}
}

// ValueFromPointer Create a Nullable from a pointer
func ValueFromPointer[T any](value *T) Nullable[T] {
	if value == nil {
		return Nullable[T]{Valid: false}
	}
	return Value(*value)
}

// Null Create a Nullable that is NULL with type
func Null[T any]() Nullable[T] {
	return Nullable[T]{}
}

// ValueOrZero Get Value, or default zero value if it is NULL
func (n Nullable[T]) ValueOrZero() T {
	if !n.Valid {
		var ref T
		return ref
	}
	return n.Data
}

func (n Nullable[T]) IsZero() bool {
	if !n.Valid {
		return true
	}
	var ref T
	return any(ref) == any(n.Data)
}

// Equal Check if this Nullable is equal to another Nullable
func (n Nullable[T]) Equal(other Nullable[T]) bool {
	switch any(n.Data).(type) {
	case time.Time:
		nValue := any(n.Data).(time.Time)
		otherValue := any(other.Data).(time.Time)
		return n.Valid == other.Valid && (!n.Valid || nValue.Equal(otherValue))
	}
	return n.ExactEqual(other)
}

// ExactEqual Check if this Nullable is exact equal to another Nullable, never using intern Equal method to check equality
func (n Nullable[T]) ExactEqual(other Nullable[T]) bool {
	return n.Valid == other.Valid && (!n.Valid || any(n.Data) == any(other.Data))
}

// String Convert value to string
func (n Nullable[T]) String() string {
	return fmt.Sprintf("%s", any(n.Data))
}

func (n Nullable[T]) GoString() string {
	var ref T
	return fmt.Sprintf("nullable.Nullable[%T]{Data:%#v,Valid:%#v}", ref, n.Data, n.Valid)
}
