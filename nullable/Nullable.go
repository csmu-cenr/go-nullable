package nullable

import (
	"fmt"
	"time"
)

// Nullable represents data that also can be NULL
type Nullable[T any] struct {
	Data    T
	IsValid bool
}

// Value Create a Nullable from a value
func Value[T any](value T) Nullable[T] {
	if any(value) == nil {
		return Nullable[T]{IsValid: false}
	}
	return Nullable[T]{Data: value, IsValid: true}
}

// ValueFromPointer Create a Nullable from a pointer
func ValueFromPointer[T any](value *T) Nullable[T] {
	if value == nil {
		return Nullable[T]{IsValid: false}
	}
	return Value(*value)
}

// Null Create a Nullable that is NULL with type
func Null[T any]() Nullable[T] {
	return Nullable[T]{}
}

// ValueOrZero Get Value, or default zero value if it is NULL
func (n Nullable[T]) ValueOrZero() T {
	if !n.IsValid {
		var ref T
		return ref
	}
	return n.Data
}

// Equal Check if this Nullable is equal to another Nullable
func (n Nullable[T]) Equal(other Nullable[T]) bool {
	switch any(n.Data).(type) {
	case time.Time:
		nValue := any(n.Data).(time.Time)
		otherValue := any(other.Data).(time.Time)
		return n.IsValid == other.IsValid && (!n.IsValid || nValue.Equal(otherValue))
	}
	return n.ExactEqual(other)
}

// ExactEqual Check if this Nullable is exact equal to another Nullable, never using intern Equal method to check equality
func (n Nullable[T]) ExactEqual(other Nullable[T]) bool {
	return n.IsValid == other.IsValid && (!n.IsValid || any(n.Data) == any(other.Data))
}

// String Convert value to string
func (n Nullable[T]) String() string {
	return fmt.Sprintf("%s", any(n.Data))
}
