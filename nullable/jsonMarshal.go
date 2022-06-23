package nullable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

var nullBytes = []byte("null")

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.HasValue {
		return nullBytes, nil
	}
	return json.Marshal(n.Value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		n.HasValue = false
		return nil
	}

	err := json.Unmarshal(data, &n.Value)
	if err == nil {
		n.HasValue = true
		return nil
	}

	switch any(n.Value).(type) {
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
	v := any(f.Value)
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
		f.Value = any(float32(n)).(T)
	case float64:
		f.Value = any(n).(T)
	}

	f.HasValue = true
	return nil
}

func unmarshalIntStringJson[T any](f *Nullable[T], data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
	}

	var size int
	v := any(f.Value)
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
		f.Value = any(int8(n)).(T)
	case int16:
		f.Value = any(int16(n)).(T)
	case int32:
		f.Value = any(int32(n)).(T)
	case int:
		f.Value = any(int(n)).(T)
	case int64:
		f.Value = any(n).(T)
	}

	f.HasValue = true
	return nil
}
