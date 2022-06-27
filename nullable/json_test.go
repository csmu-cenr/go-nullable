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
	nullBoolJSON = []byte(`{"Data":true,"Valid":true}`)

	intJSON       = []byte(`12345`)
	intStringJSON = []byte(`"12345"`)
	nullIntJSON   = []byte(`{"Int64":12345,"Valid":true}`)

	floatJSON       = []byte(`1.2345`)
	floatStringJSON = []byte(`"1.2345"`)
	floatBlankJSON  = []byte(`""`)
	nullFloatJSON   = []byte(`{"Float64":1.2345,"Valid":true}`)

	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)
)

func assertJSONEquals(t *testing.T, expected string, actual []byte, source string) {
	t.Helper()

	var actualDecoded map[string]interface{}
	var expectedDecoded map[string]interface{}

	err := json.Unmarshal(actual, &actualDecoded)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(expected), &expectedDecoded)
	assert.NoError(t, err)
	assert.Equal(t, expectedDecoded, actualDecoded, "bad %s data: %s ≠ %s\n", source, actual, expected)
}

func assertEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	assert.True(t, a.Equal(b), "Equal() of Nullable{\"%v\", Valid:%t} and Nullable{\"%v\", Valid:%t} should return true", a.Data, a.Valid, b.Data, b.Valid)
}

func assertNotEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	assert.False(t, a.Equal(b), "Equal() of Nullable{\"%v\", Valid:%t} and Nullable{\"%v\", Valid:%t} should return false", a.Data, a.Valid, b.Data, b.Valid)
}

func Test_Json_unmarshal(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\", \"Notes\": \"Some notes\", \"DepartmentId\": 12 }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.NoError(t, err)
	assert.Equal(t, 15, project.Id)
	assert.True(t, project.Notes.Valid)
	assert.Equal(t, "Some notes", project.Notes.Data)
	assert.True(t, project.DepartmentId.Valid)
	assert.Equal(t, 12, project.DepartmentId.Data)
}

func Test_Json_unmarshal_with_null(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\", \"Notes\": null, \"DepartmentId\": null }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.NoError(t, err)
	assert.Equal(t, 15, project.Id)
	assert.False(t, project.Notes.Valid)
	assert.False(t, project.DepartmentId.Valid)
}

func Test_Json_unmarshal_with_missing_values(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\" }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.NoError(t, err)
	assert.Equal(t, 15, project.Id)
	assert.False(t, project.Notes.Valid)
	assert.False(t, project.DepartmentId.Valid)
}

func Test_GoString(t *testing.T) {
	assert.Equal(t, `nullable.Nullable[bool]{Data:true,Valid:true}`, Value(true).GoString())
	assert.Equal(t, `nullable.Nullable[bool]{Data:false,Valid:false}`, Null[bool]().GoString())
	assert.Equal(t, `nullable.Nullable[float64]{Data:5.6,Valid:true}`, Value(5.6).GoString())
	assert.Equal(t, `nullable.Nullable[float32]{Data:5.6,Valid:true}`, Value(float32(5.6)).GoString())
	assert.Equal(t, `nullable.Nullable[float32]{Data:5.6,Valid:true}`, Nullable[float32]{Data: 5.6, Valid: true}.GoString())
	assert.Equal(t, `nullable.Nullable[int]{Data:5,Valid:true}`, Value(5).GoString())
	assert.Equal(t, `nullable.Nullable[int8]{Data:5,Valid:true}`, Value(int8(5)).GoString())
}
