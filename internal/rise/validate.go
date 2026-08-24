package rise

import (
	"fmt"

	"gauss-plume/internal/dispersion"
)

// ValidateConfig 校验烟囱与烟气参数。
// 逆温层底高若给出，必须严格高于烟囱顶（z_inv > hs）。
func ValidateConfig(cfg StackConfig) error {
	if err := dispersion.ValidateNonNegative("烟囱高度 hs", cfg.Height); err != nil {
		return err
	}
	if err := dispersion.ValidatePositive("烟气出口速度 ws", cfg.ExitVelocity); err != nil {
		return err
	}
	if err := dispersion.ValidatePositive("烟囱出口内径 rs", cfg.Radius); err != nil {
		return err
	}
	if err := dispersion.ValidatePositive("烟气温度 Ts", cfg.GasTemperature); err != nil {
		return err
	}
	if err := dispersion.ValidatePositive("环境温度 Ta", cfg.AmbientTemperature); err != nil {
		return err
	}
	if err := dispersion.ValidateNonNegative("位温梯度 Γ", cfg.Gradient); err != nil {
		return err
	}
	if cfg.InversionTop > 0 && cfg.InversionTop <= cfg.Height {
		return fmt.Errorf("逆温层底高必须高于烟囱顶（%g m），实际 %g m", cfg.Height, cfg.InversionTop)
	}
	return nil
}
