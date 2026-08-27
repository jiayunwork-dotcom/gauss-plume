package plume

import (
	"fmt"
	"math"
)

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

type AxisStats struct {
	MinConcentration float64
	MaxConcentration float64
	PeakIndex        int
	PeakDistance     float64
}

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
