package mathkit

import (
	"math"
)

func round(x float64) int {
	return int(math.Floor(x + 0.5))
}
