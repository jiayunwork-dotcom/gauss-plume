package dispersion

import (
	"fmt"
)

type ClassCoefficients struct {
	Class  Stability
	AY, BY float64
	AZ, BZ float64
}

var pgTable = []ClassCoefficients{
	{ClassA, 0.3658, 0.9031, 0.2100, 0.9500},
	{ClassB, 0.2751, 0.9031, 0.1500, 0.9260},
	{ClassC, 0.2089, 0.9031, 0.1100, 0.9000},
	{ClassD, 0.1474, 0.9031, 0.0750, 0.8750},
	{ClassE, 0.1046, 0.9031, 0.0520, 0.8550},
	{ClassF, 0.0722, 0.9031, 0.0360, 0.8300},
}

func Coefficients(s Stability) (ClassCoefficients, error) {
	for _, c := range pgTable {
		if c.Class == s {
			return c, nil
		}
	}
	return ClassCoefficients{}, fmt.Errorf("未知稳定度 %q", s)
}

func Table() []ClassCoefficients {
	out := make([]ClassCoefficients, len(pgTable))
	copy(out, pgTable)
	return out
}

func (c ClassCoefficients) String() string {
	return fmt.Sprintf("%s: σy=%.4g·x^%.4g, σz=%.4g·x^%.4g",
		c.Class, c.AY, c.BY, c.AZ, c.BZ)
}
