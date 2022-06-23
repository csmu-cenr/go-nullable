package nullable

import (
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

func (n Nullable[T]) String() string {
	return fmt.Sprintf("%s", any(n.Value))
}
