package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBFS(t *testing.T) {
	g := Graph{
		"A": {"B", "C"},
		"B": {"D"},
		"C": {"D"},
		"D": {},
	}
	path, ok := g.BFS("A", "D")
	assert.True(t, ok)
	assert.Equal(t, []string{"A", "B", "D"}, path)
}

func TestDFS(t *testing.T) {
	g := Graph{
		"A": {"B", "C"},
		"B": {"D"},
		"C": {"D"},
		"D": {},
	}
	path, ok := g.DFS("A", "D")
	assert.True(t, ok)
	assert.Equal(t, "A", path[0])
	assert.Equal(t, "D", path[len(path)-1])
}
