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
	return q * crosswind * ground / (math.Pi * u * sg.Y * sg.Z)
}

func MirrorTerms(H, z, sgZ float64) (primary, image float64) {
	dP := z - H
	dI := z + H
	primary = math.Exp(-dP * dP / (2 * sgZ * sgZ))
	image = math.Exp(-dI * dI / (2 * sgZ * sgZ))
	return primary, image
}
