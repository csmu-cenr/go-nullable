package nullable

import (
	"fmt"
	"time"
)

// Nullable represents data that also can be NULL
type Nullable[T any] struct {
	Data  T
	Valid bool
}

// Value Create a Nullable from a value
func Value[T any](value T) Nullable[T] {
	if any(value) == nil {
		return Nullable[T]{Valid: false}
	}
	return Nullable[T]{Data: value, Valid: true}
}

// ValueFromPointer Create a Nullable from a pointer
func ValueFromPointer[T any](value *T) Nullable[T] {
	if value == nil {
		return Nullable[T]{Valid: false}
	}
	return Value(*value)
}

// Null Create a Nullable that is NULL with type
func Null[T any]() Nullable[T] {
	return Nullable[T]{}
}

// ValueOrZero Get Value, or default zero value if it is NULL
func (n Nullable[T]) ValueOrZero() T {
	if !n.Valid {
		var ref T
		return ref
	}
	return n.Data
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

func (n Nullable[T]) GoString() string {
	var ref T
	return fmt.Sprintf("nullable.Nullable[%T]{Data:%#v,Valid:%#v}", ref, n.Data, n.Valid)
}
