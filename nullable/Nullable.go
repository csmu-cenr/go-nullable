package nullable

import (
	"fmt"
	"time"
)

type Nullable[T any] struct {
	Data    T
	IsValid bool
}

func Value[T any](value T) Nullable[T] {
	if any(value) == nil {
		return Nullable[T]{IsValid: false}
	}
	return Nullable[T]{Data: value, IsValid: true}
}

func ValueFromPointer[T any](value *T) Nullable[T] {
	if value == nil {
		return Nullable[T]{IsValid: false}
	}
	return Value(*value)
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{}
}

func (n Nullable[T]) ValueOrZero() T {
	if !n.IsValid {
		var ref T
		return ref
	}
	return n.Data
}

func (n Nullable[T]) Equal(other Nullable[T]) bool {
	switch any(n.Data).(type) {
	case time.Time:
		nValue := any(n.Data).(time.Time)
		otherValue := any(other.Data).(time.Time)
		return n.IsValid == other.IsValid && (!n.IsValid || nValue.Equal(otherValue))
	}
	return n.ExactEqual(other)
}

func (n Nullable[T]) ExactEqual(other Nullable[T]) bool {
	return n.IsValid == other.IsValid && (!n.IsValid || any(n.Data) == any(other.Data))
}

func (n Nullable[T]) String() string {
	return fmt.Sprintf("%s", any(n.Data))
}
