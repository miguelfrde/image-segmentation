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
 * graph edges. minWeight and initialCredit are the algorithm parameters.
 * For more information on this algorithm refer to either my report which link
 * is on the repo's README or to:
 * http://algo2.iti.kit.edu/wassenberg/wassenberg09parallelSegmentation.pdf
 */
func (s *Segmenter) SegmentHMSF(minWeight float64) {
	start := time.Now()
	sigma := imagenoise.EstimateStdev(s.img)
	s.smoothImage(0.8)
	s.buildGraph()

	fmt.Printf("segment... ")
	start = time.Now()
	s.resultset = disjointset.New(s.graph.TotalVertices())

	edges := s.graph.Edges()
	sort.Sort(edges)
	minWeights := s.hmsfMergeEdgesByWeight(edges, minWeight)
	regionCredit := s.hmsfComputeCredit(edges, minWeights, sigma)
	s.hmsfMergeRegionsByCredit(edges, regionCredit)

	fmt.Println(time.Since(start))
	fmt.Println("Components:", s.resultset.Components())
}

func (s *Segmenter) hmsfMergeEdgesByWeight(edges graph.EdgeList, minWeight float64) []float64 {
	minWeights := make([]float64, s.graph.TotalVertices(), s.graph.TotalVertices())
	for i := 0; i < s.graph.TotalVertices(); i++ {
		minWeights[i] = math.Inf(1)
	}
	for _, edge := range edges {
		u := s.resultset.Find(edge.U())
		v := s.resultset.Find(edge.V())
		root := u
		if u != v && edge.Weight() < minWeight {
			root = s.resultset.Union(u, v)
		}
		m := utils.MinF(minWeights[root], minWeights[u])
		m = utils.MinF(m, minWeights[v])
		minWeights[root] = utils.MinF(edge.Weight(), m)
	}
	return minWeights
}

func (s *Segmenter) hmsfComputeCredit(edges graph.EdgeList, minWeights []float64,
	sigma float64) []float64 {
	regionCredit := make([]float64, s.graph.TotalVertices(), s.graph.TotalVertices())
	for i := 0; i < s.graph.TotalVertices(); i++ {
		contrast := minWeights[s.resultset.Find(i)] - 2*sigma
		regionCredit[i] = contrast * math.Sqrt(4*math.Pi*float64(s.resultset.Size(i)))
	}
	return regionCredit
}

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
