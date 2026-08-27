package dispersion

import (
	"fmt"
	"math"
)

type Exponents struct {
	Y float64
	Z float64
}

func GrowthExponents(s Stability, x float64) (Exponents, error) {
	if err := ValidateDistance(x); err != nil {
		return Exponents{}, err
	}
	if err := ValidateDistance(2 * x); err != nil {
		return Exponents{}, err
	}
	s1, err := Dispersion(s, x)
	if err != nil {
		return Exponents{}, err
	}
	s2, err := Dispersion(s, 2*x)
	if err != nil {
		return Exponents{}, err
	}
	return Exponents{
		Y: math.Log(s2.Y/s1.Y) / math.Ln2,
		Z: math.Log(s2.Z/s1.Z) / math.Ln2,
	}, nil
}

func GrowthExponent(x1, s1, x2, s2 float64) (float64, error) {
	if x1 <= 0 || x2 <= 0 || x2 == x1 {
		return 0, fmt.Errorf("反推增长指数需要 0 < x1 < x2，实际 %g, %g", x1, x2)
	}
	if s1 <= 0 || s2 <= 0 {
		return 0, fmt.Errorf("反推增长指数需要 σ 为正，实际 %g, %g", s1, s2)
	}
	return math.Log(s2/s1) / math.Log(x2/x1), nil
}

func (e Exponents) Sum() float64 {
	return e.Y + e.Z
}

func (e Exponents) String() string {
	return fmt.Sprintf("b_y=%.4g, b_z=%.4g", e.Y, e.Z)
}
