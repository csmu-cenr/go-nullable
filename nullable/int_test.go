package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"testing"
)

func Test_Int_from_value(t *testing.T) {
	i := Value(12345)
	assertIntValue(t, i, "Data()")

	zero := Value(0)
	if !zero.Valid {
		t.Error("Data(0)", "is invalid, but should be valid")
	}
}

func Test_Int_from_pointer(t *testing.T) {
	n := 12345
	intPointer := &n
	i := ValueFromPointer(intPointer)
	assertIntValue(t, i, "ValueFromPointer()")

	null := ValueFromPointer[int](nil)
	assert.False(t, null.Valid)
}

func Test_Json_unmarshal_int(t *testing.T) {
	var i Nullable[int]
	err := json.Unmarshal(intJSON, &i)
	assert.NoError(t, err)
	assertIntValue(t, i, "int json")

	var si Nullable[int]
	err = json.Unmarshal(intStringJSON, &si)
	assert.NoError(t, err)
	assertIntValue(t, si, "int string json")

	var ni Nullable[int]
	err = json.Unmarshal(nullIntJSON, &ni)
	assert.Error(t, err)

	var bi Nullable[int]
	err = json.Unmarshal(floatBlankJSON, &bi)
	assert.Error(t, err)

	var null Nullable[int]
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var badType Nullable[int]
	err = json.Unmarshal(boolJSON, &badType)
	assert.Error(t, err)
	assert.False(t, badType.Valid)

	var invalid Nullable[int]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assert.False(t, invalid.Valid)
}

func Test_Json_unmarshal_non_integer_number(t *testing.T) {
	var i Nullable[int]
	err := json.Unmarshal(floatJSON, &i)
	assert.Error(t, err, "err should be present; non-integer number coerced to int")
}

func Test_Json_unmarshal_int64_overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Nullable[int]
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	assert.NoError(t, err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	assert.Error(t, err, "err should be present; decoded value overflows int64")
}

func Test_Text_unmarshal_int(t *testing.T) {
	var i Nullable[int]
	err := i.UnmarshalText([]byte("12345"))
	assert.NoError(t, err)
	assertIntValue(t, i, "UnmarshalText() int")

	var blank Nullable[int]
	err = blank.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assert.False(t, blank.Valid)

	var null Nullable[int]
	err = null.UnmarshalText([]byte("null"))
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var invalid Nullable[int]
	err = invalid.UnmarshalText([]byte("hello world"))
	assert.Error(t, err)
}

func Test_Json_marshal_int(t *testing.T) {
	i := Value(12345)
	data, err := json.Marshal(i)
	assert.NoError(t, err)
	assert.Equal(t, "12345", string(data))

	// invalid values should be encoded as null
	null := Nullable[int]{0, false}
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func Test_Text_marshal_int(t *testing.T) {
	i := Value(12345)
	data, err := i.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "12345", string(data))

	// invalid values should be encoded as null
	null := Nullable[int]{0, false}
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func Test_Int_ValueOrZero(t *testing.T) {
	valid := Nullable[int]{12345, true}
	if valid.ValueOrZero() != 12345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[int]{12345, false}
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func Test_Int_Equal(t *testing.T) {
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

func Test_Int_Scan(t *testing.T) {
	var i Nullable[int]
	err := i.Scan(12345)
	assert.NoError(t, err)
	assertIntValue(t, i, "scanned valid int")

	var null Nullable[int]
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func Test_IsZero_int(t *testing.T) {
	var i Nullable[int]
	assert.True(t, i.IsZero())

	var zeroInt int
	i = Value(zeroInt)
	assert.True(t, i.IsZero())

	i = Value(1)
	assert.False(t, i.IsZero())
}

func assertIntValue(t *testing.T, i Nullable[int], source string) {
	t.Helper()
	if i.Data != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", source, i.Data, 12345)
	}
	if !i.Valid {
		t.Error(source, "should be valid")
	}
}
