package rise

import (
	"gauss-plume/internal/dispersion"
)

func Compute(st dispersion.Stability, u float64, cfg StackConfig) (RiseResult, error) {
	if err := ValidateConfig(cfg); err != nil {
		return RiseResult{}, err
	}
	if err := dispersion.ValidateWind(u); err != nil {
		return RiseResult{}, err
	}
	fb, err := BuoyancyFlux(cfg)
	if err != nil {
		return RiseResult{}, err
	}
	res := RiseResult{Flux: fb}

	if !hasBuoyancy(fb) {
		res.Rise = 0
		res.EffectiveHeight = cfg.Height
		eff, capped := capHeight(cfg.Height, res.Rise, cfg.InversionTop)
		res.EffectiveHeight = eff
		res.Capped = capped
		return res, nil
	}

	if u < CalmWindThreshold {
		res.Calm = true
		res.Rise = calmRise(fb)
	} else {
		grad := cfg.Gradient
		if grad <= 0 {
			grad = DefaultGradient
		}
		res.Stable = st.IsStable()
		res.Rise = windyRise(st, u, fb, cfg.AmbientTemperature, grad)
	}

	eff, capped := capHeight(cfg.Height, res.Rise, cfg.InversionTop)
	res.EffectiveHeight = eff
	res.Capped = capped
	return res, nil
}
