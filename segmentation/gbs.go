package segmentation

import (
	"fmt"
	"github.com/miguelfrde/image-segmentation/disjointset"
	"github.com/miguelfrde/image-segmentation/graph"
	"sort"
	"time"
)

/**
 * Performs the image segmentation using the "Graph Based Segmentation"
 * algorithm. It uses sigma to apply a gaussian filter with it to the image
 * to smooth it before running the algorithm.
 * k and minSize are the algorithm parameters. For more information on this
 * algorithm refer to either my report which link is on the repo's README or
 * to: http://cs.brown.edu/~pff/papers/seg-ijcv.pdf
 */
func (s *Segmenter) SegmentGBS(sigma, k float64, minSize int) {
	s.smoothImage(sigma)
	s.buildGraph()
	fmt.Printf("segment... ")
	start := time.Now()
	s.resultset = disjointset.New(s.graph.TotalVertices())
	threshold_vals := make([]float64, s.graph.TotalVertices(), s.graph.TotalVertices())

	for v := 0; v < s.graph.TotalVertices(); v++ {
		threshold_vals[v] = k
	}

	edges := s.graph.Edges()
	sort.Sort(edges)

	s.gbsMergeFromThreshold(edges, threshold_vals, k)
	s.gbsMergeSmallRegions(edges, minSize)

	fmt.Println(time.Since(start))
	fmt.Println("Components:", s.resultset.Components())
}

/**
 * Computes the threshold used by the GBS algorithm.
 * T(c) = k/|c|
 */
func threshold(set *disjointset.DisjointSet, k float64, u int) float64 {
	return k / float64(set.Size(u))
}

/**
 * Performs the union of the regions to which the endpoints of an edge belong to if that
 * edge's weight is less than the thresholds of both regions.
 */
func (s *Segmenter) gbsMergeFromThreshold(edges graph.EdgeList, thresholds []float64, k float64) {
	for _, edge := range edges {
		u := s.resultset.Find(edge.U())
		v := s.resultset.Find(edge.V())
		uok := edge.Weight() <= thresholds[u]
		vok := edge.Weight() <= thresholds[v]
		if !s.resultset.Connected(u, v) && uok && vok {
			s.resultset.Union(u, v)
			new_threshold := edge.Weight() + threshold(s.resultset, k, s.resultset.Find(u))
			thresholds[s.resultset.Find(u)] = new_threshold
		}
	}

}

/**
 * Performs the merge of the regions to which the endpoints of an edge belong to if
 * any of these regions is less than the minimum size for all regions.
 */
func (s *Segmenter) gbsMergeSmallRegions(edges graph.EdgeList, minSize int) {
	for _, edge := range edges {
		u := s.resultset.Find(edge.U())
		v := s.resultset.Find(edge.V())
		if u != v && (s.resultset.Size(u) < minSize || s.resultset.Size(v) < minSize) {
			s.resultset.Union(u, v)
		}
	}
}
