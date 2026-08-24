package dispersion

import (
	"fmt"
	"math"
)

// Exponents 保存幂律增长指数（σ~x^b 中的 b）。
type Exponents struct {
	Y float64 // σy 的增长指数
	Z float64 // σz 的增长指数
}

// GrowthExponents 用 x 与 2x 两点反推该等级的局部增长指数：
//
//	b = ln(σ(2x)/σ(x)) / ln(2)
//
// 对幂律实现，结果应恰好等于系数表中的 BY/BZ。
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

// GrowthExponent 由两个距离与两个 σ 值反推局部增长指数：
//
//	p = ln(σ2/σ1) / ln(x2/x1)
//
// 要求 x1 < x2 且两者均大于 0。
func GrowthExponent(x1, s1, x2, s2 float64) (float64, error) {
	if x1 <= 0 || x2 <= 0 || x2 == x1 {
		return 0, fmt.Errorf("反推增长指数需要 0 < x1 < x2，实际 %g, %g", x1, x2)
	}
	if s1 <= 0 || s2 <= 0 {
		return 0, fmt.Errorf("反推增长指数需要 σ 为正，实际 %g, %g", s1, s2)
	}
	return math.Log(s2/s1) / math.Log(x2/x1), nil
}

// SumExponent 返回 σy·σz 整体的增长指数，即 σy 与 σz 指数之和。
func (e Exponents) Sum() float64 {
	return e.Y + e.Z
}

// String 返回指数的可读文本。
func (e Exponents) String() string {
	return fmt.Sprintf("b_y=%.4g, b_z=%.4g", e.Y, e.Z)
}
