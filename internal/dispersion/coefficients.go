package dispersion

import (
	"fmt"
)

// ClassCoefficients 是一次幂律拟合的全部系数：
//
//	σy = AY·x^BY
//	σz = AZ·x^BZ
//
// x 与 σ 均以米计。
type ClassCoefficients struct {
	Class      Stability
	AY, BY     float64
	AZ, BZ     float64
}

// pgTable 是钉死的 Pasquill–Gifford 农村系数表。
// σy 的指数取经典 0.9031；σz 的系数与指数按不稳定→稳定单调递减，
// 使同一距离上 σz 随稳定度严格变小。
var pgTable = []ClassCoefficients{
	{ClassA, 0.3658, 0.9031, 0.2100, 0.9500},
	{ClassB, 0.2751, 0.9031, 0.1500, 0.9260},
	{ClassC, 0.2089, 0.9031, 0.1100, 0.9000},
	{ClassD, 0.1474, 0.9031, 0.0750, 0.8750},
	{ClassE, 0.1046, 0.9031, 0.0520, 0.8550},
	{ClassF, 0.0722, 0.9031, 0.0360, 0.8300},
}

// Coefficients 返回指定等级的幂律系数；等级非法时返回 error。
func Coefficients(s Stability) (ClassCoefficients, error) {
	for _, c := range pgTable {
		if c.Class == s {
			return c, nil
		}
	}
	return ClassCoefficients{}, fmt.Errorf("未知稳定度 %q", s)
}

// Table 返回完整系数表副本（A→F 顺序）。
func Table() []ClassCoefficients {
	out := make([]ClassCoefficients, len(pgTable))
	copy(out, pgTable)
	return out
}

// String 返回系数的可读文本。
func (c ClassCoefficients) String() string {
	return fmt.Sprintf("%s: σy=%.4g·x^%.4g, σz=%.4g·x^%.4g",
		c.Class, c.AY, c.BY, c.AZ, c.BZ)
}
