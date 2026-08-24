// Package plume 实现高斯烟羽点源扩散核算内核。
//
// 浓度公式（钉死的 2π 形式，含地面镜像反射项）：
//
//	C = Q/(2π·u·σy·σz) · exp(−y²/2σy²) · [exp(−(z−H)²/2σz²) + exp(−(z+H)²/2σz²)]
//
// 其中 H 为有效源高：直接给出 height 时 H=height；给出烟囱参数时
// 叠加 Briggs 热抬升，H=hs+ΔH（可被逆温顶截断）。地面受体 z=0 时
// 两个镜像项相等，公式退化为 π 形式。
//
// 本包还提供交叉规则自检（Q 加倍、风速加倍、稳定度改变、远场衰减），
// 供离线 check 子命令与测试调用。
package plume

import (
	"encoding/json"

	"gauss-plume/internal/rise"
)

// Source 是点源描述：直接给有效源高 height，或给烟囱参数 stack 由 Briggs 抬升。
// 两者只能二选一。
type Source struct {
	Height *float64     `json:"height,omitempty"`
	Stack  *StackSource `json:"stack,omitempty"`
}

// StackSource 是烟囱参数（JSON 契约，与 rise.StackConfig 对应）。
type StackSource struct {
	Height             float64 `json:"height"`
	ExitVelocity       float64 `json:"exit_velocity"`
	Radius             float64 `json:"radius"`
	GasTemperature     float64 `json:"gas_temperature"`
	AmbientTemperature float64 `json:"ambient_temperature"`
	InversionTop       float64 `json:"inversion_top,omitempty"`
	Gradient           float64 `json:"gradient,omitempty"`
}

// toConfig 把 JSON 契约转换为 rise.StackConfig。
func (s StackSource) toConfig() rise.StackConfig {
	return rise.StackConfig{
		Height:             s.Height,
		ExitVelocity:       s.ExitVelocity,
		Radius:             s.Radius,
		GasTemperature:     s.GasTemperature,
		AmbientTemperature: s.AmbientTemperature,
		InversionTop:       s.InversionTop,
		Gradient:           s.Gradient,
	}
}

// Receptor 是受体三维坐标（米）。x 为下风向距离，y 为横向偏移，
// z 为受体高度。
type Receptor struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Point 是受体坐标的域类型。
type Point struct {
	X, Y, Z float64
}

// PointRequest 是 POST /api/conc 的请求体。
type PointRequest struct {
	Q         float64  `json:"q"`
	Source    Source   `json:"source"`
	WindSpeed float64  `json:"wind_speed"`
	Stability string   `json:"stability"`
	Receptor  Receptor `json:"receptor"`
}

// PointResponse 是单点核算的返回体。
type PointResponse struct {
	Concentration   float64 `json:"concentration"`
	SigmaY          float64 `json:"sigma_y"`
	SigmaZ          float64 `json:"sigma_z"`
	EffectiveHeight float64 `json:"effective_height"`
	PlumeRise       float64 `json:"plume_rise"`
	Capped          bool    `json:"capped"`
	Stability       string  `json:"stability"`
	WindSpeed       float64 `json:"wind_speed"`
}

// AxisGrid 是下风向网格。Start/End 为起点与终点（米），Count 为采样点数。
// Y/Z 为固定横向与高度坐标（默认地面 y=0、z=0）。
type AxisGrid struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Count int     `json:"count"`
	Y     float64 `json:"y,omitempty"`
	Z     float64 `json:"z,omitempty"`
}

// AxisRequest 是 POST /api/axis 的请求体。
type AxisRequest struct {
	Q         float64  `json:"q"`
	Source    Source   `json:"source"`
	WindSpeed float64  `json:"wind_speed"`
	Stability string   `json:"stability"`
	Axis      AxisGrid `json:"axis"`
}

// AxisPoint 是轴线上一个采样点。
type AxisPoint struct {
	X             float64 `json:"x"`
	Concentration float64 `json:"concentration"`
}

// AxisResponse 是轴线核算的返回体。
type AxisResponse struct {
	EffectiveHeight float64     `json:"effective_height"`
	PlumeRise       float64     `json:"plume_rise"`
	Capped          bool        `json:"capped"`
	Stability       string      `json:"stability"`
	WindSpeed       float64     `json:"wind_speed"`
	Points          []AxisPoint `json:"points"`
}

// PlumeState 是解析后的源状态。
type PlumeState struct {
	Height    float64 // 有效源高 H（m）
	Rise      float64 // 抬升量 ΔH（m；直接给 height 时为 0）
	Capped    bool    // 是否被逆温顶截断
	UsedStack bool    // 是否走了烟囱抬升
}

// AxisRequestFromJSON 解析轴线条目 JSON，供 CLI 复用。
func AxisRequestFromJSON(data []byte) (AxisRequest, error) {
	var req AxisRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return AxisRequest{}, err
	}
	return req, nil
}

// PointRequestFromJSON 解析单点条目 JSON，供 CLI 复用。
func PointRequestFromJSON(data []byte) (PointRequest, error) {
	var req PointRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return PointRequest{}, err
	}
	return req, nil
}
