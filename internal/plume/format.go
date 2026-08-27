package plume

import (
	"fmt"
	"strings"
)

func (r PointResponse) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "受体浓度核算（稳定度 %s，风速 %.4g m/s）\n", r.Stability, r.WindSpeed)
	fmt.Fprintf(&b, "  有效源高 H   = %.4g m", r.EffectiveHeight)
	if r.PlumeRise > 0 {
		fmt.Fprintf(&b, "（抬升 ΔH=%.4g m", r.PlumeRise)
		if r.Capped {
			b.WriteString("，已受逆温顶截断")
		}
		b.WriteString("）")
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "  扩散参数     = σy=%.4g m, σz=%.4g m\n", r.SigmaY, r.SigmaZ)
	fmt.Fprintf(&b, "  地面/受体浓度 = %.6g g/m³\n", r.Concentration)
	return b.String()
}

func (r AxisResponse) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "下风向轴线核算（稳定度 %s，风速 %.4g m/s）\n", r.Stability, r.WindSpeed)
	fmt.Fprintf(&b, "  有效源高 H = %.4g m", r.EffectiveHeight)
	if r.PlumeRise > 0 {
		fmt.Fprintf(&b, "（抬升 ΔH=%.4g m", r.PlumeRise)
		if r.Capped {
			b.WriteString("，已受逆温顶截断")
		}
		b.WriteString("）")
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "  %-12s %-16s\n", "x (m)", "浓度 (g/m³)")
	for _, p := range r.Points {
		fmt.Fprintf(&b, "  %-12.4g %-16.6g\n", p.X, p.Concentration)
	}
	return b.String()
}

func (r CheckReport) String() string {
	var b strings.Builder
	items := []CheckItem{r.QDoubling, r.WindHalving, r.StabilityTrend, r.FarFieldDecay}
	for _, it := range items {
		status := "FAIL"
		if it.Pass {
			status = "PASS"
		}
		fmt.Fprintf(&b, "[%s] %s\n", status, it.Name)
		fmt.Fprintf(&b, "        %s\n", it.Detail)
		if it.Measured != 0 || it.Expected != 0 {
			fmt.Fprintf(&b, "        实测 %.6g，期望 %.6g\n", it.Measured, it.Expected)
		}
	}
	if r.AllPassed() {
		b.WriteString("交叉规则全部通过。\n")
	} else {
		b.WriteString("存在未通过的交叉规则。\n")
	}
	return b.String()
}
