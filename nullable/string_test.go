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

func TestStringFrom(t *testing.T) {
	str := Value("test")
	assertStr(t, str, "Data() string")

	zero := Value("")
	if !zero.IsValid {
		t.Error("Data(0)", "is invalid, but should be valid")
	}
}

func TestStringFromPointer(t *testing.T) {
	s := "test"
	sPointer := &s
	str := ValueFromPointer(sPointer)
	assertStr(t, str, "ValueFromPointer() string")

	null := ValueFromPointer[string](nil)
	assert.False(t, null.IsValid)
}

func TestUnmarshalString(t *testing.T) {
	var str Nullable[string]
	err := json.Unmarshal(stringJSON, &str)
	assert.Nil(t, err)
	assertStr(t, str, "string json")

	var ns Nullable[string]
	err = json.Unmarshal(nullStringJSON, &ns)
	if err == nil {
		panic("err should not be nil")
	}

	var blank Nullable[string]
	err = json.Unmarshal(blankStringJSON, &blank)
	assert.Nil(t, err)
	if !blank.IsValid {
		t.Error("blank string should be valid")
	}

	var null Nullable[string]
	err = json.Unmarshal(nullJSON, &null)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)

	var badType Nullable[string]
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assert.False(t, badType.IsValid)

	var invalid Nullable[string]
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assert.False(t, invalid.IsValid)
}

func TestTextUnmarshalString(t *testing.T) {
	var str Nullable[string]
	err := str.UnmarshalText([]byte("test"))
	assert.Nil(t, err)
	assertStr(t, str, "UnmarshalText() string")

	var null Nullable[string]
	err = null.UnmarshalText([]byte(""))
	assert.Nil(t, err)
	assert.False(t, null.IsValid)
}

func TestMarshalString(t *testing.T) {
	str := Value("test")
	data, err := json.Marshal(str)
	assert.Nil(t, err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := Value("")
	data, err = json.Marshal(zero)
	assert.Nil(t, err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "", "string marshal text")

	null := ValueFromPointer[string](nil)
	data, err = json.Marshal(null)
	assert.Nil(t, err)
	assertJSONEquals(t, data, `null`, "null json marshal")
	data, err = null.MarshalText()
	assert.Nil(t, err)
	assertJSONEquals(t, data, "", "string marshal text")
}

func TestMarshalStringInStruct(t *testing.T) {
	obj := stringInStruct{Test: Value("")}
	data, err := json.Marshal(obj)
	assert.Nil(t, err)
	assertJSONEquals(t, data, `{"test":""}`, "null string in struct")

	obj = stringInStruct{Test: Nullable[string]{}}
	data, err = json.Marshal(obj)
	assert.Nil(t, err)
	assertJSONEquals(t, data, `{"test":null}`, "null string in struct")
}

func TestStringValueOrZero(t *testing.T) {
	valid := Nullable[string]{"test", true}
	if valid.ValueOrZero() != "test" {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := Nullable[string]{"test", false}
	if invalid.ValueOrZero() != "" {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestStringEqual(t *testing.T) {
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

func TestStringScan(t *testing.T) {
	var str Nullable[string]
	err := str.Scan("test")
	assert.Nil(t, err)
	assertStr(t, str, "scanned string")

	var null Nullable[string]
	err = null.Scan(nil)
	assert.Nil(t, err)
	assert.False(t, null.IsValid)
}

func assertStr(t *testing.T, s Nullable[string], from string) {
	if s.Data != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.Data, "test")
	}
	if !s.IsValid {
		t.Error(from, "is invalid, but should be valid")
	}
}
