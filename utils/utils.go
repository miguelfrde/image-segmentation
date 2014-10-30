/**
 * utils package contains common utility functions that are used in many
 * places in the project.
 */

package utils

import (
	"image/color"
	"math"
)

/**
 * Computes the minimum of two float values
 */
func MinF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

/**
 * Computes the minimum of two float values
 */
func MinI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/**
 * Rounds the number X.Y
 * Returns X if Y < 0.5 and X+1 if Y >= 0.5
 */
func Round(x float64) int {
	dec_part := x - math.Floor(x)
	if dec_part >= 0.5 {
		return int(x + 1)
	}
	return int(x)
}

/**
 * Returns the grayscale intensity of the color clr
 */
func Intensity(clr color.Color) float64 {
	r, g, b, _ := clr.RGBA()
	r, g, b = r>>8, g>>8, b>>8
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}
