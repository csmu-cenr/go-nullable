package nullable

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	Value    T
	HasValue bool
}

func Value[T any](value T) Nullable[T] {
	if any(value) == nil {
		return Nullable[T]{HasValue: false}
	}
	return Nullable[T]{Value: value, HasValue: true}
}

func ValueFromPtr[T any](value *T) Nullable[T] {
	if value == nil {
		return Nullable[T]{HasValue: false}
	}
	return Value(*value)
}

func (n Nullable[T]) ValueOrZero() T {
	if !n.HasValue {
		var ref T
		return ref
	}
	return n.Value
}

func (n Nullable[T]) Equal(other Nullable[T]) bool {
	if n.HasValue != other.HasValue {
		return false
	}
	if !n.HasValue {
		return true // nil == nil
	}
	return any(n.Value) == any(other.Value)
}

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
	if err != nil {
		return fmt.Errorf("null: could not unmarshal JSON: %w", err)
	}

	n.HasValue = true
	return nil
}

func (n Nullable[T]) String() string {
	return fmt.Sprintf("%s", any(n.Value))
}
