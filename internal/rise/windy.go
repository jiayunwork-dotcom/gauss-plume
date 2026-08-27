package rise

import (
	"math"

	"gauss-plume/internal/dispersion"
)

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

func LowFluxRise(fb, u float64) float64 {
	return 21.425 * math.Pow(fb, 0.75) / u
}

func HighFluxRise(fb, u float64) float64 {
	return 38.71 * math.Pow(fb, 0.60) / u
}

func StableRise(fb, u, ta, grad float64) float64 {
	s := Gravity * grad / ta
	return 2.9 * math.Cbrt(fb/(u*s))
}

func CalmRise(fb float64) float64 {
	return calmRise(fb)
}
