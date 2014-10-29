package segmentation

import (
	"github.com/miguelfrde/image-segmentation/graph"
	"math"
)

/**
 * Computes the Euclidean distance between two pixels.
 * For p1 = (r1, g1, b1) and p2 = (r2, g2, b2):
 * d = sqrt((r2 - r1)^2 + (g2 - g1)^2 + (b2 - b1)^2)
 */
func NNWeight(p1 graph.Pixel, p2 graph.Pixel) float64 {
	r1, g1, b1, _ := p1.Color.RGBA()
	r2, g2, b2, _ := p2.Color.RGBA()
	r1, g1, b1 = r1>>8, g1>>8, b1>>8
	r2, g2, b2 = r2>>8, g2>>8, b2>>8
	return math.Sqrt(math.Pow(float64(r2-r1), 2) + math.Pow(float64(g2-g1), 2) +
		math.Pow(float64(b2-b1), 2))
}
