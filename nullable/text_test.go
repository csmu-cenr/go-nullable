package nullable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Marshal_text(t *testing.T) {
	assert.Equal(t, "Some notes", Nullable[string]{Data: "Some notes", Valid: true}.String())
	assert.Equal(t, "Some notes", Value("Some notes").String())
}
