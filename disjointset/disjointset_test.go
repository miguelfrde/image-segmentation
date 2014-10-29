package disjointset

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func initSet() *DisjointSet {
	set := New(10)
	set.Union(3, 4)
	set.Union(1, 3)
	set.Union(2, 5)
	set.Union(7, 9)
	set.Union(0, 5)
	set.Union(2, 9) // 2-5-9-7-0 3-4-1 6 8
	return set
}

func TestInitializationAllDisconnected(t *testing.T) {
	set := New(5)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			assert.Equal(t, set.Connected(i, j), i == j)
		}
	}
}

func TestUnion(t *testing.T) {
	set := New(5)
	set.Union(0, 4)
	assert.True(t, set.Connected(0, 4))
}

func TestNumberOfComponents(t *testing.T) {
	set := initSet()
	assert.Equal(t, 4, set.Components())
}

func TestComponentSize(t *testing.T) {
	set := initSet()
	assert.Equal(t, 5, set.Size(2))
	assert.Equal(t, 5, set.Size(5))
	assert.Equal(t, 5, set.Size(9))
	assert.Equal(t, 5, set.Size(7))
	assert.Equal(t, 5, set.Size(0))
	assert.Equal(t, 3, set.Size(3))
	assert.Equal(t, 3, set.Size(4))
	assert.Equal(t, 3, set.Size(1))
	assert.Equal(t, 1, set.Size(6))
	assert.Equal(t, 1, set.Size(8))
}
