package segmentation

import (
	"fmt"
	"github.com/miguelfrde/image-segmentation/disjointset"
	"github.com/miguelfrde/image-segmentation/graph"
	"github.com/miguelfrde/image-segmentation/imagenoise"
	"github.com/miguelfrde/image-segmentation/utils"
	"math"
	"sort"
	"time"
)

/**
 * Performs the image segmentation using the "Heuristic for Minimum Spanning
 * Forests" algorithm. It uses the weightfn to compute the weight of the
 * graph edges. minWeight is the only parameter.
 * For more information on this algorithm refer to either my report which link
 * is on the repo's README or to:
 * http://algo2.iti.kit.edu/wassenberg/wassenberg09parallelSegmentation.pdf
 */
func (s *Segmenter) SegmentHMSF(sigmaSmooth, minWeight float64) {
	start := time.Now()
	sigma := imagenoise.EstimateStdev(s.img)
	s.smoothImage(sigmaSmooth)
	s.buildGraph()

	fmt.Printf("segment... ")
	start = time.Now()
	s.resultset = disjointset.New(s.graph.TotalVertices())

	edges := s.graph.Edges()
	sort.Sort(edges)
	setll := s.hmsfMergeEdgesByWeight(edges, minWeight)
	regionCredit := s.hmsfComputeCredit(setll, sigma)
	s.hmsfMergeRegionsByCredit(edges, regionCredit)
	regionCredit[0] = 1
	fmt.Println(time.Since(start))
	fmt.Println("Components:", s.resultset.Components())
}

/**
 * First part of HMSF algorithm. Given a minimum weight, merge edges until that
 * region exceeds that minimum weight
 */
func (s *Segmenter) hmsfMergeEdgesByWeight(edges graph.EdgeList,
	minWeight float64) *disjointset.DisjointSetLL {
	setll := disjointset.NewDisjointSetLL(s.graph.TotalVertices())
	for _, edge := range edges {
		u := s.resultset.Find(edge.U())
		v := s.resultset.Find(edge.V())
		if u != v && edge.Weight() < minWeight {
			root := s.resultset.Union(u, v)
			if root == u {
				setll.Union(root, v)
			} else {
				setll.Union(root, u)
			}
		}
	}
	return setll
}

/**
 * Compute the credit for each region. The credit is defined as:
 * Credit(R) = contrast(R) * sqrt(4 * pi * |R|)
 * where contrast(R) = minWeightInTheRegionBorder - 2 * sigma
 * where sigma is the previously computed standard deviation of the additive
 * white gaussian noise of the image.
 */
func (s *Segmenter) hmsfComputeCredit(setll *disjointset.DisjointSetLL, sigma float64) []float64 {
	regionCredit := make([]float64, s.graph.TotalVertices(), s.graph.TotalVertices())
	minWeights := s.hmfsMinWeights(setll)
	for i := 0; i < s.graph.TotalVertices(); i++ {
		contrast := minWeights[s.resultset.Find(i)] - 2*sigma
		regionCredit[i] = contrast * math.Sqrt(4*math.Pi*float64(s.resultset.Size(i)))
	}
	return regionCredit
}

/**
 * Compute the minimum weight in the border of each region.
 */
func (s *Segmenter) hmfsMinWeights(setll *disjointset.DisjointSetLL) []float64 {
	minWeights := make([]float64, s.graph.TotalVertices(), s.graph.TotalVertices())
	computed := make([]bool, s.graph.TotalVertices(), s.graph.TotalVertices())
	for v := 0; v < s.graph.TotalVertices(); v++ {
		region := s.resultset.Find(v)
		if !computed[region] {
			minWeights[region] = math.Inf(1)
			for w := range setll.Elements(region) {
				for n := range s.graph.Neighbors(w) {
					if s.resultset.Find(n) != region && s.graph.Weight(w, n) < minWeights[region] {
						minWeights[region] = s.graph.Weight(w, n)
					}
				}
			}
			computed[region] = true
		}
	}
	return minWeights
}

/**
 * Last part of the HMSF algorithm, merge regions if the credit of any of them
 * exceeds the weight of the edge connecting them.
 */
func (s *Segmenter) hmsfMergeRegionsByCredit(edges graph.EdgeList, regionCredit []float64) {
	for _, edge := range edges {
		u := s.resultset.Find(edge.U())
		v := s.resultset.Find(edge.V())
		if u != v {
			credit := utils.MinF(regionCredit[u], regionCredit[v])
			if credit > edge.Weight() {
				s.resultset.Union(u, v)
				survivor := s.resultset.Find(u)
				regionCredit[survivor] = credit - edge.Weight()
			}
		}
	}
}
