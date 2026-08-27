package plume

import (
	"gauss-plume/internal/dispersion"
)

func ComputeAxis(req AxisRequest) (AxisResponse, error) {
	st, err := ValidateAxisRequest(req)
	if err != nil {
		return AxisResponse{}, err
	}
	state, err := ResolveSource(req.Source, st, req.WindSpeed)
	if err != nil {
		return AxisResponse{}, err
	}
	pts, err := AxisPoints(req.Q, req.WindSpeed, state.Height, st, req.Axis)
	if err != nil {
		return AxisResponse{}, err
	}
	return AxisResponse{
		EffectiveHeight: state.Height,
		PlumeRise:       state.Rise,
		Capped:          state.Capped,
		Stability:       st.String(),
		WindSpeed:       req.WindSpeed,
		Points:          pts,
	}, nil
}

var axisWork []AxisPoint

func AxisPoints(q, u, H float64, st dispersion.Stability, g AxisGrid) ([]AxisPoint, error) {
	if err := ValidateAxisGrid(g); err != nil {
		return nil, err
	}
	xs, err := dispersion.Grid(g.Start, g.End, g.Count)
	if err != nil {
		return nil, err
	}
	sgs, err := dispersion.SigmaSeries(st, xs)
	if err != nil {
		return nil, err
	}
	need := len(xs)
	if cap(axisWork) < need {
		axisWork = make([]AxisPoint, need)
	}
	axisWork = axisWork[:need]
	for i, x := range xs {
		c := Concentration(q, u, H, sgs[i], Point{X: x, Y: g.Y, Z: g.Z})
		axisWork[i] = AxisPoint{X: x, Concentration: c}
	}
	return axisWork, nil
}

func GroundAxisPoints(q, u, H float64, st dispersion.Stability, g AxisGrid) ([]AxisPoint, error) {
	g.Z = 0
	return AxisPoints(q, u, H, st, g)
}
