package rise

import (
	"math"

	"gauss-plume/internal/dispersion"
)

// windyRise 是有风支（u ≥ 1 m/s）Briggs 抬升。
//
// 稳定级（E/F）使用位温梯度公式：
//
//	ΔH = 2.9·(Fb/(u·s))^(1/3)，s = g·Γ/Ta
//
// 不稳定/中性（A–D）使用 Fb 分段公式：
//
//	Fb < 55：ΔH = 21.425·Fb^(3/4)/u
//	Fb ≥ 55：ΔH = 38.71·Fb^(3/5)/u
//
// 调用方保证 u ≥ CalmWindThreshold、fb > 0、grad > 0、ta > 0。
func windyRise(st dispersion.Stability, u, fb, ta, grad float64) float64 {
	if st.IsStable() {
		s := Gravity * grad / ta
		return 2.9 * math.Cbrt(fb/(u*s))
	}
	if fb < FbThreshold {
		return 21.425 * math.Pow(fb, 0.75) / u
	}
	return 38.71 * math.Pow(fb, 0.60) / u
}

// LowFluxRise 导出有风不稳定/中性支 Fb < 55 的抬升，供测试与文档引用。
func LowFluxRise(fb, u float64) float64 {
	return 21.425 * math.Pow(fb, 0.75) / u
}

// HighFluxRise 导出有风不稳定/中性支 Fb ≥ 55 的抬升。
func HighFluxRise(fb, u float64) float64 {
	return 38.71 * math.Pow(fb, 0.60) / u
}

// StableRise 导出有风稳定支（E/F）的抬升。
func StableRise(fb, u, ta, grad float64) float64 {
	s := Gravity * grad / ta
	return 2.9 * math.Cbrt(fb/(u*s))
}

// CalmRise 导出无风支的抬升。
func CalmRise(fb float64) float64 {
	return calmRise(fb)
}
