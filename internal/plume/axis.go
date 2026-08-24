package plume

import (
	"gauss-plume/internal/dispersion"
)

// ComputeAxis 计算下风向轴线（默认地面 y=0、z=0，可由网格覆盖）的
// 浓度点列，并返回有效源高与抬升信息。
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

// AxisPoints 在网格上等距采样浓度点列。
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
	pts := make([]AxisPoint, 0, g.Count)
	for i, x := range xs {
		c := Concentration(q, u, H, sgs[i], Point{X: x, Y: g.Y, Z: g.Z})
		pts = append(pts, AxisPoint{X: x, Concentration: c})
	}
	return pts, nil
}

// GroundAxisPoints 是 AxisPoints 的地面（z=0）快捷版本。
func GroundAxisPoints(q, u, H float64, st dispersion.Stability, g AxisGrid) ([]AxisPoint, error) {
	g.Z = 0
	return AxisPoints(q, u, H, st, g)
}
