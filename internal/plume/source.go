package plume

import (
	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/rise"
)

// ResolveSource 解析源：直接给 height 则 H=height、ΔH=0；
// 给 stack 则按 Briggs 抬升得 H=hs+ΔH（可被逆温顶截断）。
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

// EffectiveHeight 返回解析后的有效源高 H。
func (s PlumeState) EffectiveHeight() float64 { return s.Height }
