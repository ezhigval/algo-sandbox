package topk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopK(t *testing.T) {
	items := []Item{
		{"a", 1}, {"b", 5}, {"c", 3}, {"d", 4},
	}
	got := TopK(items, 2)
	assert.Len(t, got, 2)
	assert.Equal(t, 5, got[0].Score)
	assert.Equal(t, 4, got[1].Score)
}
