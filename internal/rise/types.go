// Package rise 实现 Briggs 热抬升核算。
//
// 输入烟囱几何高度 hs、烟气出口速度 ws、出口内径 rs、烟气温度 Ts、
// 环境温度 Ta，先算浮力通量 Fb = g·ws·rs²·(Ts−Ta)/Ts，再按风速分段：
//
//   - 无风支（u < 1 m/s）：ΔH = 5.3·Fb^(1/4)，与 u 无关；
//   - 有风支（u ≥ 1 m/s）：
//     不稳定/中性（A–D）：Fb < 55 时 ΔH = 21.425·Fb^(3/4)/u，
//                        Fb ≥ 55 时 ΔH = 38.71·Fb^(3/5)/u；
//     稳定（E/F）：ΔH = 2.9·(Fb/(u·s))^(1/3)，s = g·Γ/Ta。
//
// 有效源高 H = hs + ΔH；给出逆温层底高 z_inv 时，H 被截断到 z_inv。
package rise

// 常数约定（钉死）。
const (
	// CalmWindThreshold 是区分无风/有风两支的风速阈值（m/s）。
	CalmWindThreshold = 1.0
	// FbThreshold 是有风支 Fb<55 与 Fb≥55 两段的分界。
	FbThreshold = 55.0
	// Gravity 是标准重力加速度（m/s²）。
	Gravity = 9.80665
	// DefaultGradient 是稳定支默认位温梯度 Γ（K/m）。
	DefaultGradient = 0.02
)

// StackConfig 描述烟囱与烟气参数（全部以 SI 单位）。
type StackConfig struct {
	Height             float64 // 烟囱几何高度 hs（m）
	ExitVelocity       float64 // 烟气出口速度 ws（m/s）
	Radius             float64 // 出口内径 rs（m）
	GasTemperature     float64 // 烟气温度 Ts（K）
	AmbientTemperature float64 // 环境温度 Ta（K）
	InversionTop       float64 // 逆温层底高 z_inv（m）；≤0 表示不限制
	Gradient           float64 // 位温梯度 Γ（K/m）；0 使用默认值
}

// RiseResult 是一次抬升核算的结果。
type RiseResult struct {
	Flux            float64 // 浮力通量 Fb（m⁴/s³）
	Rise            float64 // 抬升量 ΔH（m）
	EffectiveHeight float64 // 有效源高 hs+ΔH（m，可能被逆温顶截断）
	Capped          bool    // 是否被逆温顶截断
	Calm            bool    // 是否走了无风支
	Stable          bool    // 是否走了稳定支
}

// Unstable returns true when the rise used the windy unstable branch.
func (r RiseResult) Unstable() bool { return !r.Calm && !r.Stable && r.Flux > 0 }
