package rise

import (
	"gauss-plume/internal/dispersion"
)

// Compute 计算 Briggs 热抬升并给出有效源高。
//
// st 为稳定度（决定有风支走稳定还是不稳定/中性公式），u 为风速。
// u ≤ 0 或烟囱参数非法时返回 error。
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
		// 烟气不高于环境温度：如实报告不抬升，不把 hs 冒充已抬升的 H。
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
