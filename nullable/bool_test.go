package nullable

import (
	"encoding/json"
	"errors"
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
	assertNullBool(t, null, "ValueFromPtr(nil)")
}

func TestUnmarshalBool(t *testing.T) {
	var b Nullable[bool]
	err := json.Unmarshal(boolJSON, &b)
	maybePanic(err)
	assertBool(t, b, "bool json")

	var nb Nullable[bool]
	err = json.Unmarshal(nullBoolJSON, &nb)
	if err == nil {
		panic("err should not be nil")
	}

	var null Nullable[bool]
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullBool(t, null, "null json")

	var badType Nullable[bool]
	err = json.Unmarshal(intJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, badType, "wrong type json")

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
	maybePanic(err)
	assertBool(t, b, "UnmarshalText() bool")

	var zero Nullable[bool]
	err = zero.UnmarshalText([]byte("false"))
	maybePanic(err)
	assertFalseBool(t, zero, "UnmarshalText() false")

	var blank Nullable[bool]
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBool(t, blank, "UnmarshalText() empty bool")

	var null Nullable[bool]
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullBool(t, null, `UnmarshalText() "null"`)

	var invalid Nullable[bool]
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, invalid, "invalid json")
}

func TestMarshalBool(t *testing.T) {
	b := Value(true)
	data, err := json.Marshal(b)
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty json marshal")

	zero := Value(false)
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero json marshal")

	// invalid values should be encoded as null
	null := Nullable[bool]{Value: false, HasValue: false}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalBoolText(t *testing.T) {
	b := Value(true)
	data, err := b.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty text marshal")

	zero := Value(false)
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero text marshal")

	// invalid values should be encoded as null
	null := Nullable[bool]{HasValue: false}
	data, err = null.MarshalText()
	maybePanic(err)
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
	assertBoolEqualIsTrue(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: false}
	b2 = Nullable[bool]{Value: false, HasValue: false}
	assertBoolEqualIsTrue(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: true}
	b2 = Nullable[bool]{Value: true, HasValue: true}
	assertBoolEqualIsTrue(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: true}
	b2 = Nullable[bool]{Value: true, HasValue: false}
	assertBoolEqualIsFalse(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: false}
	b2 = Nullable[bool]{Value: true, HasValue: true}
	assertBoolEqualIsFalse(t, b1, b2)

	b1 = Nullable[bool]{Value: true, HasValue: true}
	b2 = Nullable[bool]{Value: false, HasValue: true}
	assertBoolEqualIsFalse(t, b1, b2)
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

func assertNullBool(t *testing.T, b Nullable[bool], from string) {
	if b.HasValue {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertBoolEqualIsTrue(t *testing.T, a, b Nullable[bool]) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Value{%t, HasValue:%t} and Value{%t, HasValue:%t} should return true", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}

func assertBoolEqualIsFalse(t *testing.T, a, b Nullable[bool]) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Value{%t, HasValue:%t} and Value{%t, HasValue:%t} should return false", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}
