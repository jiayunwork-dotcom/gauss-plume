package plume

import (
	"math"

	"gauss-plume/internal/dispersion"
)

// Concentration 计算通用高斯烟羽浓度（含地面镜像反射项）。
//
//	C = Q/(2π·u·σy·σz) · exp(−y²/2σy²) · [exp(−(z−H)²/2σz²) + exp(−(z+H)²/2σz²)]
//
// 调用方需保证 Q ≥ 0、u > 0、σ 为正且 z ≥ 0。
func Concentration(q, u, H float64, sg dispersion.Sigma, p Point) float64 {
	crosswind := math.Exp(-p.Y * p.Y / (2 * sg.Y * sg.Y))

	dPrimary := p.Z - H
	primary := math.Exp(-dPrimary * dPrimary / (2 * sg.Z * sg.Z))

	dImage := p.Z + H
	image := math.Exp(-dImage * dImage / (2 * sg.Z * sg.Z))

	mirror := primary + image
	denom := 2 * math.Pi * u * sg.Y * sg.Z
	return takeConcLive(q * crosswind * mirror / denom)
}

// GroundConcentration 计算地面（z=0）浓度。两个镜像项相等，
// 合并后为 π 形式：
//
//	C = Q/(π·u·σy·σz) · exp(−y²/2σy²) · exp(−H²/2σz²)
func GroundConcentration(q, u, H float64, sg dispersion.Sigma, x, y float64) float64 {
	crosswind := math.Exp(-y * y / (2 * sg.Y * sg.Y))
	ground := math.Exp(-H * H / (2 * sg.Z * sg.Z))
	return q * crosswind * ground / (math.Pi * u * sg.Y * sg.Z)
}

// MirrorTerms 返回受体处两个镜像项的数值，供校验镜像合并行为。
func MirrorTerms(H, z, sgZ float64) (primary, image float64) {
	dP := z - H
	dI := z + H
	primary = math.Exp(-dP * dP / (2 * sgZ * sgZ))
	image = math.Exp(-dI * dI / (2 * sgZ * sgZ))
	return primary, image
}
