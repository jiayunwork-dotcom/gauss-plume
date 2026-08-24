package plume

import (
	"errors"
	"fmt"
	"math"

	"gauss-plume/internal/dispersion"
)

// CheckCase 是交叉规则自检的基例：一组源、风、稳定度与网格配置。
type CheckCase struct {
	Q         float64  `json:"q"`
	Source    Source   `json:"source"`
	WindSpeed float64  `json:"wind_speed"`
	Stability string   `json:"stability"`
	Axis      AxisGrid `json:"axis"`
}

// CheckItem 是单条交叉规则的测量结果。
type CheckItem struct {
	Pass     bool    `json:"pass"`
	Name     string  `json:"name"`
	Detail   string  `json:"detail"`
	Measured float64 `json:"measured,omitempty"`
	Expected float64 `json:"expected,omitempty"`
}

// CheckReport 是全部交叉规则的汇总。
type CheckReport struct {
	QDoubling      CheckItem `json:"q_doubling"`
	WindHalving    CheckItem `json:"wind_halving"`
	StabilityTrend CheckItem `json:"stability_trend"`
	FarFieldDecay  CheckItem `json:"far_field_decay"`
}

// AllPassed 报告所有规则是否通过。
func (r CheckReport) AllPassed() bool {
	return r.QDoubling.Pass && r.WindHalving.Pass &&
		r.StabilityTrend.Pass && r.FarFieldDecay.Pass
}

// RunChecks 对基例执行四条交叉规则自检：
//
//  1. 只把 Q 加倍，轴线每个点浓度加倍；
//  2. 只把 u 加倍，轴线浓度明显下降（地面源约减半）；
//  3. 不稳定→稳定，同一距离 σz 变小且轴线形态改变；
//  4. 远场地面浓度按 σy·σz ~ x^p 衰减，p ∈ [1.6, 2.0]。
func RunChecks(c CheckCase) (CheckReport, error) {
	st, err := dispersion.ParseStability(c.Stability)
	if err != nil {
		return CheckReport{}, err
	}
	if err := ValidateQ(c.Q); err != nil {
		return CheckReport{}, err
	}
	if err := ValidateSource(c.Source); err != nil {
		return CheckReport{}, err
	}
	if err := dispersion.ValidateWind(c.WindSpeed); err != nil {
		return CheckReport{}, err
	}
	if err := ValidateAxisGrid(c.Axis); err != nil {
		return CheckReport{}, err
	}

	base, err := ComputeAxis(AxisRequest{
		Q: c.Q, Source: c.Source, WindSpeed: c.WindSpeed, Stability: c.Stability, Axis: c.Axis,
	})
	if err != nil {
		return CheckReport{}, err
	}

	// Rule 1: Q doubling.
	dq, err := ComputeAxis(AxisRequest{
		Q: 2 * c.Q, Source: c.Source, WindSpeed: c.WindSpeed, Stability: c.Stability, Axis: c.Axis,
	})
	if err != nil {
		return CheckReport{}, err
	}
	qRatio := axisMeanRatio(base.Points, dq.Points)
	qItem := CheckItem{
		Pass:     approx(qRatio, 2, 1e-6),
		Name:     "Q×2 ⇒ 每个点 C×2",
		Detail:   "轴线浓度几何平均比值应等于 2",
		Measured: qRatio,
		Expected: 2,
	}

	// Rule 2: wind doubling.
	dw, err := ComputeAxis(AxisRequest{
		Q: c.Q, Source: c.Source, WindSpeed: 2 * c.WindSpeed, Stability: c.Stability, Axis: c.Axis,
	})
	if err != nil {
		return CheckReport{}, err
	}
	wRatio := axisMeanRatio(base.Points, dw.Points)
	windOK := wRatio < 0.75
	windDetail := "轴线浓度几何平均比值应 < 0.75"
	if base.EffectiveHeight == 0 {
		windOK = approx(wRatio, 0.5, 1e-6)
		windDetail = "地面源轴线浓度约减半（比值应等于 0.5）"
	}
	wItem := CheckItem{
		Pass:     windOK,
		Name:     "u×2 ⇒ 轴线 C 明显下降",
		Detail:   windDetail,
		Measured: wRatio,
		Expected: 0.5,
	}

	// Rule 3: stability shift from unstable to stable.
	stable := dispersion.ClassF
	sf, err := ComputeAxis(AxisRequest{
		Q: c.Q, Source: c.Source, WindSpeed: c.WindSpeed, Stability: stable.String(), Axis: c.Axis,
	})
	if err != nil {
		return CheckReport{}, err
	}
	szUnstable, err := dispersion.SigmaZ(st, c.Axis.Start)
	if err != nil {
		return CheckReport{}, err
	}
	szStable, err := dispersion.SigmaZ(stable, c.Axis.Start)
	if err != nil {
		return CheckReport{}, err
	}
	trendOK := szStable < szUnstable && !axisProfilesEqual(base.Points, sf.Points)
	sItem := CheckItem{
		Pass:     trendOK,
		Name:     "不稳定→稳定 ⇒ σz 变小、形态改变",
		Detail:   fmt.Sprintf("σz(%s)=%g vs σz(%s)=%g", st, szUnstable, stable, szStable),
		Measured: szStable,
		Expected: szUnstable,
	}

	// Rule 4: far-field decay exponent of the ground axis.
	exp, err := FarFieldDecayExponent(base.Points)
	if err != nil {
		return CheckReport{}, err
	}
	dItem := CheckItem{
		Pass:     exp >= 1.6 && exp <= 2.0,
		Name:     "远场衰减指数",
		Detail:   "地面轴线浓度按 σy·σz ~ x^p 衰减，p 应在 [1.6, 2.0]",
		Measured: exp,
		Expected: 1.8,
	}

	return CheckReport{QDoubling: qItem, WindHalving: wItem, StabilityTrend: sItem, FarFieldDecay: dItem}, nil
}

// FarFieldDecayExponent 用轴线上最后两个点反推浓度衰减指数 p
// （C ~ x^(−p)）。
func FarFieldDecayExponent(pts []AxisPoint) (float64, error) {
	if len(pts) < 2 {
		return 0, errors.New("轴线至少需要两个点才能反推衰减指数")
	}
	a, b := pts[len(pts)-2], pts[len(pts)-1]
	if a.Concentration <= 0 || b.Concentration <= 0 {
		return 0, errors.New("远场浓度为 0，无法反推衰减指数")
	}
	return -math.Log(b.Concentration/a.Concentration) / math.Log(b.X/a.X), nil
}

// axisMeanRatio 返回两条轴线逐点比值的几何平均。
func axisMeanRatio(base, scaled []AxisPoint) float64 {
	sum := 0.0
	n := 0
	for i := range base {
		if base[i].Concentration > 0 {
			sum += math.Log(scaled[i].Concentration / base[i].Concentration)
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return math.Exp(sum / float64(n))
}

// axisProfilesEqual 报告两条轴线是否逐点一致。
func axisProfilesEqual(a, b []AxisPoint) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !approx(a[i].Concentration, b[i].Concentration, 1e-9) {
			return false
		}
	}
	return true
}

// approx 报告 a 与 b 的相对偏差是否在容差内。
func approx(a, b, rel float64) bool {
	return math.Abs(a-b) <= rel*math.Max(math.Abs(a), math.Abs(b))
}
