package nullable

import (
	"encoding"
	"errors"
	"fmt"
	"strconv"
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

	f, ok := any(n).(Nullable[float64])
	if ok {
		return marshalTextFloat64(f)
	}

	var ref T
	return []byte{}, fmt.Errorf("type %T cannot be marshalled to text", ref)
}

func marshalTextFloat64(f Nullable[float64]) ([]byte, error) {
	if !f.HasValue {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Value, 'f', -1, 64)), nil
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

	f, ok := any(n).(*Nullable[float64])
	if ok {
		return unmarshalTextFloat64(text, f)
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

func unmarshalTextFloat64(text []byte, f *Nullable[float64]) error {
	str := string(text)
	if str == "" || str == "null" {
		f.HasValue = false
		return nil
	}
	var err error
	f.Value, err = strconv.ParseFloat(string(text), 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	f.HasValue = true
	return err
}
