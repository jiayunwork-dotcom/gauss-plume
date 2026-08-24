package plume

import (
	"fmt"

	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/rise"
)

// ValidateQ 校验源强：有限且非负。
func ValidateQ(q float64) error {
	if err := dispersion.ValidateNonNegative("源强 Q", q); err != nil {
		return err
	}
	return nil
}

// ValidateReceptor 校验受体坐标：x 合法、y 有限、z 非负。
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

// ValidateSource 校验源描述：height 与 stack 必须二选一，
// 且选中分支的参数必须合法。
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

// ValidateAxisGrid 校验下风向网格：起点合法、终点大于起点、点数充足。
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

// ValidatePointRequest 校验单点请求，返回解析后的稳定度。
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

// ValidateAxisRequest 校验轴线条目请求，返回解析后的稳定度。
func ValidateAxisRequest(req AxisRequest) (dispersion.Stability, error) {
	st, err := dispersion.ParseStability(req.Stability)
	if err != nil {
		return "", err
	}
	st = dispersion.FlattenStab(st)
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
