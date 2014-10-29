/**
 * Package graph implements a Graph that can be either a King's graph
 * or a Grid graph. It can generate a graph from a given image.
 */
package graph

import (
	"image"
	"image/color"
	"math"
)

/**
 * Used to compute the weight of an edge when generating a graph
 * from an image
 */
type Pixel struct {
	X, Y  int
	Color color.Color
}

/**
 * Type of the functions that are used to compute the weight of
 * an edge when generating a graph from an image
 */
type WeightFn func(Pixel, Pixel) float64

/**
 * Represents a graph edge. It contains the ids of the two vertices
 * that it connects and the weight of the edge between them
 */
type Edge struct {
	u, v   int
	weight float64
}

/**
 * Return the id of one of the vertices that Edge e connects
 */
func (e *Edge) U() int {
	return e.u
}

/**
 * Return the id of one of the vertices that Edge e connects
 */
func (e *Edge) V() int {
	return e.v
}

/**
 * Return the weight of Edge e
 */
func (e *Edge) Weight() float64 {
	return e.weight
}

/**
 * Used to store all the edges that the graph contains.
 * It can be used with sort.Sort
 */
type EdgeList []Edge

/**
 * Return the number of edges that the EdgeList edges stores
 */
func (edges EdgeList) Len() int {
	return len(edges)
}

/**
 * Returns true if the weight of the edge i is less than the weight
 * of the edge j
 */
func (edges EdgeList) Less(i, j int) bool {
	return edges[i].weight < edges[j].weight
}

/**
 * In place swap of two edges i and j in the edge list
 */
func (edges EdgeList) Swap(i, j int) {
	edges[i], edges[j] = edges[j], edges[i]
}

/**
 * Used to recongise which type of graph to generate
 */
type GraphType int

const (
	GRIDGRAPH  GraphType = iota
	KINGSGRAPH GraphType = iota
)

/**
 * Graph datatype. Contains a list of edges, the graph width and height and
 * if it's a King's graph or a Grid graph
 */
type Graph struct {
	edges         EdgeList
	width, height int
	graphType     GraphType
}

/**
 * Returns a new width x height King's or Grid graph. It assigns a weight of Infinity
 * to all edges
 */
func New(width, height int, graphType GraphType) *Graph {
	g := new(Graph)
	g.width = width
	g.height = height
	g.graphType = graphType
	g.edges = make(EdgeList, 0, g.TotalEdges())

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			p := x + y*g.width
			for n := range g.Neighbors(p) {
				g.edges = append(g.edges, Edge{u: p, v: n, weight: math.Inf(1)})
			}
		}
	}
	return g
}

/**
 * Returns a new graph that represents the image img. The graph will be either
 * a King's grph or a Grid graph. It will compute the edge weights using the
 * provided function weight.
 */
func FromImage(img image.Image, weight WeightFn, graphType GraphType) *Graph {
	g := new(Graph)
	g.height = img.Bounds().Max.Y
	g.width = img.Bounds().Max.X
	g.graphType = graphType
	g.edges = make(EdgeList, 0, g.TotalEdges())

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			p := x + y*g.width
			pixel := Pixel{X: x, Y: y, Color: img.At(x, y)}
			for n := range g.Neighbors(p) {
				x2, y2 := n%g.width, n/g.width
				pixel2 := Pixel{X: x2, Y: x2, Color: img.At(x2, y2)}
				w := weight(pixel, pixel2)
				g.edges = append(g.edges, Edge{u: p, v: n, weight: w})
			}
		}
	}
	return g
}

/**
 * Return the ids of the vertices to which v is adjacent
 */
func (g *Graph) Neighbors(v int) <-chan int {
	ch := make(chan int)
	go func() {
		x, y := v%g.width, v/g.width
		if x+1 < g.width {
			ch <- v + 1
		}
		if y+1 < g.height {
			ch <- v + g.width
		}
		if g.graphType == KINGSGRAPH {
			if y-1 >= 0 && x+1 < g.width {
				ch <- v - g.width + 1
			}
			if y+1 < g.height && x+1 < g.width {
				ch <- v + g.width + 1
			}
		}
		close(ch)
	}()
	return ch
}

/**
 * Returns the width of a graph
 */
func (g *Graph) Width() int {
	return g.width
}

/**
 * Returns the height of a graph
 */
func (g *Graph) Height() int {
	return g.height
}

/**
 * Returns the total number of edges that the graph has
 */
func (g *Graph) TotalEdges() int {
	if g.graphType == KINGSGRAPH {
		return 4*g.width*g.height - 3*(g.width+g.height) + 2
	}
	return (g.width-1)*g.height + g.width*(g.height-1)
}

/**
 * Returns the total number of vertices that the graph has
 */
func (g *Graph) TotalVertices() int {
	return g.width * g.height
}

/**
 * Returns all the edges that the graph has
 */
func (g *Graph) Edges() EdgeList {
	return g.edges
}
