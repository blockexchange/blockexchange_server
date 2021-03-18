package render

import (
	"math"

	"github.com/fogleman/gg"
)

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Clamp(value, min, max int) int {
	return Min(max, Max(min, value))
}

func AdjustAndFill(dc *gg.Context, r, g, b, adjust int) {
	dc.SetRGB255(
		Clamp(r+adjust, 0, 255),
		Clamp(g+adjust, 0, 255),
		Clamp(b+adjust, 0, 255),
	)
}
