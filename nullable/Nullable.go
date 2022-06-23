package nullable

import (
	"fmt"
	"time"
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
	switch any(n.Value).(type) {
	case time.Time:
		nValue := any(n.Value).(time.Time)
		otherValue := any(other.Value).(time.Time)
		return n.HasValue == other.HasValue && (!n.HasValue || nValue.Equal(otherValue))
	}
	return n.ExactEqual(other)
}

func (n Nullable[T]) ExactEqual(other Nullable[T]) bool {
	return n.HasValue == other.HasValue && (!n.HasValue || any(n.Value) == any(other.Value))
}

func (n Nullable[T]) String() string {
	return fmt.Sprintf("%s", any(n.Value))
}
