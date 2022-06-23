package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"testing"
)

func TestIntFrom(t *testing.T) {
	i := Value(12345)
	assertInt(t, i, "Value()")

	zero := Value(0)
	if !zero.HasValue {
		t.Error("Value(0)", "is invalid, but should be valid")
	}
}

func TestIntFromPtr(t *testing.T) {
	n := 12345
	intPointer := &n
	i := ValueFromPtr(intPointer)
	assertInt(t, i, "ValueFromPtr()")

	null := ValueFromPtr[int](nil)
	assertNull(t, null, "ValueFromPtr(nil)")
}

func TestUnmarshalInt(t *testing.T) {
	var i Nullable[int]
	err := json.Unmarshal(intJSON, &i)
	assert.Nil(t, err)
	assertInt(t, i, "int json")

	var si Nullable[int]
	err = json.Unmarshal(intStringJSON, &si)
	assert.Nil(t, err)
	assertInt(t, si, "int string json")

	var ni Nullable[int]
	err = json.Unmarshal(nullIntJSON, &ni)
	if err == nil {
		panic("err should not be nill")
	}

	var bi Nullable[int]
	err = json.Unmarshal(floatBlankJSON, &bi)
	if err == nil {
		panic("err should not be nill")
	}

	var null Nullable[int]
	err = json.Unmarshal(nullJSON, &null)
	assert.Nil(t, err)
	assertNull(t, null, "null json")

	var badType Nullable[int]
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNull(t, badType, "wrong type json")

	var invalid Nullable[int]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNull(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i Nullable[int]
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Nullable[int]
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	assert.Nil(t, err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalInt(t *testing.T) {
	var i Nullable[int]
	err := i.UnmarshalText([]byte("12345"))
	assert.Nil(t, err)
	assertInt(t, i, "UnmarshalText() int")

	var blank Nullable[int]
	err = blank.UnmarshalText([]byte(""))
	assert.Nil(t, err)
	assertNull(t, blank, "UnmarshalText() empty int")

	var null Nullable[int]
	err = null.UnmarshalText([]byte("null"))
	assert.Nil(t, err)
	assertNull(t, null, `UnmarshalText() "null"`)

	var invalid Nullable[int]
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalInt(t *testing.T) {
	i := Value(12345)
	data, err := json.Marshal(i)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := Nullable[int]{0, false}
	data, err = json.Marshal(null)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIntText(t *testing.T) {
	i := Value(12345)
	data, err := i.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := Nullable[int]{0, false}
	data, err = null.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestIntValueOrZero(t *testing.T) {
	valid := Nullable[int]{12345, true}
	if valid.ValueOrZero() != 12345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[int]{12345, false}
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestIntEqual(t *testing.T) {
	int1 := Nullable[int]{10, false}
	int2 := Nullable[int]{10, false}
	assertEqual(t, int1, int2)

	int1 = Nullable[int]{10, false}
	int2 = Nullable[int]{20, false}
	assertEqual(t, int1, int2)

	int1 = Nullable[int]{10, true}
	int2 = Nullable[int]{10, true}
	assertEqual(t, int1, int2)

	int1 = Nullable[int]{10, true}
	int2 = Nullable[int]{10, false}
	assertNotEqual(t, int1, int2)

	int1 = Nullable[int]{10, false}
	int2 = Nullable[int]{10, true}
	assertNotEqual(t, int1, int2)

	int1 = Nullable[int]{10, true}
	int2 = Nullable[int]{20, true}
	assertNotEqual(t, int1, int2)
}

func assertInt(t *testing.T, i Nullable[int], from string) {
	if i.Value != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Value, 12345)
	}
	if !i.HasValue {
		t.Error(from, "is invalid, but should be valid")
	}
}
