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
	falseJSON    = []byte(`false`)
	nullBoolJSON = []byte(`{"Value":true,"HasValue":true}`)

	intJSON       = []byte(`12345`)
	intStringJSON = []byte(`"12345"`)
	nullIntJSON   = []byte(`{"Int64":12345,"Valid":true}`)

	floatJSON       = []byte(`1.2345`)
	floatStringJSON = []byte(`"1.2345"`)
	floatBlankJSON  = []byte(`""`)
	nullFloatJSON   = []byte(`{"Float64":1.2345,"Valid":true}`)
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func assertNull[T any](t *testing.T, f Nullable[T], from string) {
	if f.HasValue {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}

func assertEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Nullable{\"%v\", Valid:%t} and Nullable{\"%v\", Valid:%t} should return true", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}

func assertNotEqual[T any](t *testing.T, a, b Nullable[T]) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Nullable{\"%v\", Valid:%t} and Nullable{\"%v\", Valid:%t} should return false", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}

func Test_Json_unmarshal(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\", \"Notes\": \"Some notes\", \"DepartmentId\": 12 }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.Nil(t, err)
	assert.Equal(t, 15, project.Id)
	assert.True(t, project.Notes.HasValue)
	assert.Equal(t, "Some notes", project.Notes.Value)
	assert.True(t, project.DepartmentId.HasValue)
	assert.Equal(t, 12, project.DepartmentId.Value)
}

func Test_Json_unmarshal_with_null(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\", \"Notes\": null, \"DepartmentId\": null }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.Nil(t, err)
	assert.Equal(t, 15, project.Id)
	assert.False(t, project.Notes.HasValue)
	assert.False(t, project.DepartmentId.HasValue)
}

func Test_Json_unmarshal_with_missing_values(t *testing.T) {
	var project testProject
	jsonData := "{ \"Id\": 15, \"Number\": \"1234\" }"
	err := json.Unmarshal([]byte(jsonData), &project)
	assert.Nil(t, err)
	assert.Equal(t, 15, project.Id)
	assert.False(t, project.Notes.HasValue)
	assert.False(t, project.DepartmentId.HasValue)
}
