package plume

import (
	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/rise"
)

func ResolveSource(src Source, st dispersion.Stability, u float64) (PlumeState, error) {
	if err := ValidateSource(src); err != nil {
		return PlumeState{}, err
	}
	if src.Height != nil {
		return PlumeState{Height: *src.Height}, nil
	}
	res, err := rise.Compute(st, u, src.Stack.toConfig())
	if err != nil {
		return PlumeState{}, err
	}
	return PlumeState{
		Height:    res.EffectiveHeight,
		Rise:      res.Rise,
		Capped:    res.Capped,
		UsedStack: true,
	}, nil
}

func (s PlumeState) EffectiveHeight() float64 { return s.Height }
