package plume

import (
	"fmt"

	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/rise"
)

func ValidateQ(q float64) error {
	if err := dispersion.ValidateNonNegative("源强 Q", q); err != nil {
		return err
	}
	return nil
}

func ValidateReceptor(r Receptor) error {
	if err := dispersion.ValidateDistance(r.X); err != nil {
		return err
	}
	if err := dispersion.ValidateFinite("受体横向坐标 y", r.Y); err != nil {
		return err
	}
	if err := dispersion.ValidateNonNegative("受体高度 z", r.Z); err != nil {
		return err
	}
	return nil
}

func ValidateSource(src Source) error {
	if src.Height != nil && src.Stack != nil {
		return fmt.Errorf("height 与 stack 只能二选一，不能同时给出")
	}
	if src.Height != nil {
		return dispersion.ValidateNonNegative("有效源高 height", *src.Height)
	}
	if src.Stack != nil {
		return rise.ValidateConfig(src.Stack.toConfig())
	}
	return fmt.Errorf("缺少有效源高：必须提供 height 或 stack")
}

func ValidateAxisGrid(g AxisGrid) error {
	if err := dispersion.ValidateDistance(g.Start); err != nil {
		return err
	}
	if err := dispersion.ValidateDistance(g.End); err != nil {
		return err
	}
	if g.End <= g.Start {
		return fmt.Errorf("轴线终点 end 必须大于起点 start（%g m），实际 %g m", g.Start, g.End)
	}
	if g.Count < 2 {
		return fmt.Errorf("轴线网格点数 count 必须 ≥ 2，实际 %d", g.Count)
	}
	if g.Count > 10000 {
		return fmt.Errorf("轴线网格点数 count 过大（上限 10000），实际 %d", g.Count)
	}
	if err := dispersion.ValidateFinite("轴线横向坐标 y", g.Y); err != nil {
		return err
	}
	if err := dispersion.ValidateNonNegative("轴线高度 z", g.Z); err != nil {
		return err
	}
	return nil
}

func ValidatePointRequest(req PointRequest) (dispersion.Stability, error) {
	st, err := dispersion.ParseStability(req.Stability)
	if err != nil {
		return "", err
	}
	if err := ValidateQ(req.Q); err != nil {
		return "", err
	}
	if err := dispersion.ValidateWind(req.WindSpeed); err != nil {
		return "", err
	}
	if err := ValidateSource(req.Source); err != nil {
		return "", err
	}
	if err := ValidateReceptor(req.Receptor); err != nil {
		return "", err
	}
	return st, nil
}

func ValidateAxisRequest(req AxisRequest) (dispersion.Stability, error) {
	st, err := dispersion.ParseStability(req.Stability)
	if err != nil {
		return "", err
	}
	if err := ValidateQ(req.Q); err != nil {
		return "", err
	}
	if err := dispersion.ValidateWind(req.WindSpeed); err != nil {
		return "", err
	}
	if err := ValidateSource(req.Source); err != nil {
		return "", err
	}
	if err := ValidateAxisGrid(req.Axis); err != nil {
		return "", err
	}
	return st, nil
}
