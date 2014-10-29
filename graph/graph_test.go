package graph

import (
	"github.com/stretchr/testify/assert"
	"image"
	_ "image/png"
	"os"
	"testing"
)

/*
 * Helper functions
 */

func loadGraphFromImage(imgName string, graphType GraphType) *Graph {
	f, _ := os.Open(imgName)
	defer f.Close()
	img, _, _ := image.Decode(f)
	return FromImage(img, func(p, q Pixel) float64 {
		return 1.0
	}, graphType)
}

/*
 * Tests
 */

func TestInitializationGridGraph(t *testing.T) {
	graph := New(5, 6, GRIDGRAPH)
	assert.Equal(t, 5, graph.Width())
	assert.Equal(t, 6, graph.Height())
	assert.Equal(t, 30, graph.TotalVertices())
	assert.Equal(t, 49, graph.TotalEdges())
}

func TestEdgesReturnAllEdgesGridGraph(t *testing.T) {
	graph := New(5, 6, GRIDGRAPH)
	assert.Equal(t, 49, len(graph.Edges()))
}

func TestGridGraphFromImageInitialization(t *testing.T) {
	graph := loadGraphFromImage("../test/test.png", GRIDGRAPH)
	assert.Equal(t, 100, graph.Width())
	assert.Equal(t, 100, graph.Height())
	assert.Equal(t, 19800, graph.TotalEdges())
	assert.Equal(t, 10000, graph.TotalVertices())
	assert.Equal(t, len(graph.Edges()), graph.TotalEdges())
}

func TestInitializationKingsGraph(t *testing.T) {
	graph := New(5, 6, KINGSGRAPH)
	assert.Equal(t, 30, graph.TotalVertices())
	assert.Equal(t, 89, graph.TotalEdges())
}

func TestEdgesReturnAllEdgesKingsGraph(t *testing.T) {
	graph := New(5, 6, KINGSGRAPH)
	assert.Equal(t, graph.TotalEdges(), len(graph.Edges()))
}

func TestKingsGraphsFromImageInitialization(t *testing.T) {
	graph := loadGraphFromImage("../test/test.png", KINGSGRAPH)
	assert.Equal(t, 100, graph.Width())
	assert.Equal(t, 100, graph.Height())
	assert.Equal(t, 39402, graph.TotalEdges())
	assert.Equal(t, 10000, graph.TotalVertices())
	assert.Equal(t, len(graph.Edges()), graph.TotalEdges())
}
