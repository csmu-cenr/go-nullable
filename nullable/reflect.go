package nullable

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// // findFieldByJSONTag finds a struct field by its JSON tag.
// func findFieldByJSONTag(dataType reflect.Type, jsonTag string) (reflect.StructField, bool) {
// 	for i := 0; i < dataType.NumField(); i++ {
// 		field := dataType.Field(i)
// 		tag := strings.Split(field.Tag.Get("json"), ",")[0]
// 		if tag == jsonTag {
// 			return field, true
// 		}
// 	}
// 	return reflect.StructField{}, false
// }

// Changes a struct with Nullable[T] by returning a map[string] from fields specified
// func StructToMap(data any, fields []string) (map[string]any, error) {

// 	function := `StructToMap`

// 	result := map[string]any{}

// 	if data == nil {
// 		m := ErrorMessage{
// 			Details:  `data cannot be nil`,
// 			ErrorNo:  http.StatusBadRequest,
// 			Function: function,
// 			Message:  BAD_REQUEST,
// 		}
// 		return result, m
// 	}

// 	// Dereference pointers if necessary
// 	check := reflect.ValueOf(data)
// 	if check.Kind() == reflect.Ptr && !check.IsNil() {
// 		check = check.Elem()
// 		data = check
// 	}

// 	// Get the type and value of the input data
// 	dataType := reflect.TypeOf(data)
// 	dataValue := reflect.ValueOf(data)

// 	// Ensure the input is a struct
// 	if dataType.Kind() != reflect.Struct {
// 		return nil, fmt.Errorf("input is not a struct")
// 	}

// 	count := dataType.NumField()
// 	for i := 0; i < count; i++ {
// 		field := dataType.Field(i)
// 		if field.PkgPath == "" {
// 			tag := strings.Split(field.Tag.Get(JSON), ",")[0]
// 			if tag == "" {
// 				tag = field.Name
// 			} else {
// 				tag = strings.Split(tag, COMMA)[0]
// 			}
// 			if stringInStrings(tag, fields) {
// 				fieldValue := dataValue.Field(i).Interface()
// 				result[tag] = fieldValue
// 			}
// 		}
// 	}
// 	return result, nil

// }

// func StructToMap(data any, fields []string) (map[string]any, error) {
// 	function := `StructToMap`

// 	result := map[string]any{}

// 	if data == nil {
// 		m := ErrorMessage{
// 			Details:  `data cannot be nil`,
// 			ErrorNo:  http.StatusBadRequest,
// 			Function: function,
// 			Message:  BAD_REQUEST,
// 		}
// 		return result, m
// 	}

// 	// Dereference pointers if necessary
// 	v := reflect.ValueOf(data)
// 	if v.Kind() == reflect.Ptr {
// 		if v.IsNil() {
// 			return nil, fmt.Errorf("input pointer is nil")
// 		}
// 		v = v.Elem()
// 	}

// 	t := v.Type()

// 	// Ensure the input is a struct
// 	if t.Kind() != reflect.Struct {
// 		return nil, fmt.Errorf("input is not a struct")
// 	}

// 	count := t.NumField()
// 	for i := 0; i < count; i++ {
// 		field := t.Field(i)
// 		if field.PkgPath == "" { // exported
// 			tag := strings.Split(field.Tag.Get(JSON), ",")[0]
// 			if tag == "" {
// 				tag = field.Name
// 			} else {
// 				tag = strings.Split(tag, COMMA)[0]
// 			}

// 			if stringInStrings(tag, fields) {
// 				fieldValue := v.Field(i).Interface()
// 				result[tag] = fieldValue
// 			}
// 		}
// 	}

// 	return result, nil
// }

func StructToMap(data any, fields []string) (map[string]any, error) {
	const function = "StructToMap"

	result := make(map[string]any)

	if data == nil {
		m := ErrorMessage{
			Details:  "data cannot be nil",
			ErrorNo:  http.StatusBadRequest,
			Function: function,
			Message:  BAD_REQUEST,
		}
		return result, m
	}

	v := reflect.ValueOf(data)
	// Dereference pointers until we reach a non-pointer
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("input pointer is nil")
		}
		v = v.Elem()
	}
	t := v.Type()

	// ─────────────────────────────────────
	// Edge case: non-struct (primitive, slice, map, etc.)
	// ─────────────────────────────────────
	if t.Kind() != reflect.Struct {
		// There is no JSON tag here, so we can only use the provided field names as keys.
		if len(fields) == 0 {
			// No fields requested → nothing to map
			return result, nil
		}
		for _, name := range fields {
			name = strings.TrimSpace(name)
			name = strings.Trim(name, `"'`+"`")
			if name == "" {
				continue
			}
			result[name] = v.Interface()
		}
		return result, nil
	}

	// ─────────────────────────────────────
	// Normal case: struct
	// ─────────────────────────────────────
	n := t.NumField()
	for i := 0; i < n; i++ {
		sf := t.Field(i)

		// Skip unexported fields
		if sf.PkgPath != "" {
			continue
		}

		// Extract JSON tag (or fall back to field name)
		tag := sf.Tag.Get("json")
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = sf.Name
		} else {
			tag = strings.Split(tag, ",")[0] // strip ",omitempty" etc.
		}
		if tag == "" {
			continue
		}
		// Only map if this tag is requested
		if !FieldsContainsName(fields, tag) {
			continue
		}

		fieldValue := v.Field(i).Interface()
		result[tag] = fieldValue
	}

	return result, nil
}

// StructSetToMapSet converts a set of structs to a set of maps with selected fields.
func StructSetToMapSet(data any, fields []string) ([]map[string]any, error) {
	var result []map[string]any

	// Get the type and value of the input data
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	// Ensure the input is a slice
	if dataType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a slice")
	}

	// Iterate over the elements of the slice
	for i := 0; i < dataValue.Len(); i++ {
		element := dataValue.Index(i).Interface()

		// Convert the struct to a map with selected fields
		elementMap, err := StructToMap(element, fields)
		if err != nil {
			return nil, fmt.Errorf("error converting struct to map: %v", err)
		}

		// Append the map to the result slice
		result = append(result, elementMap)
	}

	return result, nil
}
