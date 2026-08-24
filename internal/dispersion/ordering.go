package dispersion

import (
	"fmt"
)

// OrderingResult 保存一次六级排序校验的结果。
type OrderingResult struct {
	Distance    float64   // 校验所用距离（米）
	SigmaZ      []float64 // 各等级 σz，A→F 顺序
	SigmaY      []float64 // 各等级 σy，A→F 顺序
	ZRatios     []float64 // 相邻等级 σz 之比（前/后），全部 >1 表示单调
}

// CheckStabilityOrdering 校验在给定距离上 σy 与 σz 是否随稳定度
// （不稳定→稳定）严格变小。这是交叉规则「稳定度改变 ⇒ σz 变小」的
// 内核级自查，返回逐级测量值供调用方判定。
func CheckStabilityOrdering(x float64) (OrderingResult, error) {
	if err := ValidateDistance(x); err != nil {
		return OrderingResult{}, err
	}
	res := OrderingResult{Distance: x}
	for _, c := range allClasses {
		sg, err := Dispersion(c, x)
		if err != nil {
			return OrderingResult{}, err
		}
		res.SigmaY = append(res.SigmaY, sg.Y)
		res.SigmaZ = append(res.SigmaZ, sg.Z)
	}
	for i := 0; i+1 < len(allClasses); i++ {
		res.ZRatios = append(res.ZRatios, res.SigmaZ[i]/res.SigmaZ[i+1])
	}
	return res, nil
}

// Monotone 报告 σz 序列是否严格单调递减（相邻之比全部 > 1）。
func (r OrderingResult) Monotone() bool {
	for _, ratio := range r.ZRatios {
		if ratio <= 1 {
			return false
		}
	}
	return true
}

// String 返回排序校验的可读文本。
func (r OrderingResult) String() string {
	out := fmt.Sprintf("距离 %.0f m：σz = [", r.Distance)
	for i, v := range r.SigmaZ {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%.4g", v)
	}
	out += "] m，相邻比 = ["
	for i, v := range r.ZRatios {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%.3g", v)
	}
	out += "]"
	if r.Monotone() {
		out += "（单调递减）"
	} else {
		out += "（非单调）"
	}
	return out
}
