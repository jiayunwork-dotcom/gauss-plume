package rise

import "math"

// calmRise 是无风支（u < 1 m/s）Briggs 抬升：
//
//	ΔH = 5.3·Fb^(1/4)
//
// 与风速无关。调用方保证 fb > 0。
func calmRise(fb float64) float64 {
	return 5.3 * math.Pow(fb, 0.25)
}
