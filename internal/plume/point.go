package plume

import (
	"gauss-plume/internal/dispersion"
)

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
	p := Point{X: req.Receptor.X, Y: req.Receptor.Y, Z: req.Receptor.Z}
	c := Concentration(req.Q, req.WindSpeed, state.Height, sg, p)
	return PointResponse{
		Concentration:   c,
		SigmaY:          sg.Y,
		SigmaZ:          sg.Z,
		EffectiveHeight: state.Height,
		PlumeRise:       state.Rise,
		Capped:          state.Capped,
		Stability:       st.String(),
		WindSpeed:       req.WindSpeed,
	}, nil
}
