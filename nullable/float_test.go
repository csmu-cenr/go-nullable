package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestFloatFrom(t *testing.T) {
	f := Value(1.2345)
	assertFloat(t, f, "Value()")

	zero := Value(0)
	if !zero.HasValue {
		t.Error("Value(0)", "is invalid, but should be valid")
	}
}

func TestFloatFromPtr(t *testing.T) {
	n := float64(1.2345)
	iptr := &n
	f := ValueFromPtr(iptr)
	assertFloat(t, f, "ValueFromPtr()")

	null := ValueFromPtr[float64](nil)
	assertNullFloat(t, null, "ValueFromPtr(nil)")
}

func TestUnmarshalFloat(t *testing.T) {
	var f Nullable[float64]
	err := json.Unmarshal(floatJSON, &f)
	assert.Nil(t, err)
	assertFloat(t, f, "float json")

	var sf Nullable[float64]
	err = json.Unmarshal(floatStringJSON, &sf)
	assert.Nil(t, err)
	assertFloat(t, sf, "string float json")

	var nf Nullable[float64]
	err = json.Unmarshal(nullFloatJSON, &nf)
	assert.Error(t, err)

	var null Nullable[float64]
	err = json.Unmarshal(nullJSON, &null)
	assert.Nil(t, err)
	assertNullFloat(t, null, "null json")

	var blank Nullable[float64]
	err = json.Unmarshal(floatBlankJSON, &blank)
	if err == nil {
		panic("expected error")
	}

	var badType Nullable[float64]
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}

	var invalid Nullable[float64]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalFloat(t *testing.T) {
	var f Nullable[float64]
	err := f.UnmarshalText([]byte("1.2345"))
	assert.Nil(t, err)
	assertFloat(t, f, "UnmarshalText() float")

	var blank Nullable[float64]
	err = blank.UnmarshalText([]byte(""))
	assert.Nil(t, err)
	assertNullFloat(t, blank, "UnmarshalText() empty float")

	var null Nullable[float64]
	err = null.UnmarshalText([]byte("null"))
	assert.Nil(t, err)
	assertNullFloat(t, null, `UnmarshalText() "null"`)

	var invalid Nullable[float64]
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalFloat(t *testing.T) {
	f := Value(1.2345)
	data, err := json.Marshal(f)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := Nullable[float64]{0, false}
	data, err = json.Marshal(null)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalFloatText(t *testing.T) {
	f := Value(1.2345)
	data, err := f.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := Nullable[float64]{0, false}
	data, err = null.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestFloatInfNaN(t *testing.T) {
	nan := Nullable[float64]{math.NaN(), true}
	_, err := nan.MarshalJSON()
	if err == nil {
		t.Error("expected error for NaN, got nil")
	}

	inf := Nullable[float64]{math.Inf(1), true}
	_, err = inf.MarshalJSON()
	if err == nil {
		t.Error("expected error for Inf, got nil")
	}
}

func TestFloatValueOrZero(t *testing.T) {
	valid := Nullable[float64]{1.2345, true}
	if valid.ValueOrZero() != 1.2345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[float64]{1.2345, false}
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestFloatEqual(t *testing.T) {
	f1 := Nullable[float64]{10, false}
	f2 := Nullable[float64]{10, false}
	assertFloatEqualIsTrue(t, f1, f2)

	f1 = Nullable[float64]{10, false}
	f2 = Nullable[float64]{20, false}
	assertFloatEqualIsTrue(t, f1, f2)

	f1 = Nullable[float64]{10, true}
	f2 = Nullable[float64]{10, true}
	assertFloatEqualIsTrue(t, f1, f2)

	f1 = Nullable[float64]{10, true}
	f2 = Nullable[float64]{10, false}
	assertFloatEqualIsFalse(t, f1, f2)

	f1 = Nullable[float64]{10, false}
	f2 = Nullable[float64]{10, true}
	assertFloatEqualIsFalse(t, f1, f2)

	f1 = Nullable[float64]{10, true}
	f2 = Nullable[float64]{20, true}
	assertFloatEqualIsFalse(t, f1, f2)
}

func assertFloat(t *testing.T, f Nullable[float64], from string) {
	if f.Value != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Value, 1.2345)
	}
	if !f.HasValue {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat(t *testing.T, f Nullable[float64], from string) {
	if f.HasValue {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertFloatEqualIsTrue(t *testing.T, a, b Nullable[float64]) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Nullable[float64{%v, HasValue:%t} and Nullable[float64{%v, HasValue:%t} should return true", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}

func assertFloatEqualIsFalse(t *testing.T, a, b Nullable[float64]) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Float{%v, HasValue:%t} and Float{%v, HasValue:%t} should return false", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}
