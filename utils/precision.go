package utils

import (
	"fmt"
	"math"
)

const FloatPrecision = 4

func SetPrecision(n float64, precision int) float64 {
	tmp := math.Pow(10, float64(precision))

	return math.Round(n*tmp) / tmp
}

func PrecisionString(n float64) string {
	return fmt.Sprintf("%v", SetPrecision(n, FloatPrecision))
}
