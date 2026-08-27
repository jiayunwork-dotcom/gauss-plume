package plume

import (
	"encoding/json"

	"gauss-plume/internal/rise"
)

type Source struct {
	Height *float64     `json:"height,omitempty"`
	Stack  *StackSource `json:"stack,omitempty"`
}

type StackSource struct {
	Height             float64 `json:"height"`
	ExitVelocity       float64 `json:"exit_velocity"`
	Radius             float64 `json:"radius"`
	GasTemperature     float64 `json:"gas_temperature"`
	AmbientTemperature float64 `json:"ambient_temperature"`
	InversionTop       float64 `json:"inversion_top,omitempty"`
	Gradient           float64 `json:"gradient,omitempty"`
}

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

type Receptor struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Point struct {
	X, Y, Z float64
}

type PointRequest struct {
	Q         float64  `json:"q"`
	Source    Source   `json:"source"`
	WindSpeed float64  `json:"wind_speed"`
	Stability string   `json:"stability"`
	Receptor  Receptor `json:"receptor"`
}

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

type AxisGrid struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Count int     `json:"count"`
	Y     float64 `json:"y,omitempty"`
	Z     float64 `json:"z,omitempty"`
}

type AxisRequest struct {
	Q         float64  `json:"q"`
	Source    Source   `json:"source"`
	WindSpeed float64  `json:"wind_speed"`
	Stability string   `json:"stability"`
	Axis      AxisGrid `json:"axis"`
}

type AxisPoint struct {
	X             float64 `json:"x"`
	Concentration float64 `json:"concentration"`
}

type AxisResponse struct {
	EffectiveHeight float64     `json:"effective_height"`
	PlumeRise       float64     `json:"plume_rise"`
	Capped          bool        `json:"capped"`
	Stability       string      `json:"stability"`
	WindSpeed       float64     `json:"wind_speed"`
	Points          []AxisPoint `json:"points"`
}

type PlumeState struct {
	Height    float64
	Rise      float64
	Capped    bool
	UsedStack bool
}

func AxisRequestFromJSON(data []byte) (AxisRequest, error) {
	var req AxisRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return AxisRequest{}, err
	}
	return req, nil
}

func PointRequestFromJSON(data []byte) (PointRequest, error) {
	var req PointRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return PointRequest{}, err
	}
	return req, nil
}
