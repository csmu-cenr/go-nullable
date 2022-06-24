package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoolFromValue(t *testing.T) {
	b := Value(true)
	assertBool(t, b, "Data()")

	zero := Value(false)
	if !zero.IsValid {
		t.Error("Data(false)", "is invalid, but should be valid")
	}
}

func TestBoolFromPointer(t *testing.T) {
	n := true
	boolPointer := &n
	b := ValueFromPointer(boolPointer)
	assertBool(t, b, "ValueFromPointer()")

	null := ValueFromPointer[bool](nil)
	assert.False(t, null.IsValid)
}

func TestUnmarshalBool(t *testing.T) {
	var b Nullable[bool]
	err := json.Unmarshal(boolJSON, &b)
	assert.Nil(t, err)
	assertBool(t, b, "bool json")

	var nb Nullable[bool]
	err = json.Unmarshal(nullBoolJSON, &nb)
	assert.NotNil(t, err)

	var null Nullable[bool]
	err = json.Unmarshal(nullJSON, &null)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)

	var badType Nullable[bool]
	err = json.Unmarshal(intJSON, &badType)
	assert.NotNil(t, err)
	assert.False(t, badType.IsValid)

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

	var falseBool Nullable[bool]
	err = falseBool.UnmarshalText([]byte("false"))
	assert.Nil(t, err)
	assertFalseBool(t, falseBool, "UnmarshalText() false")

	var fromBlankString Nullable[bool]
	err = fromBlankString.UnmarshalText([]byte(""))
	assert.Nil(t, err)
	assert.False(t, fromBlankString.IsValid)

	var fromNullString Nullable[bool]
	err = fromNullString.UnmarshalText([]byte("null"))
	assert.Nil(t, err)
	assert.False(t, fromNullString.IsValid)

	var invalid Nullable[bool]
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("err should not be nil")
	}
	assert.False(t, invalid.IsValid)
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
	null := Null[bool]()
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
	null := Null[bool]()
	data, err = null.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBoolValueOrZero(t *testing.T) {
	valid := Value(true)
	if valid.ValueOrZero() != true {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[bool]{Data: true, IsValid: false}
	if invalid.ValueOrZero() != false {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestBoolEqual(t *testing.T) {
	b1 := Nullable[bool]{Data: true, IsValid: false}
	b2 := Nullable[bool]{Data: true, IsValid: false}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, IsValid: false}
	b2 = Nullable[bool]{Data: false, IsValid: false}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, IsValid: true}
	b2 = Nullable[bool]{Data: true, IsValid: true}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, IsValid: true}
	b2 = Nullable[bool]{Data: true, IsValid: false}
	assertNotEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, IsValid: false}
	b2 = Nullable[bool]{Data: true, IsValid: true}
	assertNotEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, IsValid: true}
	b2 = Nullable[bool]{Data: false, IsValid: true}
	assertNotEqual(t, b1, b2)
}

func TestBoolScan(t *testing.T) {
	var b Nullable[bool]
	err := b.Scan(true)
	assert.Nil(t, err)
	assertBool(t, b, "scanned bool")

	var null Nullable[bool]
	err = null.Scan(nil)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)
}

func assertBool(t *testing.T, b Nullable[bool], from string) {
	if b.Data != true {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Data, true)
	}
	if !b.IsValid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFalseBool(t *testing.T, b Nullable[bool], from string) {
	if b.Data != false {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Data, false)
	}
	if !b.IsValid {
		t.Error(from, "is invalid, but should be valid")
	}
}
