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
)

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStr(t *testing.T, s Nullable[string], from string) {
	if s.Value != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.Value, "test")
	}
	if !s.HasValue {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s Nullable[string], from string) {
	if s.HasValue {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}

func assertStringEqualIsTrue(t *testing.T, a, b Nullable[string]) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of String{\"%v\", Valid:%t} and String{\"%v\", Valid:%t} should return true", a.Value, a.HasValue, b.Value, b.HasValue)
	}
}

func assertStringEqualIsFalse(t *testing.T, a, b Nullable[string]) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of String{\"%v\", Valid:%t} and String{\"%v\", Valid:%t} should return false", a.Value, a.HasValue, b.Value, b.HasValue)
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