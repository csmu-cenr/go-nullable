package nullable

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	timeString1   = "2012-12-21T21:21:21Z"
	timeString2   = "2012-12-21T22:21:21+01:00" // Same time as timeString1 but in a different timezone
	timeString3   = "2018-08-19T01:02:03Z"
	timeJSON      = []byte(`"` + timeString1 + `"`)
	nullTimeJSON  = []byte(`null`)
	timeValue1, _ = time.Parse(time.RFC3339, timeString1)
	timeValue2, _ = time.Parse(time.RFC3339, timeString2)
	timeValue3, _ = time.Parse(time.RFC3339, timeString3)
	timeObject    = []byte(`{"Time":"2012-12-21T21:21:21Z","Valid":true}`)
	nullObject    = []byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`)
	badObject     = []byte(`{"hello": "world"}`)
)

func Test_Json_unmarshal_time(t *testing.T) {
	var ti Nullable[time.Time]
	err := json.Unmarshal(timeJSON, &ti)
	assert.NoError(t, err)
	assertTime(t, ti, "UnmarshalJSON() json")

	var null Nullable[time.Time]
	err = json.Unmarshal(nullTimeJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var fromObject Nullable[time.Time]
	err = json.Unmarshal(timeObject, &fromObject)
	assert.Error(t, err)

	var nullFromObj Nullable[time.Time]
	err = json.Unmarshal(nullObject, &nullFromObj)
	assert.Error(t, err)

	var invalid Nullable[time.Time]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assert.False(t, invalid.Valid)

	var bad Nullable[time.Time]
	err = json.Unmarshal(badObject, &bad)
	assert.Error(t, err)
	assert.False(t, bad.Valid)

	var wrongType Nullable[time.Time]
	err = json.Unmarshal(intJSON, &wrongType)
	assert.Error(t, err)
	assert.False(t, wrongType.Valid)
}

func Test_Text_unmarshal_time(t *testing.T) {
	ti := Value(timeValue1)
	txt, err := ti.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, timeString1, string(txt))

	var unmarshal Nullable[time.Time]
	err = unmarshal.UnmarshalText(txt)
	assert.NoError(t, err)
	assertTime(t, unmarshal, "unmarshal text")

	var null Nullable[time.Time]
	err = null.UnmarshalText(nullJSON)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
	txt, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(txt))

	var invalid Nullable[time.Time]
	err = invalid.UnmarshalText([]byte("hello world"))
	assert.Error(t, err)
	assert.False(t, invalid.Valid)
}

func Test_Json_marshal_time(t *testing.T) {
	ti := Value(timeValue1)
	data, err := json.Marshal(ti)
	assert.NoError(t, err)
	assert.Equal(t, timeJSON, data)

	ti.Valid = false
	data, err = json.Marshal(ti)
	assert.NoError(t, err)
	assert.Equal(t, nullJSON, data)
}

func Test_Time_from_value(t *testing.T) {
	ti := Value(timeValue1)
	assertTime(t, ti, "Data() time.Time")
}

func Test_Time_from_pointer(t *testing.T) {
	ti := ValueFromPointer[time.Time](&timeValue1)
	assertTime(t, ti, "ValueFromPointer[time.Time() time")

	null := ValueFromPointer[time.Time](nil)
	assert.False(t, null.Valid)
}

func Test_Time_ValueOrZero(t *testing.T) {
	valid := Value(timeValue1)
	if valid.ValueOrZero() != valid.Data || valid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := valid
	invalid.Valid = false
	if !invalid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func Test_Time_Equal(t *testing.T) {
	t1 := Nullable[time.Time]{timeValue1, false, false}
	t2 := Nullable[time.Time]{timeValue2, false, false}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false, false}
	t2 = Nullable[time.Time]{timeValue3, false, false}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue2, true, false}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue1, true, false}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue2, false, false}
	assertNotEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false, false}
	t2 = Nullable[time.Time]{timeValue2, true, false}
	assertNotEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue3, true, false}
	assertNotEqual(t, t1, t2)
}

func Test_Time_ExactEqual(t *testing.T) {
	t1 := Nullable[time.Time]{timeValue1, false, false}
	t2 := Nullable[time.Time]{timeValue1, false, false}
	assertExactEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false, false}
	t2 = Nullable[time.Time]{timeValue2, false, false}
	assertExactEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue1, true, false}
	assertExactEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue1, false, false}
	assertNotExactEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false, false}
	t2 = Nullable[time.Time]{timeValue1, true, false}
	assertNotExactEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue2, true, false}
	assertNotExactEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true, false}
	t2 = Nullable[time.Time]{timeValue3, true, false}
	assertNotExactEqual(t, t1, t2)
}

func Test_Time_Scan(t *testing.T) {
	var ti Nullable[time.Time]
	err := ti.Scan(timeValue1)
	assert.NoError(t, err)
	assertTime(t, ti, "scanned time")
	if v, err := ti.Value(); v != timeValue1 || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null Nullable[time.Time]
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong Nullable[time.Time]
	err = wrong.Scan(int64(42))
	assert.NotNil(t, err)
}

func Test_IsZero_time(t *testing.T) {
	var tm Nullable[time.Time]
	assert.True(t, tm.IsZero())

	var zeroTime time.Time
	tm = Value(zeroTime)
	assert.True(t, tm.IsZero())

	tm = Value(time.Now())
	assert.False(t, tm.IsZero())
}

func assertTime(t *testing.T, ti Nullable[time.Time], from string) {
	if ti.Data != timeValue1 {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Data, timeValue1)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertExactEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	if !a.ExactEqual(b) {
		t.Errorf("ExactEqual() of Time{%v, Valid:%t} and Time{%v, Valid:%t} should return true", a.Data, a.Valid, b.Data, b.Valid)
	}
}

func assertNotExactEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	if a.ExactEqual(b) {
		t.Errorf("ExactEqual() of Time{%v, Valid:%t} and Time{%v, Valid:%t} should return false", a.Data, a.Valid, b.Data, b.Valid)
	}
}
