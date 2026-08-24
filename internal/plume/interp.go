package plume

import (
	"fmt"
	"math"
)

// AxisConcentrationAt 在已算好的轴线点列上，对指定距离 x 做
// 对数线性插值得到浓度。x 超出网格范围时返回 error。
func AxisConcentrationAt(pts []AxisPoint, x float64) (float64, error) {
	if len(pts) < 2 {
		return 0, fmt.Errorf("轴线至少需要两个点才能插值，实际 %d 个", len(pts))
	}
	if x < pts[0].X || x > pts[len(pts)-1].X {
		return 0, fmt.Errorf("插值距离 %g 超出轴线范围 [%g, %g]", x, pts[0].X, pts[len(pts)-1].X)
	}
	for i := 0; i+1 < len(pts); i++ {
		if x < pts[i].X {
			return 0, fmt.Errorf("插值距离 %g 不在任何网格区间内", x)
		}
		if x == pts[i].X {
			return pts[i].Concentration, nil
		}
		if x > pts[i].X && x <= pts[i+1].X {
			return logInterp(pts[i], pts[i+1], x), nil
		}
	}
	return pts[len(pts)-1].Concentration, nil
}

// logInterp 在两个相邻点之间做对数线性插值。
// 浓度取对数后线性，避免线性插值在近源区产生负值。
func logInterp(a, b AxisPoint, x float64) float64 {
	if x <= a.X {
		return a.Concentration
	}
	if x >= b.X {
		return b.Concentration
	}
	t := (x - a.X) / (b.X - a.X)
	if a.Concentration <= 0 || b.Concentration <= 0 {
		return a.Concentration*(1-t) + b.Concentration*t
	}
	la, lb := math.Log(a.Concentration), math.Log(b.Concentration)
	return math.Exp(la + (lb-la)*t)
}

// AxisStats 统计轴线浓度的最小值、最大值与位置。
type AxisStats struct {
	MinConcentration float64
	MaxConcentration float64
	PeakIndex        int
	PeakDistance     float64
}

// Stats 计算轴线点列的基本统计。
func Stats(pts []AxisPoint) (AxisStats, error) {
	if len(pts) == 0 {
		return AxisStats{}, fmt.Errorf("轴线为空，无法统计")
	}
	st := AxisStats{MinConcentration: pts[0].Concentration, MaxConcentration: pts[0].Concentration}
	for i, p := range pts {
		if p.Concentration < st.MinConcentration {
			st.MinConcentration = p.Concentration
		}
		if p.Concentration > st.MaxConcentration {
			st.MaxConcentration = p.Concentration
			st.PeakIndex = i
			st.PeakDistance = p.X
		}
	}
	return st, nil
}
