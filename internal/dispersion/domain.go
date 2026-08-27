package dispersion

import (
	"fmt"
	"math"
)

const (
	MaxDistance = 1e6
)

func ValidateDistance(x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return fmt.Errorf("下风向距离必须是有限数值，实际 %v", x)
	}
	if x <= 0 {
		return fmt.Errorf("下风向距离必须大于 0，实际 %v", x)
	}
	if x > MaxDistance {
		return fmt.Errorf("下风向距离过大（上限 %g 米），实际 %v", MaxDistance, x)
	}
	return nil
}

func ValidateWind(u float64) error {
	if math.IsNaN(u) || math.IsInf(u, 0) {
		return fmt.Errorf("风速必须是有限数值，实际 %v", u)
	}
	if u <= 0 {
		return fmt.Errorf("风速必须大于 0，实际 %v", u)
	}
	return nil
}

func ValidateFinite(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%s 必须是有限数值，实际 %v", name, v)
	}
	return nil
}

func ValidateNonNegative(name string, v float64) error {
	if err := ValidateFinite(name, v); err != nil {
		return err
	}
	if v < 0 {
		return fmt.Errorf("%s 必须大于等于 0，实际 %v", name, v)
	}
	return nil
}

func ValidatePositive(name string, v float64) error {
	if err := ValidateFinite(name, v); err != nil {
		return err
	}
	if v <= 0 {
		return fmt.Errorf("%s 必须大于 0，实际 %v", name, v)
	}
	return nil
}
