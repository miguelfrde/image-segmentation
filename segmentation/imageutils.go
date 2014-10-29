package segmentation

import (
	"github.com/miguelfrde/image-segmentation/disjointset"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"time"
)

/**
 * Color type that implements color.RGBA
 */
type ImageColor struct {
	r, g, b float32
}

/**
 * Returns the RGBA values of the given color
 */
func (color *ImageColor) RGBA() (uint32, uint32, uint32, uint32) {
	return color.getColor(color.r), color.getColor(color.g), color.getColor(color.b), color.getColor(255)
}

/**
 * Compute the RGB value that Image expects
 */
func (color *ImageColor) getColor(val float32) uint32 {
	return uint32(val)<<8 + 0xFF
}

/**
 * Returns a random color ImageColor where 0 <= r,g,b <= 255
 */
func randomColor() ImageColor {
	rand.Seed(time.Now().UTC().UnixNano())
	color := ImageColor{}
	color.r = float32(rand.Intn(256))
	color.g = float32(rand.Intn(256))
	color.b = float32(rand.Intn(256))
	return color
}

/**
 * Returns the image that is generated from the given disjoint set.
 * If the randomColors parameter is true, a random color will be assigned
 * to each segment of the image. If it's false, then the mean color of
 * the pixels in the original image will be assigned to each segment.
 */
func imageFromDisjointSet(set *disjointset.DisjointSet,
	originalimg image.Image, randomColors bool) image.Image {
	resultimg := image.NewNRGBA(originalimg.Bounds())
	width := originalimg.Bounds().Max.X
	height := originalimg.Bounds().Max.Y
	meanColors := make([]ImageColor, width*height, width*height)
	if randomColors {
		for u := 0; u < width*height; u++ {
			meanColors[u] = randomColor()
		}
	} else {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				c := set.Find(x + y*width)
				r, g, b, _ := originalimg.At(x, y).RGBA()
				meanColors[c].r += float32(r>>8) / float32(set.Size(c))
				meanColors[c].g += float32(g>>8) / float32(set.Size(c))
				meanColors[c].b += float32(b>>8) / float32(set.Size(c))
			}
		}
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			resultimg.Set(x, y, &meanColors[set.Find(x+y*width)])
		}
	}
	return resultimg
}
