package rise

func BuoyancyFlux(cfg StackConfig) (float64, error) {
	if err := ValidateConfig(cfg); err != nil {
		return 0, err
	}
	dt := cfg.GasTemperature - cfg.AmbientTemperature
	return Gravity * cfg.ExitVelocity * cfg.Radius * cfg.Radius * dt / cfg.GasTemperature, nil
}

func hasBuoyancy(fb float64) bool {
	return fb > 0
}
