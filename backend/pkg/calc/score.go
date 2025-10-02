package calc

import "math"

func CalculateAccessibilityScore(scores []int) float64 {
	if len(scores) == 0 {
		return 0
	}
	total := 0
	for _, score := range scores {
		total += score
	}
	average := float64(total) / float64(len(scores))
	return toFixed(average, 1)
}
func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}