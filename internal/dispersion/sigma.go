package dispersion

import (
	"fmt"
	"math"
)

type Sigma struct {
	Y float64
	Z float64
}

func SigmaY(s Stability, x float64) (float64, error) {
	c, err := Coefficients(s)
	if err != nil {
		return 0, err
	}
	if err := ValidateDistance(x); err != nil {
		return 0, err
	}
	return c.AY * math.Pow(x, c.BY), nil
}

func SigmaZ(s Stability, x float64) (float64, error) {
	c, err := Coefficients(s)
	if err != nil {
		return 0, err
	}
	if err := ValidateDistance(x); err != nil {
		return 0, err
	}
	return c.AZ * math.Pow(x, c.BZ), nil
}

func Dispersion(s Stability, x float64) (Sigma, error) {
	c, err := Coefficients(s)
	if err != nil {
		return Sigma{}, err
	}
	if err := ValidateDistance(x); err != nil {
		return Sigma{}, err
	}
	return Sigma{
		Y: c.AY * math.Pow(x, c.BY),
		Z: c.AZ * math.Pow(x, c.BZ),
	}, nil
}

func (sg Sigma) String() string {
	return fmt.Sprintf("σy=%.4g m, σz=%.4g m", sg.Y, sg.Z)
}

func (sg Sigma) Max() float64 {
	if sg.Y > sg.Z {
		return sg.Y
	}
	return sg.Z
}

func (sg Sigma) Product() float64 {
	return sg.Y * sg.Z
}
