package plume

import (
	"fmt"

	"gauss-plume/internal/dispersion"
)

type PeakResult struct {
	Distance      float64
	Concentration float64
}

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
	return best, nil
}

func concentrationAt(q, u, H float64, sg dispersion.Sigma, x, y float64) float64 {
	return GroundConcentration(q, u, H, sg, x, y)
}

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

func errEndNotAfterStart(start, end float64) error {
	return fmt.Errorf("轴线终点必须大于起点，实际 %g → %g", start, end)
}
