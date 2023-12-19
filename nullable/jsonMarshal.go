package nullable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var nullBytes = []byte("null")

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	if n.Selected {
		return json.Marshal(n.Data)
	}
	return OmitEmpty(n.Data)
}

// https://stackoverflow.com/a/77304217/162358
func OmitEmpty(data any) ([]byte, error) {
	value := reflect.ValueOf(data)
	t := value.Type()
	sf := make([]reflect.StructField, 0)
	// modify the 'for i' snippet for more complicated cases
	for i := 0; i < t.NumField(); i++ {
		sf = append(sf, t.Field(i))
		tag := t.Field(i).Tag
		if !strings.Contains(string(tag), ",omitempty") {
			r := regexp.MustCompile(`json:"\s*(.*?)\s*"`)
			matches := r.FindAllStringSubmatch(string(tag), -1)
			for _, v := range matches {
				tagKey := v[1]
				sf[i].Tag = reflect.StructTag(fmt.Sprintf(`json:"%s,omitempty"`, tagKey))
			}
		}
	}
	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)
	return json.Marshal(newValue.Interface())
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {

	n.Selected = true

	if bytes.Equal(data, nullBytes) {
		n.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &n.Data)
	if err == nil {
		n.Valid = true
		return nil
	}

	switch any(n.Data).(type) {
	case float32, float64:
		return unmarshalFloatStringJson(n, data)
	case int, int8, int16, int32, int64:
		return unmarshalIntStringJson(n, data)
	}

	return fmt.Errorf("null: could not unmarshal JSON: %w", err)
}

func unmarshalFloatStringJson[T any](f *Nullable[T], data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
	}

	var size int
	v := any(f.Data)
	switch v.(type) {
	case float32:
		size = 32
	case float64:
		size = 64
	}

	n, err := strconv.ParseFloat(str, size)
	if err != nil {
		return fmt.Errorf("null: couldn't convert string to float: %w", err)
	}

	switch v.(type) {
	case float32:
		f.Data = any(float32(n)).(T)
	case float64:
		f.Data = any(n).(T)
	}

	f.Valid = true

	return nil
}

func unmarshalIntStringJson[T any](f *Nullable[T], data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
	}

	var size int
	v := any(f.Data)
	switch v.(type) {
	case int8:
		size = 8
	case int16:
		size = 16
	case int32, int:
		size = 32
	case int64:
		size = 64
	}

	n, err := strconv.ParseInt(str, 10, size)
	if err != nil {
		return fmt.Errorf("null: couldn't convert string to float: %w", err)
	}

	switch v.(type) {
	case int8:
		f.Data = any(int8(n)).(T)
	case int16:
		f.Data = any(int16(n)).(T)
	case int32:
		f.Data = any(int32(n)).(T)
	case int:
		f.Data = any(int(n)).(T)
	case int64:
		f.Data = any(n).(T)
	}

	f.Valid = true
	return nil
}
