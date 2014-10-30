package segmentation

import (
	"github.com/miguelfrde/image-segmentation/graph"
	"github.com/miguelfrde/image-segmentation/utils"
	"math"
)

/**
 * Computes the Euclidean distance between two pixels.
 * For p1 = (r1, g1, b1) and p2 = (r2, g2, b2):
 * d = sqrt((r2 - r1)^2 + (g2 - g1)^2 + (b2 - b1)^2)
 */
func NNWeight(p1 graph.Pixel, p2 graph.Pixel) float64 {
	ur1, ug1, ub1, _ := p1.Color.RGBA()
	ur2, ug2, ub2, _ := p2.Color.RGBA()
	r1, g1, b1 := float64(ur1>>8), float64(ug1>>8), float64(ub1>>8)
	r2, g2, b2 := float64(ur2>>8), float64(ug2>>8), float64(ub2>>8)
	return math.Sqrt(math.Pow(float64(r2-r1), 2) + math.Pow(float64(g2-g1), 2) +
		math.Pow(float64(b2-b1), 2))
}

/**
 * Computes the absolute difference between the two pixels intensities.
 */
func IntensityDifference(p1 graph.Pixel, p2 graph.Pixel) float64 {
	return math.Abs(utils.Intensity(p2.Color) - utils.Intensity(p1.Color))
}
