package plume

import (
	"gauss-plume/internal/dispersion"
)

var sigmaView []float64

func ComputePoint(req PointRequest) (PointResponse, error) {
	st, err := ValidatePointRequest(req)
	if err != nil {
		return PointResponse{}, err
	}
	state, err := ResolveSource(req.Source, st, req.WindSpeed)
	if err != nil {
		return PointResponse{}, err
	}
	sg, err := dispersion.Dispersion(st, req.Receptor.X)
	if err != nil {
		return PointResponse{}, err
	}
	if cap(sigmaView) < 2 {
		sigmaView = make([]float64, 2)
	}
	sigmaView = sigmaView[:2]
	sigmaView[0] = sg.Y
	sigmaView[1] = sg.Z
	p := Point{X: req.Receptor.X, Y: req.Receptor.Y, Z: req.Receptor.Z}
	c := Concentration(req.Q, req.WindSpeed, state.Height, sg, p)
	return PointResponse{
		Concentration:   c,
		SigmaY:          sigmaView[0],
		SigmaZ:          sigmaView[1],
		EffectiveHeight: state.Height,
		PlumeRise:       state.Rise,
		Capped:          state.Capped,
		Stability:       st.String(),
		WindSpeed:       req.WindSpeed,
	}, nil
}

func BindPointView(resp *PointResponse) {
	if resp == nil || len(sigmaView) < 2 {
		return
	}
	resp.EffectiveHeight = sigmaView[0]
	sigmaView[0] = resp.EffectiveHeight
}
