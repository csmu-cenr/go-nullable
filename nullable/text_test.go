package nullable

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Marshal_text(t *testing.T) {
	assert.Equal(t, "Some notes", fmt.Sprintf("%s", Nullable[string]{Data: "Some notes", IsValid: true}))
	assert.Equal(t, "Some notes", fmt.Sprintf("%s", Value("Some notes")))
}
