package nullable

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testProject struct {
	Id     int
	Guid   string
	Number string
	Notes  Nullable[string]
	Name   string

	ResponsibleCoworkerId Nullable[int]
	DepartmentId          Nullable[int]
}

var (
	nullJSON    = []byte(`null`)
	invalidJSON = []byte(`:)`)

	boolJSON     = []byte(`true`)
	nullBoolJSON = []byte(`{"Data":true,"IsValid":true}`)

	intJSON       = []byte(`12345`)
	intStringJSON = []byte(`"12345"`)
	nullIntJSON   = []byte(`{"Int64":12345,"Valid":true}`)

	floatJSON       = []byte(`1.2345`)
	floatStringJSON = []byte(`"1.2345"`)
	floatBlankJSON  = []byte(`""`)
	nullFloatJSON   = []byte(`{"Float64":1.2345,"Valid":true}`)

	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","IsValid":true}`)
)

func assertJSONEquals(t *testing.T, data []byte, cmp string, source string) {
	t.Helper()
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", source, data, cmp)
	}
}

func assertEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Nullable{\"%v\", Valid:%t} and Nullable{\"%v\", Valid:%t} should return true", a.Data, a.IsValid, b.Data, b.IsValid)
	}
}

func assertNotEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Nullable{\"%v\", Valid:%t} and Nullable{\"%v\", Valid:%t} should return false", a.Data, a.IsValid, b.Data, b.IsValid)
	}
}

func Test_Json_unmarshal(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\", \"Notes\": \"Some notes\", \"DepartmentId\": 12 }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.Nil(t, err)
	assert.Equal(t, 15, project.Id)
	assert.True(t, project.Notes.IsValid)
	assert.Equal(t, "Some notes", project.Notes.Data)
	assert.True(t, project.DepartmentId.IsValid)
	assert.Equal(t, 12, project.DepartmentId.Data)
}

func Test_Json_unmarshal_with_null(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\", \"Notes\": null, \"DepartmentId\": null }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.Nil(t, err)
	assert.Equal(t, 15, project.Id)
	assert.False(t, project.Notes.IsValid)
	assert.False(t, project.DepartmentId.IsValid)
}

func Test_Json_unmarshal_with_missing_values(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\" }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.Nil(t, err)
	assert.Equal(t, 15, project.Id)
	assert.False(t, project.Notes.IsValid)
	assert.False(t, project.DepartmentId.IsValid)
}

func Test_GoString(t *testing.T) {
	assert.Equal(t, `nullable.Nullable[bool]{Data:true,IsValid:true}`, Value(true).GoString())
	assert.Equal(t, `nullable.Nullable[bool]{Data:false,IsValid:false}`, Null[bool]().GoString())
	assert.Equal(t, `nullable.Nullable[float64]{Data:5.6,IsValid:true}`, Value(5.6).GoString())
	assert.Equal(t, `nullable.Nullable[float32]{Data:5.6,IsValid:true}`, Value(float32(5.6)).GoString())
	assert.Equal(t, `nullable.Nullable[float32]{Data:5.6,IsValid:true}`, Nullable[float32]{Data: 5.6, IsValid: true}.GoString())
	assert.Equal(t, `nullable.Nullable[int]{Data:5,IsValid:true}`, Value(5).GoString())
	assert.Equal(t, `nullable.Nullable[int8]{Data:5,IsValid:true}`, Value(int8(5)).GoString())
}
