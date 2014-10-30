package disjointset

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func initSetLL() *DisjointSetLL {
	set := NewDisjointSetLL(10)
	set.Union(3, 4)
	set.Union(1, 3)
	set.Union(2, 5)
	set.Union(7, 9)
	set.Union(0, 5)
	set.Union(2, 9) // 2-5-9-7-0 3-4-1 6 8
	return set
}

func TestInitializationAllDisconnectedLL(t *testing.T) {
	set := NewDisjointSetLL(5)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			assert.Equal(t, set.Connected(i, j), i == j)
		}
	}
}

func TestUnionLL(t *testing.T) {
	set := NewDisjointSetLL(5)
	set.Union(0, 4)
	assert.True(t, set.Connected(0, 4))
}

func TestNumberOfComponentsLL(t *testing.T) {
	set := initSetLL()
	assert.Equal(t, 4, set.TotalComponents())
}

func TestComponentSizeLL(t *testing.T) {
	set := initSetLL()
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

func TestElementsINterator(t *testing.T) {
	set := initSetLL()
	var testElements func([]int)
	testElements = func(vertices []int) {
		for _, v := range vertices {
			visited := make([]bool, 10, 10)
			total := 0
			for e := range set.Elements(v) {
				assert.False(t, visited[e])
				visited[e] = true
				total++
			}
			assert.Equal(t, len(vertices), total)
		}
	}
	testElements([]int{2, 5, 9, 7, 0})
	testElements([]int{3, 4, 1})
	testElements([]int{6})
	testElements([]int{8})
}
