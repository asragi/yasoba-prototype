package frontend

import (
	"image/color"
	"math"
)

func ColorRainbow(t float64) color.RGBA {
	return color.RGBA{
		R: uint8(math.Sin(t)*127 + 128),
		G: uint8(math.Sin(t+2*math.Pi/3)*127 + 128),
		B: uint8(math.Sin(t+4*math.Pi/3)*127 + 128),
		A: 255,
	}
}
