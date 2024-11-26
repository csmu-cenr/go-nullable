package nullable

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	BAD_REQUEST                                    = `bad request`
	COMMA                                          = `,`
	LEFT_SQUARE_BRACKET                            = `[`
	LEFT_AND_RIGHT_MUST_BE_STRUCTS                 = `left and right must be structs`
	LEFT_AND_RIGHT_MUST_HAVE_EQUAL_NO_OF_FIELDS    = `left and right must have qual no of fields`
	LEFT_AND_RIGHT_NAME_TYPE_AND_TAG_MUST_BE_EQUAL = `left and right name, type and tag must be equal`
	JSON                                           = `json`
	NOTHING                                        = ``
	MODIFY_READ_ONLY                               = `modify read only`
	MODIFIED                                       = `Modified`
	NIL_POINTER                                    = `nil pointer`
	READ_ONLY                                      = `read_only`
	READONLY                                       = `ReadOnly`
	SELECTED                                       = `Selected`
	SET_DATA                                       = `set data`
	SET_NULLABLE                                   = `set nullable`
	UNEXPECTED_ERROR                               = `unexpected error`
	VARIABLE_MUST_BE_A_STRUCT                      = `variable must be a struct`
)

// Nullable represents data that also can be NULL
type Nullable[T any] struct {
	Data     T
	Modified bool
	ReadOnly bool
	Selected bool
	Valid    bool
}

func (n Nullable[T]) GoString() string {
	var ref T
	return fmt.Sprintf("nullable.Nullable[%T]{Data:%#v,Valid:%#v,Selected:%#v,ReadOnly:%#v}", ref, n.Data, n.Valid, n.Selected, n.ReadOnly)
}

// Set assigns a Nullable[model] as well as selected and valid.
func (n *Nullable[T]) Set(data Nullable[T]) error {
	if n == nil {
		message := ErrorMessage{Message: NIL_POINTER, Attempted: SET_NULLABLE, Details: data}
		return message
	}
	if n.ReadOnly {
		m := ErrorMessage{ErrorNo: http.StatusBadRequest, Message: BAD_REQUEST, Attempted: MODIFY_READ_ONLY, Details: data}
		return m
	}
	// Compare current data with the new data
	if !reflect.DeepEqual(n.Data, data.Data) {
		n.Data = data.Data
		n.Modified = true
	}

	n.Valid = true
	n.Selected = true
	return nil
}

// SetData assigns a value as well as selected and valid.
func (n *Nullable[T]) SetData(data T) error {
	if n == nil {
		m := ErrorMessage{ErrorNo: http.StatusBadRequest, Message: NIL_POINTER, Attempted: SET_DATA, Details: data}
		return m
	}
	if n.ReadOnly {
		m := ErrorMessage{ErrorNo: http.StatusBadRequest, Message: BAD_REQUEST, Attempted: MODIFY_READ_ONLY, Details: data}
		return m
	}

	// Compare current data with the new data
	if !reflect.DeepEqual(n.Data, data) {
		n.Data = data
		n.Modified = true
	}

	n.Valid = true
	n.Selected = true
	return nil
}

// True returns whether or not the data is true.
func (n *Nullable[T]) True() bool {
	if n == nil {
		return false
	}

	switch data := any(n.Data).(type) {
	case string:
		switch data {
		case "true", "yes", "y":
			return true
		case "TRUE", "YES", "Y":
			return true
		default:
			return false
		}
	case float32:
		if data == 1 {
			return true
		}
		return false
	case float64:
		if data == 1 {
			return true
		}
		return false
	case int:
		if data == 1 {
			return true
		}
		return false
	case int32:
		if data == 1 {
			return true
		}
		return false
	case int64:
		if data == 1 {
			return true
		}
		return false
	case bool:
		return data
	default:
		return false
	}
}

// ValueOrZero Get Value, or default zero value if it is NULL
func (n Nullable[T]) ValueOrZero() T {
	if !n.Valid {
		var ref T
		return ref
	}
	return n.Data
}

// IsEmpty is syntactic sugar for IsZero
func (n Nullable[T]) IsEmpty() bool {
	if !n.Valid {
		return true
	}
	var ref T
	return any(ref) == any(n.Data)
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

// CopyLeftToRight
func CopyLeftToRight(left, right reflect.Value, keepRight bool, setRightSelectedFalse bool, nullablesOnly bool, fields []string) error {

	functionName := `CopyLeftToRight`

	// Dereference pointers if necessary
	if left.Kind() == reflect.Ptr {
		left = left.Elem()
	}
	if right.Kind() == reflect.Ptr {
		right = right.Elem()
	}

	// Ensure both left and right are structs
	if left.Kind() != reflect.Struct || right.Kind() != reflect.Struct {
		m := ErrorMessage{
			Details:  LEFT_AND_RIGHT_MUST_BE_STRUCTS,
			ErrorNo:  http.StatusBadRequest,
			Function: functionName,
			Message:  BAD_REQUEST,
		}
		return m
	}

	leftType := left.Type()
	for i := 0; i < left.NumField(); i++ {

		field := leftType.Field(i)
		fieldName := strings.Split(field.Tag.Get(`json`), COMMA)[0]
		if fieldName == NOTHING {
			fieldName = field.Name
		}

		leftField := left.Field(i)
		rightField := right.Field(i)

		if leftField.Kind() == reflect.Ptr {
			leftField = leftField.Elem()
		}
		if rightField.Kind() == reflect.Ptr {
			rightField = rightField.Elem()
		}

		process := len(fields) == 0
		if len(fields) > 0 {
			process = FieldsContainsName(fields, fieldName)
		}
		if !process {
			continue
		}

		if IsNullable(leftField) && IsNullable(rightField) {

			// Ensure both fields are structs
			if leftField.Kind() != reflect.Struct || rightField.Kind() != reflect.Struct {
				continue
			}
			rightSelected := rightField.FieldByName(SELECTED).Bool()

			if rightSelected {
				if setRightSelectedFalse {
					err := SetNullableField(false, SELECTED, rightField)
					if err != nil {
						message := ErrorMessage{
							Attempted: `SetValueToNullableField`,
							Details:   err.Error(),
							ErrorNo:   http.StatusBadRequest,
							Function:  functionName,
							Message:   `setRightSelectedFalse`,
						}
						return message
					}
				}
				if keepRight {
					continue
				}
			}
			err := setReflectValueFromToField(leftField, rightField, "Data")
			if err != nil {
				message := ErrorMessage{Message: UNEXPECTED_ERROR, Details: "Data", Attempted: `setReflectValueFromToField`, InnerError: err}
				return message
			}
			err = setReflectValueFromToField(leftField, rightField, "Valid")
			if err != nil {
				message := ErrorMessage{Message: UNEXPECTED_ERROR, Details: "Valid", Attempted: `setReflectValueFromToField`, InnerError: err}
				return message
			}

		} else if !nullablesOnly {
			// process non nullables and primatives
			if leftField.IsValid() {
				if rightField.IsValid() {
					if rightField.CanSet() {
						rightField.Set(leftField)
					}
				}
			}
		}
	}

	return nil
}

func fieldNameIsInFields(fieldName string, fields []string) bool {
	for i := 0; i < len(fields); i++ {
		if fieldName == fields[i] {
			return true
		}
	}
	return false
}

func FieldsContainsName(fields []string, name string) bool {
	for _, field := range fields {
		if strings.EqualFold(name, field) {
			return true
		}
	}
	return false
}

func FindModifiedFields(data any) []string {
	result := []string{}

	val := reflect.ValueOf(data)
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
		if field.Kind() == reflect.Struct && hasField(fieldType.Type, SELECTED) {
			modified := field.FieldByName("Modified")
			if modified.IsValid() && modified.Bool() {
				result = append(result, jsonTag)
			}
		}
	}

	return result
}

func FindSelectedFields(data any) []string {
	result := []string{}

	val := reflect.ValueOf(data)
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
		if field.Kind() == reflect.Struct && hasField(fieldType.Type, SELECTED) {
			selected := field.FieldByName(SELECTED)
			if selected.IsValid() && selected.Bool() {
				result = append(result, jsonTag)
			}
		}
	}

	return result
}

func GetModifiedFieldTags(model any) []string {
	result := []string{}

	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		structField := modelType.Field(i)
		fieldName := structField.Name
		jsonTag := structField.Tag.Get("json")
		if strings.Contains(jsonTag, ",") {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}
		if jsonTag == "" {
			jsonTag = fieldName
		}
		if field.Kind() == reflect.Struct && hasField(structField.Type, "Modified") {
			modified := field.FieldByName("Modified")
			if modified.IsValid() && modified.Bool() {
				result = append(result, jsonTag)
			}
		}
	}

	return result
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

func GetSelectedFields(any interface{}, fields []string) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(any)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name
		key := fieldType.Tag.Get("json")
		if strings.Contains(key, ",") {
			key = strings.Split(key, ",")[0]
		}
		if key == "" {
			key = fieldName
		}
		fieldNameInFields := fieldNameIsInFields(key, fields)
		if field.Kind() == reflect.Struct && hasField(fieldType.Type, SELECTED) {
			selectedField := field.FieldByName(SELECTED)
			if (selectedField.IsValid() && selectedField.Bool()) || fieldNameInFields {
				result[key] = field.Interface()
			}
		} else if field.Kind() == reflect.Struct {
			value := field.Interface()
			mapped := GetSelectedFields(value, fields)
			if len(mapped) > 0 {
				result[key] = mapped
			}
		} else if field.Kind() == reflect.Slice {
			if field.Len() > 0 {
				value := field.Interface()
				result[key] = GetSelectedFieldsSlice(value, fields)
			}
		} else if fieldNameInFields {
			result[key] = field.Interface()
		}
	}

	return result
}

func getTag(field reflect.StructField) string {
	return strings.Split(field.Tag.Get(JSON), COMMA)[0]
}

func hasField(typ reflect.Type, fieldName string) bool {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			return true
		}
	}
	return false
}

// IsSelectedEqual checks all the selected Nullable fields in the left instance against the right object
func IsSelectedEqual(left, right reflect.Value) bool {

	//functionName := `IsSelectedEqual`

	// Dereference pointers if necessary
	if left.Kind() == reflect.Ptr {
		left = left.Elem()
	}
	if right.Kind() == reflect.Ptr {
		right = right.Elem()
	}

	// Ensure both left and right are structs
	if left.Kind() != reflect.Struct || right.Kind() != reflect.Struct {
		return false
	}

	leftType := left.Type()
	rightTyoe := right.Type()
	for i := 0; i < left.NumField(); i++ {

		leftFieldType := leftType.Field(i)
		rightFieldType := rightTyoe.Field(i)
		if strings.EqualFold(leftFieldType.Name, rightFieldType.Name) {
			return false
		}

		leftField := left.Field(i)
		rightField := right.Field(i)

		if leftField.Kind() == reflect.Ptr {
			leftField = leftField.Elem()
		}
		if rightField.Kind() == reflect.Ptr {
			rightField = rightField.Elem()
		}

		if IsNullable(leftField) && IsNullable(rightField) {

			// Ensure both fields are structs
			if leftField.Kind() != reflect.Struct || rightField.Kind() != reflect.Struct {
				continue
			}

			leftSelected := leftField.FieldByName(SELECTED).Bool()
			rightSelected := rightField.FieldByName(SELECTED).Bool()

			if leftSelected == rightSelected {
				leftData := leftField.FieldByName("Data")
				rightData := leftField.FieldByName("Data")
				equal := reflect.DeepEqual(leftData, rightData)
				if !equal {
					return false
				}
			}
		}
	}

	return true
}

// IsNullable
func IsNullable(model reflect.Value) bool {
	name := strings.Split(model.Type().Name(), LEFT_SQUARE_BRACKET)[0]
	return name == `Nullable`
}

// Modified returns true if one or more Nullable fields are selected. Otherwise false is returned
func Modified(model any) bool {

	// Get the type and inputValue of the input data
	inputValue := reflect.ValueOf(model)

	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	if inputValue.Kind() == reflect.Struct {

		var field reflect.Value

		for i := 0; i < inputValue.NumField(); i++ {

			field = inputValue.Field(i) // Get the reflecion value

			if field.Type().Kind() == reflect.Struct {
				if IsNullable(field) {
					modified := field.FieldByName(MODIFIED)
					if modified.IsValid() && modified.Kind() == reflect.Bool && modified.Bool() {
						return true
					}
				}
			}
		}
	}

	return false
}

// ModifiedFields returns all nullable Null structs that are modified
func ModifiedFields(input any) []string {

	result := []string{}

	inputValue := reflect.ValueOf(input)
	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	if inputValue.Kind() == reflect.Struct {

		inputType := reflect.TypeOf(input)
		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
		}

		var tag string
		for i := 0; i < inputValue.NumField(); i++ {

			field := inputType.Field(i)
			value := inputValue.Field(i)

			tag = getTag(field)
			if tag == NOTHING {
				tag = field.Name
			}

			if value.Type().Kind() == reflect.Struct {
				if IsNullable(value) {
					readOnly := value.FieldByName(READ_ONLY)
					if readOnly.IsValid() && readOnly.Kind() == reflect.Bool && readOnly.Bool() {
						continue
					}
					modified := value.FieldByName("Modified")
					if modified.IsValid() && modified.Kind() == reflect.Bool && modified.Bool() {
						result = append(result, tag)
					}
				}
			}

		}
	}

	return result
}

// SelectedFields is a generic function to get JSON tags of selected Nullable fields
func SelectedFields[T any](input T, includeNonNullable bool) []string {
	var fields []string

	inputValue := reflect.ValueOf(input)
	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	if inputValue.Kind() == reflect.Struct {

		inputType := reflect.TypeOf(input)
		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
		}

		tag := NOTHING
		for i := 0; i < inputValue.NumField(); i++ {

			field := inputType.Field(i)
			value := inputValue.Field(i)

			tag = getTag(field)
			if tag == NOTHING {
				tag = field.Name
			}

			if value.Type().Kind() == reflect.Struct {
				if IsNullable(value) {
					selected := value.FieldByName(`Selected`)
					if selected.IsValid() && selected.Kind() == reflect.Bool && selected.Bool() {
						fields = append(fields, tag)
					}
				} else if includeNonNullable {
					fields = append(fields, tag)
				}
			} else if includeNonNullable {
				fields = append(fields, tag)
			}

		}

	}

	return fields
}

// setBooleanFields sets every field in fields to the target, all others are set to not the target.
func setBooleanFields(instance reflect.Value, fields []string, fieldName string, target bool, not bool) error {
	functionName := `SetNullableBooleanFields`
	var err error

	if len(fields) == 0 {
		return nil
	}

	// Dereference pointers if necessary
	if instance.Kind() == reflect.Ptr {
		instance = instance.Elem()
	}

	// Ensure instance is a struct
	if instance.Kind() != reflect.Struct {
		m := ErrorMessage{
			Details:  VARIABLE_MUST_BE_A_STRUCT,
			ErrorNo:  http.StatusBadRequest,
			Function: functionName,
			Message:  BAD_REQUEST,
		}
		return m
	}

	process := false
	value := false
	instanceType := instance.Type()
	for i := 0; i < instance.NumField(); i++ {

		field := instanceType.Field(i)
		tag := strings.Split(field.Tag.Get(`json`), COMMA)[0]
		if tag == NOTHING {
			tag = field.Name
		}

		instanceField := instance.Field(i)
		if instanceField.Kind() == reflect.Ptr {
			instanceField = instanceField.Elem()
		}
		if instanceField.Kind() != reflect.Struct {
			continue
		}
		if !IsNullable(instanceField) {
			continue
		}

		process = FieldsContainsName(fields, tag)
		value = instanceField.FieldByName(fieldName).Bool()

		if process {
			switch {
			case value != target: // 0
				err = SetNullableField(target, fieldName, instanceField)
			default: // 1
				continue
			}
		}

		if not && !process {
			switch {
			case value != target: // 0
				continue
			default: // 1
				err = SetNullableField(!target, fieldName, instanceField)
			}
		}

		if err != nil {
			m := ErrorMessage{
				Attempted:  `setReflectValueFromToField`,
				Details:    "Data",
				ErrorNo:    http.StatusInternalServerError,
				InnerError: err,
				Message:    UNEXPECTED_ERROR,
			}
			return m
		}
	}

	return nil
}

// SetModifiedBooleanFields calls SetBooleanFields with 'Modified' as the field name
func SetModifiedBooleanFields(instance reflect.Value, fields []string, target bool, not bool) error {
	return setBooleanFields(instance, fields, `Modified`, target, not)
}

// SetModifiedIfDifferent sets any field in the left struct to modified if different from the right
func SetModifiedIfDifferent(modify, base reflect.Value) error {

	functionName := `CopyLeftToRight`

	// Dereference pointers if necessary
	if modify.Kind() == reflect.Ptr {
		modify = modify.Elem()
	}
	if base.Kind() == reflect.Ptr {
		base = base.Elem()
	}

	// Ensure both left and right are structs
	if modify.Kind() != reflect.Struct || base.Kind() != reflect.Struct {
		m := ErrorMessage{
			Details:  LEFT_AND_RIGHT_MUST_BE_STRUCTS,
			ErrorNo:  http.StatusBadRequest,
			Function: functionName,
			Message:  BAD_REQUEST,
		}
		return m
	}

	if modify.NumField() != base.NumField() {
		m := ErrorMessage{
			Details:  LEFT_AND_RIGHT_MUST_HAVE_EQUAL_NO_OF_FIELDS,
			ErrorNo:  http.StatusBadRequest,
			Function: functionName,
			Message:  BAD_REQUEST,
		}
		return m
	}

	modifyType := modify.Type()
	baseType := base.Type()

	for i := 0; i < modify.NumField(); i++ {

		modifyTypeField := modifyType.Field(i)
		baseTypeField := baseType.Field(i)
		if !strings.EqualFold(modifyTypeField.Name, baseTypeField.Name) && modifyTypeField.Type != baseTypeField.Type && modifyTypeField.Tag != baseTypeField.Tag {
			m := ErrorMessage{
				Details:  LEFT_AND_RIGHT_NAME_TYPE_AND_TAG_MUST_BE_EQUAL,
				ErrorNo:  http.StatusBadRequest,
				Function: functionName,
				Message:  BAD_REQUEST,
			}
			return m
		}
		modifiyField := modify.Field(i)
		baseField := base.Field(i)

		if modifiyField.Kind() == reflect.Ptr {
			modifiyField = modifiyField.Elem()
		}
		if baseField.Kind() == reflect.Ptr {
			baseField = baseField.Elem()
		}
		if IsNullable(modifiyField) {
			if modifiyField.IsValid() && baseField.IsValid() {
				modifySelectedField := modifiyField.FieldByName(SELECTED)
				if !(modifySelectedField.IsValid() || modifySelectedField.CanInterface()) {
					continue
				}
				modifySelected := modifySelectedField.Bool()
				if modifySelected {
					modifyReadOnlyField := modifiyField.FieldByName("ReadOnly")
					if !(modifyReadOnlyField.IsValid() || modifyReadOnlyField.CanInterface()) {
						continue
					}
					modifyReadOnly := modifyReadOnlyField.Bool()
					if modifyReadOnly {
						continue
					}
					modifyData := modifiyField.FieldByName("Data")
					if !(modifyData.IsValid() || modifyData.CanInterface()) {
						continue
					}
					baseData := baseField.FieldByName("Data")
					if !(baseData.IsValid() || baseData.CanInterface()) {
						continue
					}
					if !reflect.DeepEqual(modifyData, baseData) {
						baseModified := baseField.FieldByName("Modified")
						if !(baseModified.IsValid() || baseModified.CanInterface()) {
							continue
						}
						if baseModified.IsValid() && baseModified.Kind() == reflect.Bool {
							err := SetNullableField(true, `Modified`, baseModified)
							if err != nil {
								m := ErrorMessage{
									Attempted: `SetNullableField`,
									Details:   fmt.Sprintf(`FieldName: %s Err: %+v`, modifyType.Name(), err),
									ErrorNo:   http.StatusInternalServerError,
									Function:  functionName,
									Message:   UNEXPECTED_ERROR,
								}
								return m
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// SetModifiedIfSelected is used for data that
// is unmarshalled but needs to be saved.
func SetModifiedIfSelected(model any) error {

	function := `SetModifiedIfSelected`

	// Get the type and inputValue of the input data
	inputValue := reflect.ValueOf(model)

	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	if inputValue.Kind() != reflect.Struct {
		return nil
	}

	structType := inputValue.Type()

	for i := 0; i < inputValue.NumField(); i++ {
		field := inputValue.Field(i)
		fieldName := structType.Field(i).Name

		if field.Type().Kind() != reflect.Struct {
			continue
		}

		if !IsNullable(field) {
			continue
		}

		selected := field.FieldByName(SELECTED)
		modified := field.FieldByName("Modified")

		// Ensure fields are valid and settable
		if selected.IsValid() && selected.Kind() == reflect.Bool && selected.Bool() {
			if modified.IsValid() && modified.Kind() == reflect.Bool {
				err := SetNullableField(true, `Modified`, field)
				if err != nil {
					m := ErrorMessage{
						Attempted: `SetNullableField`,
						Details:   fmt.Sprintf(`FieldName: %s Err: %+v`, fieldName, err),
						ErrorNo:   http.StatusInternalServerError,
						Function:  function,
						Message:   UNEXPECTED_ERROR,
					}
					return m
				}
			}
		}

	}

	return nil
}

// SetNullableField sets the value to a specific field in a Nullable struct
func SetNullableField(value any, fieldName string, nullableField reflect.Value) error {
	if nullableField.Kind() == reflect.Ptr {
		nullableField = nullableField.Elem()
	}

	if nullableField.Kind() != reflect.Struct {
		return fmt.Errorf("nullableField is not a struct")
	}

	field := nullableField.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	val := reflect.ValueOf(value)

	switch fieldName {
	case "Data":
		if field.Type() != val.Type() {
			return fmt.Errorf("value type %s does not match field type %s", val.Type(), field.Type())
		}
		field.Set(val)
	case "Valid", SELECTED, "Modified":
		if val.Kind() != reflect.Bool {
			return fmt.Errorf("value type %s does not match field type bool", val.Kind())
		}
		field.SetBool(val.Bool())
	default:
		return fmt.Errorf("unsupported field name: %s", fieldName)
	}

	return nil
}

// SetReadOnlyBooleanFields calls SetBooleanFields with "ReadOnly" as the field name
func SetReadOnlyBooleanFields(instance reflect.Value, fields []string, target bool, not bool) error {
	return setBooleanFields(instance, fields, READ_ONLY, target, not)
}

// setReflectValueFromToField sets the value from the `from` struct to the `to` struct for the specified field.
func setReflectValueFromToField(from, to reflect.Value, field string) error {
	// Dereference the pointers if necessary
	if from.Kind() == reflect.Ptr {
		from = from.Elem()
	}
	if to.Kind() == reflect.Ptr {
		to = to.Elem()
	}

	// Ensure both from and to are structs
	if from.Kind() != reflect.Struct || to.Kind() != reflect.Struct {
		return errors.New("both 'from' and 'to' must be structs or pointers to structs")
	}

	// Get the field from both structs
	fromField := from.FieldByName(field)
	toField := to.FieldByName(field)

	// Check if the field exists in both structs
	if !fromField.IsValid() {
		return fmt.Errorf("field %s not found in 'from' struct", field)
	}
	if !toField.IsValid() {
		return fmt.Errorf("field %s not found in 'to' struct", field)
	}

	// Ensure the toField is settable
	if !toField.CanSet() {
		return fmt.Errorf("field %s cannot be set in 'to' struct", field)
	}

	// Set the value
	toField.Set(fromField)
	return nil
}

// SetSelectedBooleanFields calls SetBooleanFields with "Selected" as the field name
func SetSelectedBooleanFields(instance reflect.Value, fields []string, target bool, not bool) error {
	return setBooleanFields(instance, fields, `Selected`, target, not)
}

// Null Create a Nullable that is NULL with type
func Null[T any]() Nullable[T] {
	return Nullable[T]{}
}

// Value Create a Nullable from a value
func Value[T any](value T) Nullable[T] {
	if any(value) == nil {
		return Nullable[T]{Valid: false, Selected: true}
	}
	return Nullable[T]{Data: value, Valid: true, Selected: true}
}

// ValueFromPointer Create a Nullable from a pointer
func ValueFromPointer[T any](value *T) Nullable[T] {
	if value == nil {
		return Nullable[T]{Valid: false, Selected: true}
	}
	return Value(*value)
}
