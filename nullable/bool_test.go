package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoolFrom(t *testing.T) {
	b := Value(true)
	assertBool(t, b, "Value()")

	zero := Value(false)
	if !zero.HasValue {
		t.Error("Value(false)", "is invalid, but should be valid")
	}
}

func TestBoolFromPtr(t *testing.T) {
	n := true
	boolPointer := &n
	b := ValueFromPtr(boolPointer)
	assertBool(t, b, "ValueFromPtr()")

	null := ValueFromPtr[bool](nil)
	assertNull(t, null, "ValueFromPtr(nil)")
}

func TestUnmarshalBool(t *testing.T) {
	var b Nullable[bool]
	err := json.Unmarshal(boolJSON, &b)
	assert.Nil(t, err)
	assertBool(t, b, "bool json")

	var nb Nullable[bool]
	err = json.Unmarshal(nullBoolJSON, &nb)
	if err == nil {
		panic("err should not be nil")
	}

	var null Nullable[bool]
	err = json.Unmarshal(nullJSON, &null)
	assert.Nil(t, err)
	assertNull(t, null, "null json")

	var badType Nullable[bool]
	err = json.Unmarshal(intJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNull(t, badType, "wrong type json")

	var invalid Nullable[bool]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalBool(t *testing.T) {
	var b Nullable[bool]
	err := b.UnmarshalText([]byte("true"))
	assert.Nil(t, err)
	assertBool(t, b, "UnmarshalText() bool")

	var zero Nullable[bool]
	err = zero.UnmarshalText([]byte("false"))
	assert.Nil(t, err)
	assertFalseBool(t, zero, "UnmarshalText() false")

	var blank Nullable[bool]
	err = blank.UnmarshalText([]byte(""))
	assert.Nil(t, err)
	assertNull(t, blank, "UnmarshalText() empty bool")

	var null Nullable[bool]
	err = null.UnmarshalText([]byte("null"))
	assert.Nil(t, err)
	assertNull(t, null, `UnmarshalText() "null"`)

	var invalid Nullable[bool]
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("err should not be nil")
	}
	assertNull(t, invalid, "invalid json")
}

func TestMarshalBool(t *testing.T) {
	b := Value(true)
	data, err := json.Marshal(b)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "true", "non-empty json marshal")

	zero := Value(false)
	data, err = json.Marshal(zero)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "false", "zero json marshal")

	// invalid values should be encoded as null
	null := Nullable[bool]{Value: false, HasValue: false}
	data, err = json.Marshal(null)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalBoolText(t *testing.T) {
	b := Value(true)
	data, err := b.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "true", "non-empty text marshal")

	zero := Value(false)
	data, err = zero.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "false", "zero text marshal")

	// invalid values should be encoded as null
	null := Nullable[bool]{HasValue: false}
	data, err = null.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBoolValueOrZero(t *testing.T) {
	valid := Value(true)
	if valid.ValueOrZero() != true {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[bool]{Value: true, HasValue: false}
	if invalid.ValueOrZero() != false {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestBoolEqual(t *testing.T) {
	b1 := Nullable[bool]{Value: true, HasValue: false}
	b2 := Nullable[bool]{Value: true, HasValue: false}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: false}
	b2 = Nullable[bool]{Value: false, HasValue: false}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: true}
	b2 = Nullable[bool]{Value: true, HasValue: true}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: true}
	b2 = Nullable[bool]{Value: true, HasValue: false}
	assertNotEqual(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: false}
	b2 = Nullable[bool]{Value: true, HasValue: true}
	assertNotEqual(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: true}
	b2 = Nullable[bool]{Value: false, HasValue: true}
	assertNotEqual(t, b1, b2)
}

func assertBool(t *testing.T, b Nullable[bool], from string) {
	if b.Value != true {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Value, true)
	}
	if !b.HasValue {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFalseBool(t *testing.T, b Nullable[bool], from string) {
	if b.Value != false {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Value, false)
	}
	if !b.HasValue {
		t.Error(from, "is invalid, but should be valid")
	}
}
