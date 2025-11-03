package nullable

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// findFieldByJSONTag finds a struct field by its JSON tag.
func findFieldByJSONTag(dataType reflect.Type, jsonTag string) (reflect.StructField, bool) {
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		if tag == jsonTag {
			return field, true
		}
	}
	return reflect.StructField{}, false
}

// Changes a struct with Nullable[T] by returning a map[string] from fields specified
func StructToMap(data any, fields []string) (map[string]any, error) {

	function := `dmdata.StructToMap`

	result := map[string]any{}

	if data == nil {
		m := ErrorMessage{
			Details:  `data cannot be nil`,
			ErrorNo:  http.StatusBadRequest,
			Function: function,
			Message:  BAD_REQUEST,
		}
		return result, m
	}

	// Dereference pointers if necessary
	check := reflect.ValueOf(data)
	if check.Kind() == reflect.Ptr && !check.IsNil() {
		check = check.Elem()
		data = check
	}

	// Get the type and value of the input data
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	// Ensure the input is a struct
	if dataType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	// If fields is empty, map all fields in the struct
	if len(fields) == 0 {
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			if field.PkgPath == "" {
				fieldValue := dataValue.Field(i).Interface()
				jsonTag := field.Tag.Get(JSON)
				if jsonTag == "" {
					jsonTag = field.Name
				} else {
					jsonTag = strings.Split(jsonTag, COMMA)[0]
				}
				result[jsonTag] = fieldValue
			}
		}
		return result, nil
	}

	// Iterate over the fields to be selected
	for _, fieldName := range fields {
		// Find the field by JSON tag
		field, found := findFieldByJSONTag(dataType, fieldName)
		if !found {
			// Ignore fields not found in the struct
			continue
		}

		// field names must be exported to access the value
		if field.PkgPath == "" {
			// Get the field value
			fieldValue := dataValue.FieldByName(field.Name).Interface()

			// Add the field and value to the result map
			result[fieldName] = fieldValue
		} else {
			result[fieldName] = nil
		}
	}

	return result, nil
}

// StructListToMapList converts a list of structs to a list of maps with selected fields.
func StructMultipleToMapMultiple(data interface{}, fields []string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

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
