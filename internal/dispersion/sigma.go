package dispersion

import (
	"fmt"
	"math"
)

// Sigma 保存某一下风向距离处的横向与垂直扩散参数（米）。
type Sigma struct {
	Y float64 // 横向扩散参数 σy
	Z float64 // 垂直扩散参数 σz
}

// SigmaY 评估 σy = AY·x^BY。
func SigmaY(s Stability, x float64) (float64, error) {
	c, err := Coefficients(s)
	if err != nil {
		return 0, err
	}
	if err := ValidateDistance(x); err != nil {
		return 0, err
	}
	return c.AY * math.Pow(x, c.BY), nil
}

// SigmaZ 评估 σz = AZ·x^BZ。
func SigmaZ(s Stability, x float64) (float64, error) {
	c, err := Coefficients(s)
	if err != nil {
		return 0, err
	}
	if err := ValidateDistance(x); err != nil {
		return 0, err
	}
	return c.AZ * math.Pow(x, c.BZ), nil
}

// Dispersion 同时评估 σy 与 σz，返回结构体；距离或等级非法时返回 error。
func Dispersion(s Stability, x float64) (Sigma, error) {
	c, err := Coefficients(s)
	if err != nil {
		return Sigma{}, err
	}
	if err := ValidateDistance(x); err != nil {
		return Sigma{}, err
	}
	return Sigma{
		Y: c.AY * math.Pow(x, c.BY),
		Z: c.AZ * math.Pow(x, c.BZ),
	}, nil
}

// String 返回扩散参数的可读文本。
func (sg Sigma) String() string {
	return fmt.Sprintf("σy=%.4g m, σz=%.4g m", sg.Y, sg.Z)
}

// Max 返回两个扩散参数中的较大者。
func (sg Sigma) Max() float64 {
	if sg.Y > sg.Z {
		return sg.Y
	}
	return sg.Z
}

// Product 返回 σy·σz（体积展宽度量），远场浓度按此尺度衰减。
func (sg Sigma) Product() float64 {
	return sg.Y * sg.Z
}
