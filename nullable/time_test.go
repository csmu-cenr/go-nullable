package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
	timeObject    = []byte(`{"Time":"2012-12-21T21:21:21Z","HasValue":true}`)
	nullObject    = []byte(`{"Time":"0001-01-01T00:00:00Z","HasValue":false}`)
	badObject     = []byte(`{"hello": "world"}`)
)

func TestUnmarshalTimeJSON(t *testing.T) {
	var ti Nullable[time.Time]
	err := json.Unmarshal(timeJSON, &ti)
	assert.Nil(t, err)
	assertTime(t, ti, "UnmarshalJSON() json")

	var null Nullable[time.Time]
	err = json.Unmarshal(nullTimeJSON, &null)
	assert.Nil(t, err)
	assertNull(t, null, "null time json")

	var fromObject Nullable[time.Time]
	err = json.Unmarshal(timeObject, &fromObject)
	if err == nil {
		panic("expected error")
	}

	var nullFromObj Nullable[time.Time]
	err = json.Unmarshal(nullObject, &nullFromObj)
	if err == nil {
		panic("expected error")
	}

	var invalid Nullable[time.Time]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNull(t, invalid, "invalid from object json")

	var bad Nullable[time.Time]
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNull(t, bad, "bad from object json")

	var wrongType Nullable[time.Time]
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNull(t, wrongType, "wrong type object json")
}

func TestUnmarshalTimeText(t *testing.T) {
	ti := Value(timeValue1)
	txt, err := ti.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, txt, timeString1, "marshal text")

	var unmarshal Nullable[time.Time]
	err = unmarshal.UnmarshalText(txt)
	assert.Nil(t, err)
	assertTime(t, unmarshal, "unmarshal text")

	var null Nullable[time.Time]
	err = null.UnmarshalText(nullJSON)
	assert.Nil(t, err)
	assertNull(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, txt, "", "marshal null text")

	var invalid Nullable[time.Time]
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNull(t, invalid, "bad string")
}

func TestMarshalTime(t *testing.T) {
	ti := Value(timeValue1)
	data, err := json.Marshal(ti)
	assert.Nil(t, err)
	assertJSONEquals(t, data, string(timeJSON), "non-empty json marshal")

	ti.HasValue = false
	data, err = json.Marshal(ti)
	assert.Nil(t, err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestTimeFrom(t *testing.T) {
	ti := Value(timeValue1)
	assertTime(t, ti, "Value() time.Time")
}

func TestTimeFromPtr(t *testing.T) {
	ti := ValueFromPtr[time.Time](&timeValue1)
	assertTime(t, ti, "ValueFromPtr[time.Time() time")

	null := ValueFromPtr[time.Time](nil)
	assertNull(t, null, "ValueFromPtr[time.Time(nil)")
}

func TestTimeValueOrZero(t *testing.T) {
	valid := Value(timeValue1)
	if valid.ValueOrZero() != valid.Value || valid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := valid
	invalid.HasValue = false
	if !invalid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestTimeEqual(t *testing.T) {
	t1 := Nullable[time.Time]{timeValue1, false}
	t2 := Nullable[time.Time]{timeValue2, false}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false}
	t2 = Nullable[time.Time]{timeValue3, false}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue2, true}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue1, true}
	assertEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue2, false}
	assertNotEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false}
	t2 = Nullable[time.Time]{timeValue2, true}
	assertNotEqual(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue3, true}
	assertNotEqual(t, t1, t2)
}

func TestTimeExactEqual(t *testing.T) {
	t1 := Nullable[time.Time]{timeValue1, false}
	t2 := Nullable[time.Time]{timeValue1, false}
	assertTimeExactEqualIsTrue(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false}
	t2 = Nullable[time.Time]{timeValue2, false}
	assertTimeExactEqualIsTrue(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue1, true}
	assertTimeExactEqualIsTrue(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue1, false}
	assertTimeExactEqualIsFalse(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, false}
	t2 = Nullable[time.Time]{timeValue1, true}
	assertTimeExactEqualIsFalse(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue2, true}
	assertTimeExactEqualIsFalse(t, t1, t2)

	t1 = Nullable[time.Time]{timeValue1, true}
	t2 = Nullable[time.Time]{timeValue3, true}
	assertTimeExactEqualIsFalse(t, t1, t2)
}

func assertTime(t *testing.T, ti Nullable[time.Time], from string) {
	if ti.Value != timeValue1 {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Value, timeValue1)
	}
	if !ti.HasValue {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertTimeExactEqualIsTrue(t *testing.T, a, b Nullable[time.Time]) {
	t.Helper()
	if !a.ExactEqual(b) {
		t.Errorf("ExactEqual() of Time{%v, HasValue:%t} and Time{%v, HasValue:%t} should return true", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}

func assertTimeExactEqualIsFalse(t *testing.T, a, b Nullable[time.Time]) {
	t.Helper()
	if a.ExactEqual(b) {
		t.Errorf("ExactEqual() of Time{%v, HasValue:%t} and Time{%v, HasValue:%t} should return false", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}
