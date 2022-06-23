package nullable

import (
	"encoding"
	"errors"
	"fmt"
)

func (n Nullable[T]) MarshalText() ([]byte, error) {
	if !n.HasValue {
		return []byte{}, nil
	}

	value := any(n.Value)
	txt, ok := value.(encoding.TextMarshaler)
	if ok {
		return txt.MarshalText()
	}

	b, ok := any(n).(Nullable[bool])
	if ok {
		return marshalTextBool(b)
	}

	var ref T
	return []byte{}, fmt.Errorf("type %T cannot be marshalled to text", ref)
}

func marshalTextBool(b Nullable[bool]) ([]byte, error) {
	if !b.HasValue {
		return []byte{}, nil
	}
	if !b.Value {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (n *Nullable[T]) UnmarshalText(text []byte) error {
	value := any(&n.Value)
	txt, ok := value.(encoding.TextUnmarshaler)
	if ok {
		return txt.UnmarshalText(text)
	}

	b, ok := any(n).(*Nullable[bool])
	if ok {
		return unmarshalTextBool(text, b)
	}

	var ref T
	return fmt.Errorf("type %T unmarshal", ref)
}

func unmarshalTextBool(text []byte, b *Nullable[bool]) error {
	str := string(text)
	switch str {
	case "", "null":
		b.HasValue = false
		return nil
	case "true":
		b.Value = true
	case "false":
		b.Value = false
	default:
		return errors.New("null: invalid input for UnmarshalText:" + str)
	}
	b.HasValue = true
	return nil
}
