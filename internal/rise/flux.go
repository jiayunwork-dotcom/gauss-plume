package rise

// BuoyancyFlux 计算浮力通量：
//
//	Fb = g·ws·rs²·(Ts−Ta)/Ts
//
// 参数非法时返回 error。Ts ≤ Ta 时 Fb ≤ 0，表示无浮力。
func BuoyancyFlux(cfg StackConfig) (float64, error) {
	if err := ValidateConfig(cfg); err != nil {
		return 0, err
	}
	dt := cfg.GasTemperature - cfg.AmbientTemperature
	return Gravity * cfg.ExitVelocity * cfg.Radius * cfg.Radius * dt / cfg.GasTemperature, nil
}

// hasBuoyancy 报告烟气是否存在浮力（Fb 严格为正）。
func hasBuoyancy(fb float64) bool {
	return fb > 0
}
