package imagenoise

import (
	"fmt"
	"github.com/miguelfrde/image-segmentation/utils"
	"github.com/miguelfrde/imaging"
	"image"
	"image/color"
	"math"
	"runtime"
	"time"
)

/**
 * Block sizes for block based noise variance estimation
 */
const (
	BLOCK_WIDTH  = 16
	BLOCK_HEIGHT = 3
)

/**
 * Used for parallel compputeStdevs
 */
type indexValuePair struct {
	i     int
	value float64
}

/**
 * Returns the grayscale intensity of the color clr
 */
func Intensity(clr color.Color) float64 {
	r, g, b, _ := clr.RGBA()
	r, g, b = r>>8, g>>8, b>>8
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

/**
 * Estimates the standard deviation of the additive white gaussian noise
 * in the image.
 * Based on: "Block Based Noise Estimation Using Adaptive Gaussian Filtering"
 */
func EstimateStdev(img image.Image) float64 {
	fmt.Printf("estimate noise stdev...")
	start := time.Now()
	blocks := imageToBlocks(img)
	blocks, minstdev := computeHomogeneousBlocksAndMinStdev(blocks)
	filteredBlocks := filterBlocks(blocks, minstdev)
	sigma := stdevOfBlockDiffs(blocks, filteredBlocks)
	fmt.Println(time.Since(start))
	fmt.Println("Noise stdev =", sigma)
	return sigma
}

/**
 * Returns the a list of blocks, that are subimages of size
 * BLOCK_WIDTH x BLOCK_HEIGHT of the image img
 */
func imageToBlocks(img image.Image) []image.Image {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	blocks := make([]image.Image, 0,
		((w-w%BLOCK_WIDTH)/BLOCK_WIDTH)*((h-h%BLOCK_HEIGHT)/BLOCK_HEIGHT))
	pimg := imaging.ToNRGBA(img)
	for y := 0; y < h; y += BLOCK_HEIGHT {
		for x := 0; x < w; x += BLOCK_WIDTH {
			if w-w%BLOCK_WIDTH == x || h-h%BLOCK_HEIGHT == y {
				continue
			}
			maxX, maxY := x+BLOCK_WIDTH, y+BLOCK_HEIGHT
			blocks = append(blocks, pimg.SubImage(
				image.Rect(x, y, maxX, maxY)))
		}
	}
	return blocks
}

/**
 * From all the blocks, returns the ones which standard deviation is
 * almost the minimum standard deviation of all of them
 */
func computeHomogeneousBlocksAndMinStdev(blocks []image.Image) ([]image.Image, float64) {
	minStdev, stdevs := computeStdevs(blocks)
	roundFn := utils.Round
	if int(math.Floor(minStdev)) == utils.Round(minStdev) {
		roundFn = func(x float64) int {
			return int(math.Floor(x))
		}
	}
	homogeneousBlocks := make([]image.Image, 0, len(blocks))
	for i, stdev := range stdevs {
		if roundFn(stdev) == roundFn(minStdev) {
			homogeneousBlocks = append(homogeneousBlocks, blocks[i])
		}
	}
	return homogeneousBlocks, minStdev
}

/**
 * Computes the standard deviations of the intensities of all blocks and
 * the minimum standard deviation
 */
func computeStdevs(blocks []image.Image) (float64, []float64) {
	stdevs := make([]float64, len(blocks), len(blocks))
	minStdev := math.Inf(1)
	ch := make(chan indexValuePair)
	for i, block := range blocks {
		go func(i int, block image.Image) {
			variance, mean := 0.0, 0.0
			minX, minY := block.Bounds().Min.X, block.Bounds().Min.Y
			maxX, maxY := block.Bounds().Max.X, block.Bounds().Max.Y
			total := float64((maxX - minX) * (maxY - minY))
			for i := minX; i < maxX; i++ {
				for j := minY; j < maxY; j++ {
					mean += Intensity(block.At(i, j))
				}
			}
			mean /= total
			for i := minX; i < maxX; i++ {
				for j := minY; j < maxY; j++ {
					variance += math.Pow(Intensity(block.At(i, j))-mean, 2)
				}
			}
			ch <- indexValuePair{i: i, value: math.Sqrt(variance / total)}
		}(i, block)
	}
	for i := 0; i < len(blocks); i++ {
		pair := <-ch
		stdevs[pair.i] = pair.value
		minStdev = utils.MinF(pair.value, minStdev)
	}
	close(ch)
	return minStdev, stdevs
}

/**
 * Filter all blocks using a gaussian filter with sigma = stdev
 */
func filterBlocks(blocks []image.Image, stdev float64) []image.Image {
	filteredBlocks := make([]image.Image, len(blocks), len(blocks))
	for i, block := range blocks {
		filteredBlocks[i] = imaging.Blur(block, stdev, 5)
	}
	return filteredBlocks
}

/**
 * Returns the standard deviation of the differences between the intensities
 * of the filtered block and the original block
 */
func stdevOfBlockDiffs(origBlocks, filtBlocks []image.Image) float64 {
	diffs, mean := diffsAndMeanDiff(origBlocks, filtBlocks)
	n := float64(len(diffs))
	cpus := runtime.NumCPU()
	ch := make(chan float64)
	for i := 0; i < cpus; i++ {
		go func(start int) {
			sum := 0.0
			for i := start; i < len(diffs); i += cpus {
				sum += math.Pow(diffs[i]-mean, 2) / n
			}
			ch <- sum
		}(i)
	}
	variance := 0.0
	for i := 0; i < cpus; i++ {
		variance += <-ch
	}
	return math.Sqrt(variance)
}

/**
 * Returns the differences between the intensities between each original block and
 * the filtered one. Also return the mean of those differences.
 */
func diffsAndMeanDiff(origBlocks, filtBlocks []image.Image) ([]float64, float64) {
	total := float64(len(origBlocks) * BLOCK_WIDTH * BLOCK_HEIGHT)
	diffs := make([]float64, int(total), int(total))
	mean := 0.0
	for b := 0; b < len(origBlocks); b++ {
		minX, minY := origBlocks[b].Bounds().Min.X, origBlocks[b].Bounds().Min.Y
		maxX, maxY := origBlocks[b].Bounds().Max.X, origBlocks[b].Bounds().Max.Y
		for i := minX; i < maxX; i++ {
			for j := minY; j < maxY; j++ {
				intensityFiltered := Intensity(filtBlocks[b].At(i-minX, j-minY))
				intensityOriginal := Intensity(origBlocks[b].At(i, j))
				diffs[b] = math.Abs(intensityFiltered - intensityOriginal)
				mean += diffs[b] / total
			}
		}
	}
	return diffs, mean
}

/*
SERIAL VERSIONS OF SOME PARALLEL FUNCTIONS:

func computeStdevs(blocks []image.Image) (float64, []float64) {
	stdevs := make([]float64, len(blocks), len(blocks))
	minStdev := math.Inf(1)
	for i, block := range blocks {
		variance, mean := 0.0, 0.0
		minX, minY := block.Bounds().Min.X, block.Bounds().Min.Y
		maxX, maxY := block.Bounds().Max.X, block.Bounds().Max.Y
		total := float64((maxX - minX) * (maxY - minY))
		for i := minX; i < maxX; i++ {
			for j := minY; j < maxY; j++ {
				mean += Intensity(block.At(i, j))
			}
		}
		mean /= total
		for i := minX; i < maxX; i++ {
			for j := minY; j < maxY; j++ {
				variance += math.Pow(Intensity(block.At(i, j))-mean, 2)
			}
		}
		stdevs[i] = math.Sqrt(variance / total)
		minStdev = utils.MinF(stdevs[i], minStdev)
	}
	return minStdev, stdevs
}

func stdevOfBlockDiffs(origBlocks, filtBlocks []image.Image) float64 {
	diffs, mean := diffsAndMeanDiff(origBlocks, filtBlocks)
	start := time.Now()
	n := float64(len(diffs))
	variance := 0.0
	for i := 0; i < len(diffs); i++ {
		variance += math.Pow(diffs[i]-mean, 2) / n
	}
	return math.Sqrt(variance)
}
*/
