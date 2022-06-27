package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_Float_from_value(t *testing.T) {
	f := Value(1.2345)
	assertFloat64(t, f, "Data()")

	zero := Value(0)
	if !zero.Valid {
		t.Error("Data(0)", "is invalid, but should be valid")
	}
}

func Test_Float_from_pointer(t *testing.T) {
	n := 1.2345
	iPointer := &n
	f := ValueFromPointer(iPointer)
	assertFloat64(t, f, "ValueFromPointer()")

	null := ValueFromPointer[float64](nil)
	assert.False(t, null.Valid)

	null32 := ValueFromPointer[float32](nil)
	assert.False(t, null32.Valid)
}

func Test_Json_unmarshal_float64(t *testing.T) {
	var f Nullable[float64]
	err := json.Unmarshal(floatJSON, &f)
	assert.NoError(t, err)
	assertFloat64(t, f, "float json")

	var sf Nullable[float64]
	err = json.Unmarshal(floatStringJSON, &sf)
	assert.NoError(t, err)
	assertFloat64(t, sf, "string float json")

	var nf Nullable[float64]
	err = json.Unmarshal(nullFloatJSON, &nf)
	assert.Error(t, err)

	var null Nullable[float64]
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var blank Nullable[float64]
	err = json.Unmarshal(floatBlankJSON, &blank)
	assert.Error(t, err)

	var badType Nullable[float64]
	err = json.Unmarshal(boolJSON, &badType)
	assert.Error(t, err)

	var invalid Nullable[float64]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
}

func Test_Json_unmarshal_float32(t *testing.T) {
	var f Nullable[float32]
	err := json.Unmarshal(floatJSON, &f)
	assert.NoError(t, err)
	assertFloat32(t, f, "float json")

	var sf Nullable[float32]
	err = json.Unmarshal(floatStringJSON, &sf)
	assert.NoError(t, err)
	assertFloat32(t, sf, "string float json")

	var nf Nullable[float32]
	err = json.Unmarshal(nullFloatJSON, &nf)
	assert.Error(t, err)

	var null Nullable[float32]
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var blank Nullable[float32]
	err = json.Unmarshal(floatBlankJSON, &blank)
	assert.NotNil(t, err)

	var badType Nullable[float32]
	err = json.Unmarshal(boolJSON, &badType)
	assert.NotNil(t, err)

	var invalid Nullable[float32]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
}

func Test_Text_unmarshal_float64(t *testing.T) {
	var f Nullable[float64]
	err := f.UnmarshalText([]byte("1.2345"))
	assert.NoError(t, err)
	assertFloat64(t, f, "UnmarshalText() float")

	var blank Nullable[float64]
	err = blank.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assert.False(t, blank.Valid)

	var null Nullable[float64]
	err = null.UnmarshalText([]byte("null"))
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var invalid Nullable[float64]
	err = invalid.UnmarshalText([]byte("hello world"))
	assert.Error(t, err)
}

func Test_Json_marshal_float(t *testing.T) {
	f := Value(1.2345)
	data, err := json.Marshal(f)
	assert.NoError(t, err)
	assert.Equal(t, "1.2345", string(data))

	// invalid values should be encoded as null
	null := Nullable[float64]{0, false}
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func Test_Json_marshal_float32(t *testing.T) {
	f := Value[float32](1.2345)
	data, err := json.Marshal(f)
	assert.NoError(t, err)
	assert.Equal(t, "1.2345", string(data))

	// invalid values should be encoded as null
	null := Nullable[float32]{0, false}
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(data))
}

func Test_Text_marshal_float64(t *testing.T) {
	f := Value(1.2345)
	data, err := f.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "1.2345", string(data))

	// invalid values should be encoded as null
	null := Nullable[float64]{0, false}
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func Test_Float_Inf_and_NaN(t *testing.T) {
	nan := Nullable[float64]{math.NaN(), true}
	_, err := nan.MarshalJSON()
	assert.Error(t, err, "expected error for NaN")

	inf := Nullable[float64]{math.Inf(1), true}
	_, err = inf.MarshalJSON()
	assert.Error(t, err, "expected error for Inf")
}

func Test_Float_ValueOrZero(t *testing.T) {
	valid := Nullable[float64]{1.2345, true}
	if valid.ValueOrZero() != 1.2345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[float64]{1.2345, false}
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func Test_Float_Equal(t *testing.T) {
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

func Test_Float64_Scan(t *testing.T) {
	var f Nullable[float64]
	err := f.Scan(1.2345)
	assert.NoError(t, err)
	assertFloat64(t, f, "scanned float")

	var sf Nullable[float64]
	err = sf.Scan("1.2345")
	assert.NoError(t, err)
	assertFloat64(t, sf, "scanned string float")

	var null Nullable[float64]
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func Test_Float32_Scan(t *testing.T) {
	var f Nullable[float32]
	err := f.Scan(1.2345)
	assert.NoError(t, err)
	assertFloat32(t, f, "scanned float")

	var sf Nullable[float32]
	err = sf.Scan("1.2345")
	assert.NoError(t, err)
	assertFloat32(t, sf, "scanned string float")

	var null Nullable[float32]
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func Test_IsZero_float(t *testing.T) {
	var f Nullable[float64]
	assert.True(t, f.IsZero())

	var zeroFloat float64
	f = Value(zeroFloat)
	assert.True(t, f.IsZero())

	f = Value(1.0)
	assert.False(t, f.IsZero())
}

func assertFloat64(t *testing.T, f Nullable[float64], from string) {
	if f.Data != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Data, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFloat32(t *testing.T, f Nullable[float32], from string) {
	if f.Data != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Data, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}
