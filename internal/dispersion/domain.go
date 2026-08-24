package dispersion

import (
	"fmt"
	"math"
)

// 扩散计算域的边界（米）。幂律参数在超远距离失去意义，
// 因此对下风向距离设上限，超出即报错。
const (
	// MaxDistance 是允许的最大下风向距离（1000 km）。
	MaxDistance = 1e6
)

// ValidateDistance 校验下风向距离：有限、大于 0 且不超过 MaxDistance。
// 幂律公式在 x→0 时发散，故 x≤0 一律拒绝。
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

// ValidateWind 校验风速：有限且大于 0。
func ValidateWind(u float64) error {
	if math.IsNaN(u) || math.IsInf(u, 0) {
		return fmt.Errorf("风速必须是有限数值，实际 %v", u)
	}
	if u <= 0 {
		return fmt.Errorf("风速必须大于 0，实际 %v", u)
	}
	return nil
}

// ValidateFinite 校验任意标量参数为有限数值。
func ValidateFinite(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%s 必须是有限数值，实际 %v", name, v)
	}
	return nil
}

// ValidateNonNegative 校验参数有限且非负。
func ValidateNonNegative(name string, v float64) error {
	if err := ValidateFinite(name, v); err != nil {
		return err
	}
	if v < 0 {
		return fmt.Errorf("%s 必须大于等于 0，实际 %v", name, v)
	}
	return nil
}

// ValidatePositive 校验参数有限且严格为正。
func ValidatePositive(name string, v float64) error {
	if err := ValidateFinite(name, v); err != nil {
		return err
	}
	if v <= 0 {
		return fmt.Errorf("%s 必须大于 0，实际 %v", name, v)
	}
	return nil
}
