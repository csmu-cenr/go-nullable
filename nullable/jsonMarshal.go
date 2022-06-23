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

	f, ok := any(n).(*Nullable[float64])
	if ok {
		return unmarshalFloat64StringJson(f, data)
	}

	return fmt.Errorf("null: could not unmarshal JSON: %w", err)
}

func unmarshalFloat64StringJson(f *Nullable[float64], data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
	}
	n, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("null: couldn't convert string to float: %w", err)
	}

	f.Value = n
	f.HasValue = true
	return nil
}
