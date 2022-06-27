package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Bool_from_value(t *testing.T) {
	b := Value(true)
	assertBool(t, b, "Data()")

	zero := Value(false)
	if !zero.Valid {
		t.Error("Data(false)", "is invalid, but should be valid")
	}
}

func Test_Bool_from_pointer(t *testing.T) {
	n := true
	boolPointer := &n
	b := ValueFromPointer(boolPointer)
	assertBool(t, b, "ValueFromPointer()")

	null := ValueFromPointer[bool](nil)
	assert.False(t, null.Valid)
}

func Test_Json_unmarshal_bool(t *testing.T) {
	var b Nullable[bool]
	err := json.Unmarshal(boolJSON, &b)
	assert.NoError(t, err)
	assertBool(t, b, "bool json")

	var nb Nullable[bool]
	err = json.Unmarshal(nullBoolJSON, &nb)
	assert.NotNil(t, err)

	var null Nullable[bool]
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var badType Nullable[bool]
	err = json.Unmarshal(intJSON, &badType)
	assert.NotNil(t, err)
	assert.False(t, badType.Valid)

	var invalid Nullable[bool]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
}

func Test_Text_unmarshal_bool(t *testing.T) {
	var b Nullable[bool]
	err := b.UnmarshalText([]byte("true"))
	assert.NoError(t, err)
	assertBool(t, b, "UnmarshalText() bool")

	var falseBool Nullable[bool]
	err = falseBool.UnmarshalText([]byte("false"))
	assert.NoError(t, err)
	assertFalseBool(t, falseBool, "UnmarshalText() false")

	var fromBlankString Nullable[bool]
	err = fromBlankString.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assert.False(t, fromBlankString.Valid)

	var fromNullString Nullable[bool]
	err = fromNullString.UnmarshalText([]byte("null"))
	assert.NoError(t, err)
	assert.False(t, fromNullString.Valid)

	var invalid Nullable[bool]
	err = invalid.UnmarshalText([]byte(":D"))
	assert.Error(t, err)
	assert.False(t, invalid.Valid)
}

func Test_Json_marshal_bool(t *testing.T) {
	b := Value(true)
	data, err := json.Marshal(b)
	assert.NoError(t, err)
	assert.Equal(t, "true", string(data))

	zero := Value(false)
	data, err = json.Marshal(zero)
	assert.NoError(t, err)
	assert.Equal(t, "false", string(data))

	// invalid values should be encoded as null
	null := Null[bool]()
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func Test_Text_marshal_bool(t *testing.T) {
	b := Value(true)
	data, err := b.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "true", string(data))

	zero := Value(false)
	data, err = zero.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "false", string(data))

	// invalid values should be encoded as null
	null := Null[bool]()
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func Test_Bool_ValueOrZero(t *testing.T) {
	valid := Value(true)
	if valid.ValueOrZero() != true {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[bool]{Data: true, Valid: false}
	if invalid.ValueOrZero() != false {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func Test_Bool_Equal(t *testing.T) {
	b1 := Nullable[bool]{Data: true, Valid: false}
	b2 := Nullable[bool]{Data: true, Valid: false}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, Valid: false}
	b2 = Nullable[bool]{Data: false, Valid: false}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, Valid: true}
	b2 = Nullable[bool]{Data: true, Valid: true}
	assertEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, Valid: true}
	b2 = Nullable[bool]{Data: true, Valid: false}
	assertNotEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, Valid: false}
	b2 = Nullable[bool]{Data: true, Valid: true}
	assertNotEqual(t, b1, b2)

	b1 = Nullable[bool]{Data: true, Valid: true}
	b2 = Nullable[bool]{Data: false, Valid: true}
	assertNotEqual(t, b1, b2)
}

func Test_Bool_Scan(t *testing.T) {
	var b Nullable[bool]
	err := b.Scan(true)
	assert.NoError(t, err)
	assertBool(t, b, "scanned bool")

	var null Nullable[bool]
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func Test_IsZero_bool(t *testing.T) {
	var b Nullable[bool]
	assert.True(t, b.IsZero())

	var zeroBool bool
	b = Value(zeroBool)
	assert.True(t, b.IsZero())

	b = Value(true)
	assert.False(t, b.IsZero())
}

func assertBool(t *testing.T, b Nullable[bool], source string) {
	if b.Data != true {
		t.Errorf("bad %s bool: %v ≠ %v\n", source, b.Data, true)
	}
	if !b.Valid {
		t.Error(source, "is invalid, but should be valid")
	}
}

func assertFalseBool(t *testing.T, b Nullable[bool], from string) {
	if b.Data != false {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Data, false)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}
