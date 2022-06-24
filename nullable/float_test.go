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
	assertFloat64(t, f, "Data()")

	zero := Value(0)
	if !zero.IsValid {
		t.Error("Data(0)", "is invalid, but should be valid")
	}
}

func TestFloatFromPointer(t *testing.T) {
	n := float64(1.2345)
	iPointer := &n
	f := ValueFromPointer(iPointer)
	assertFloat64(t, f, "ValueFromPointer()")

	null := ValueFromPointer[float64](nil)
	assert.False(t, null.IsValid)

	null32 := ValueFromPointer[float32](nil)
	assert.False(t, null32.IsValid)
}

func TestUnmarshalFloat(t *testing.T) {
	var f Nullable[float64]
	err := json.Unmarshal(floatJSON, &f)
	assert.Nil(t, err)
	assertFloat64(t, f, "float json")

	var sf Nullable[float64]
	err = json.Unmarshal(floatStringJSON, &sf)
	assert.Nil(t, err)
	assertFloat64(t, sf, "string float json")

	var nf Nullable[float64]
	err = json.Unmarshal(nullFloatJSON, &nf)
	assert.Error(t, err)

	var null Nullable[float64]
	err = json.Unmarshal(nullJSON, &null)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)

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

func TestUnmarshalFloat32(t *testing.T) {
	var f Nullable[float32]
	err := json.Unmarshal(floatJSON, &f)
	assert.Nil(t, err)
	assertFloat32(t, f, "float json")

	var sf Nullable[float32]
	err = json.Unmarshal(floatStringJSON, &sf)
	assert.Nil(t, err)
	assertFloat32(t, sf, "string float json")

	var nf Nullable[float32]
	err = json.Unmarshal(nullFloatJSON, &nf)
	assert.Error(t, err)

	var null Nullable[float32]
	err = json.Unmarshal(nullJSON, &null)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)

	var blank Nullable[float32]
	err = json.Unmarshal(floatBlankJSON, &blank)
	if err == nil {
		panic("expected error")
	}

	var badType Nullable[float32]
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}

	var invalid Nullable[float32]
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
	assertFloat64(t, f, "UnmarshalText() float")

	var blank Nullable[float64]
	err = blank.UnmarshalText([]byte(""))
	assert.Nil(t, err)
	assert.False(t, blank.IsValid)

	var null Nullable[float64]
	err = null.UnmarshalText([]byte("null"))
	assert.Nil(t, err)
	assert.False(t, null.IsValid)

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

func TestMarshalFloat32(t *testing.T) {
	f := Value[float32](1.2345)
	data, err := json.Marshal(f)
	assert.Nil(t, err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := Nullable[float32]{0, false}
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
	assertEqual(t, f1, f2)

	f1 = Nullable[float64]{10, false}
	f2 = Nullable[float64]{20, false}
	assertEqual(t, f1, f2)

	f1 = Nullable[float64]{10, true}
	f2 = Nullable[float64]{10, true}
	assertEqual(t, f1, f2)

	f1 = Nullable[float64]{10, true}
	f2 = Nullable[float64]{10, false}
	assertNotEqual(t, f1, f2)

	f1 = Nullable[float64]{10, false}
	f2 = Nullable[float64]{10, true}
	assertNotEqual(t, f1, f2)

	f1 = Nullable[float64]{10, true}
	f2 = Nullable[float64]{20, true}
	assertNotEqual(t, f1, f2)
}

func TestFloat64Scan(t *testing.T) {
	var f Nullable[float64]
	err := f.Scan(1.2345)
	assert.Nil(t, err)
	assertFloat64(t, f, "scanned float")

	var sf Nullable[float64]
	err = sf.Scan("1.2345")
	assert.Nil(t, err)
	assertFloat64(t, sf, "scanned string float")

	var null Nullable[float64]
	err = null.Scan(nil)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)
}

func TestFloat32Scan(t *testing.T) {
	var f Nullable[float32]
	err := f.Scan(1.2345)
	assert.Nil(t, err)
	assertFloat32(t, f, "scanned float")

	var sf Nullable[float32]
	err = sf.Scan("1.2345")
	assert.Nil(t, err)
	assertFloat32(t, sf, "scanned string float")

	var null Nullable[float32]
	err = null.Scan(nil)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)
}

func assertFloat64(t *testing.T, f Nullable[float64], from string) {
	if f.Data != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Data, 1.2345)
	}
	if !f.IsValid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFloat32(t *testing.T, f Nullable[float32], from string) {
	if f.Data != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Data, 1.2345)
	}
	if !f.IsValid {
		t.Error(from, "is invalid, but should be valid")
	}
}
