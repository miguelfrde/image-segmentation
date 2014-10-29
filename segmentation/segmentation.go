/**
 * This is probably the most important package in the repository. It contains
 * all the segmentation algorithms and the utilities to transform from an
 * image to a graph and from a DisjointSet forest to an image.
 */
package segmentation

import (
	"fmt"
	"github.com/miguelfrde/image-segmentation/disjointset"
	"github.com/miguelfrde/image-segmentation/graph"
	"github.com/miguelfrde/imaging"
	"image"
	"time"
)

/**
 * Type used to run all the segmentation algorithms.
 * It stores the graph, the resultset, the original image, the graph
 * obtained from the image and if it will generate a result image with
 * random colors.
 */
type Segmenter struct {
	randomColors bool
	img          image.Image
	graph        *graph.Graph
	resultset    *disjointset.DisjointSet
	graphType    graph.GraphType
	weightfn     graph.WeightFn
}

/**
 * Returns a new Segmenter, generates a graph of the given graph type from
 * the given image using the given weight function to compute the edge
 * weights.
 */
func New(img image.Image, graphType graph.GraphType,
	weightfn graph.WeightFn) *Segmenter {
	s := new(Segmenter)
	s.randomColors = false
	s.img = img
	s.weightfn = weightfn
	s.graphType = graphType
	return s
}

func (s *Segmenter) smoothImage(sigma float64) {
	fmt.Printf("blur image... ")
	start := time.Now()
	s.img = imaging.Blur(s.img, sigma, 4)
	fmt.Println(time.Since(start))
}

func (s *Segmenter) buildGraph() {
	fmt.Printf("build graph... ")
	start := time.Now()
	s.graph = graph.FromImage(s.img, s.weightfn, s.graphType)
	fmt.Println(time.Since(start))
}

/**
 * Sets the random color attribute to true or false according to val
 */
func (s *Segmenter) SetRandomColors(val bool) {
	s.randomColors = val
}

/**
 * Returns the result image. Returns nil if no segmentation algorithm
 * has been executed before.
 */
func (s *Segmenter) GetResultImage() image.Image {
	if s.resultset == nil {
		return nil
	}
	fmt.Printf("build image... ")
	start := time.Now()
	resultimg := imageFromDisjointSet(s.resultset, s.img, s.randomColors)
	fmt.Println(time.Since(start))
	return resultimg
}
