package nullable

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type task struct {
	TaskId     int           `json:"task_id"`
	ProjectId  Nullable[int] `json:"project_id"`
	Subject    string        `json:"subject"`
	CategoryId Nullable[int] `json:"category_id"`
}

func Test_Json_of_struct(t *testing.T) {
	tsk := task{
		TaskId:     1,
		ProjectId:  Value(5),
		Subject:    "kjell",
		CategoryId: Nullable[int]{},
	}
	jsonData, err := json.Marshal(tsk)
	assert.NoError(t, err)
	assertJSONEquals(t, jsonData, `{"task_id":1,"project_id":5,"subject":"kjell","category_id":null}`, "struct marshal")
}

func Test_Struct_from_JSON(t *testing.T) {
	jsonData := []byte(`{"task_id":1,"project_id":5,"subject":"kjell","category_id":null}`)
	var tsk task
	err := json.Unmarshal(jsonData, &tsk)
	assert.NoError(t, err)
	assert.Equal(t, 1, tsk.TaskId)
	assert.True(t, tsk.ProjectId.Valid)
	assert.Equal(t, 5, tsk.ProjectId.Data)
	assert.Equal(t, "kjell", tsk.Subject)
	assert.False(t, tsk.CategoryId.Valid)
}
