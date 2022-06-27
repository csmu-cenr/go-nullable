package nullable

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stringInStruct struct {
	Test Nullable[string] `json:"test,omitempty"`
}

func Test_String_from_value(t *testing.T) {
	str := Value("test")
	assertStr(t, str, "Data() string")

	zero := Value("")
	if !zero.Valid {
		t.Error("Data(0)", "is invalid, but should be valid")
	}
}

func Test_String_from_pointer(t *testing.T) {
	s := "test"
	sPointer := &s
	str := ValueFromPointer(sPointer)
	assertStr(t, str, "ValueFromPointer() string")

	null := ValueFromPointer[string](nil)
	assert.False(t, null.Valid)
}

func Test_Json_unmarshal_string(t *testing.T) {
	var str Nullable[string]
	err := json.Unmarshal(stringJSON, &str)
	assert.NoError(t, err)
	assertStr(t, str, "string json")

	var ns Nullable[string]
	err = json.Unmarshal(nullStringJSON, &ns)
	assert.Error(t, err)

	var blank Nullable[string]
	err = json.Unmarshal(blankStringJSON, &blank)
	assert.NoError(t, err)
	assert.True(t, blank.Valid)

	var null Nullable[string]
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err)
	assert.False(t, null.Valid)

	var badType Nullable[string]
	err = json.Unmarshal(boolJSON, &badType)
	assert.Error(t, err)
	assert.False(t, badType.Valid)

	var invalid Nullable[string]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assert.False(t, invalid.Valid)
}

func Test_Text_unmarshal_string(t *testing.T) {
	var str Nullable[string]
	err := str.UnmarshalText([]byte("test"))
	assert.NoError(t, err)
	assertStr(t, str, "UnmarshalText() string")

	var null Nullable[string]
	err = null.UnmarshalText([]byte(""))
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func Test_Json_marshal_string(t *testing.T) {
	str := Value("test")
	data, err := json.Marshal(str)
	assert.NoError(t, err)
	assert.Equal(t, `"test"`, string(data))
	data, err = str.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))

	// empty values should be encoded as an empty string
	zero := Value("")
	data, err = json.Marshal(zero)
	assert.NoError(t, err)
	assert.Equal(t, `""`, string(data))
	data, err = zero.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))

	null := ValueFromPointer[string](nil)
	data, err = json.Marshal(null)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(data))
	data, err = null.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func Test_Json_marshal_string_in_struct(t *testing.T) {
	obj := stringInStruct{Test: Value("")}
	data, err := json.Marshal(obj)
	assert.NoError(t, err)
	assertJSONEquals(t, `{"test":""}`, data, "null string in struct")

	obj = stringInStruct{Test: Nullable[string]{}}
	data, err = json.Marshal(obj)
	assert.NoError(t, err)
	assertJSONEquals(t, `{"test":null}`, data, "null string in struct")
}

func Test_String_ValueOrZero(t *testing.T) {
	valid := Nullable[string]{"test", true}
	if valid.ValueOrZero() != "test" {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[string]{"test", false}
	if invalid.ValueOrZero() != "" {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func Test_String_equal(t *testing.T) {
	str1 := Nullable[string]{"foo", false}
	str2 := Nullable[string]{"foo", false}
	assertEqual(t, str1, str2)

	str1 = Nullable[string]{"foo", false}
	str2 = Nullable[string]{"bar", false}
	assertEqual(t, str1, str2)

	str1 = Nullable[string]{"foo", true}
	str2 = Nullable[string]{"foo", true}
	assertEqual(t, str1, str2)

	str1 = Nullable[string]{"foo", true}
	str2 = Nullable[string]{"foo", false}
	assertNotEqual(t, str1, str2)

	str1 = Nullable[string]{"foo", false}
	str2 = Nullable[string]{"foo", true}
	assertNotEqual(t, str1, str2)

	str1 = Nullable[string]{"foo", true}
	str2 = Nullable[string]{"bar", true}
	assertNotEqual(t, str1, str2)
}

func Test_String_scan(t *testing.T) {
	var str Nullable[string]
	err := str.Scan("test")
	assert.NoError(t, err)
	assertStr(t, str, "scanned string")

	var null Nullable[string]
	err = null.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, null.Valid)
}

func Test_IsZero_string(t *testing.T) {
	var str Nullable[string]
	assert.True(t, str.IsZero())

	var zeroStr string
	str = Value(zeroStr)
	assert.True(t, str.IsZero())

	str = Value("asdf")
	assert.False(t, str.IsZero())
}

func assertStr(t *testing.T, s Nullable[string], from string) {
	if s.Data != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.Data, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}
