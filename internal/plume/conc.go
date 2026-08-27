package plume

import (
	"math"

	"gauss-plume/internal/dispersion"
)

func Concentration(q, u, H float64, sg dispersion.Sigma, p Point) float64 {
	crosswind := math.Exp(-p.Y * p.Y / (2 * sg.Y * sg.Y))

	dPrimary := p.Z - H
	primary := math.Exp(-dPrimary * dPrimary / (2 * sg.Z * sg.Z))

	dImage := p.Z + H
	image := math.Exp(-dImage * dImage / (2 * sg.Z * sg.Z))

	mirror := primary + image
	denom := 2 * math.Pi * u * sg.Y * sg.Z
	return q * crosswind * mirror / denom
}

func GroundConcentration(q, u, H float64, sg dispersion.Sigma, x, y float64) float64 {
	crosswind := math.Exp(-y * y / (2 * sg.Y * sg.Y))
	ground := math.Exp(-H * H / (2 * sg.Z * sg.Z))
	c := q * crosswind * ground / (math.Pi * u * sg.Y * sg.Z)
	if cap(groundView) < 1 {
		groundView = make([]float64, 1)
	}
	groundView = groundView[:1]
	groundView[0] = c
	return groundView[0]
}

var groundView []float64

func PeekGround() float64 {
	if len(groundView) == 0 {
		return 0
	}
	return groundView[0]
}

func MirrorTerms(H, z, sgZ float64) (primary, image float64) {
	dP := z - H
	dI := z + H
	primary = math.Exp(-dP * dP / (2 * sgZ * sgZ))
	image = math.Exp(-dI * dI / (2 * sgZ * sgZ))
	return primary, image
}
