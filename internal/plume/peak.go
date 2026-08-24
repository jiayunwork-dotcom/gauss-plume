package plume

import (
	"fmt"

	"gauss-plume/internal/dispersion"
)

// PeakResult 是沿轴线的峰值位置与浓度。
type PeakResult struct {
	Distance      float64 // 峰值出现的下风向距离 x*（m）
	Concentration float64 // 峰值浓度 Cmax（g/m³）
}

// PeakOnAxis 在 [start, end] 上定位地面轴线的峰值。先用粗网格扫描，
// 再在峰值区间内做对数线性细化。浓度随距离通常为单峰或单调，
// 结果对网格点数不敏感。
func PeakOnAxis(q, u, H float64, st dispersion.Stability, start, end float64) (PeakResult, error) {
	if err := dispersion.ValidateDistance(start); err != nil {
		return PeakResult{}, err
	}
	if err := dispersion.ValidateDistance(end); err != nil {
		return PeakResult{}, err
	}
	if end <= start {
		return PeakResult{}, errEndNotAfterStart(start, end)
	}
	coarse, err := dispersion.Grid(start, end, 256)
	if err != nil {
		return PeakResult{}, err
	}
	sgs, err := dispersion.SigmaSeries(st, coarse)
	if err != nil {
		return PeakResult{}, err
	}
	best := PeakResult{Distance: coarse[0], Concentration: concentrationAt(q, u, H, sgs[0], coarse[0], 0)}
	for i := 1; i < len(coarse); i++ {
		c := concentrationAt(q, u, H, sgs[i], coarse[i], 0)
		if c > best.Concentration {
			best = PeakResult{Distance: coarse[i], Concentration: c}
		}
	}

	// Refine around the coarse peak using a local log-linear scan.
	lo := start
	hi := end
	if best.Distance > start {
		lo = best.Distance - (best.Distance-start)*0.2
	}
	if best.Distance < end {
		hi = best.Distance + (end-best.Distance)*0.2
	}
	fine, err := dispersion.Grid(lo, hi, 512)
	if err != nil {
		return best, nil
	}
	fgs, err := dispersion.SigmaSeries(st, fine)
	if err != nil {
		return best, nil
	}
	for i := range fine {
		c := concentrationAt(q, u, H, fgs[i], fine[i], 0)
		if c > best.Concentration {
			best = PeakResult{Distance: fine[i], Concentration: c}
		}
	}
	bindPeakLive("cmax", best.Concentration)
	return best, nil
}

// concentrationAt 是地面 y=0 的浓度求值。
func concentrationAt(q, u, H float64, sg dispersion.Sigma, x, y float64) float64 {
	return GroundConcentration(q, u, H, sg, x, y)
}

// ConcentrationAt 在任意距离上求地面轴线浓度（直接由公式计算）。
func ConcentrationAt(q, u, H float64, st dispersion.Stability, x, y float64) (float64, error) {
	if err := dispersion.ValidateDistance(x); err != nil {
		return 0, err
	}
	sg, err := dispersion.Dispersion(st, x)
	if err != nil {
		return 0, err
	}
	return GroundConcentration(q, u, H, sg, x, y), nil
}

// errEndNotAfterStart 构造网格端点错误。
func errEndNotAfterStart(start, end float64) error {
	return fmt.Errorf("轴线终点必须大于起点，实际 %g → %g", start, end)
}
