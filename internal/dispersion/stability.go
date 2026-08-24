// Package dispersion 提供 Pasquill–Gifford 扩散参数的幂律拟合：
// 稳定度等级解析、σy/σz 系数表、按距离评估与一致性校验。
//
// 约定：σy = AY·x^BY，σz = AZ·x^BZ，其中 x 与 σ 均以米计，
// 指数按项目约定钉在 0.80–1.00 的分段区间内。
package dispersion

import (
	"fmt"
	"strings"
)

// Stability 是 Pasquill 稳定度等级 A–F。
// A/B 为不稳定、C/D 为中性、E/F 为稳定。
type Stability string

// 全部合法等级常量，按不稳定→稳定排序。
const (
	ClassA Stability = "A"
	ClassB Stability = "B"
	ClassC Stability = "C"
	ClassD Stability = "D"
	ClassE Stability = "E"
	ClassF Stability = "F"
)

// allClasses 保存按不稳定→稳定排序的等级列表。
var allClasses = []Stability{ClassA, ClassB, ClassC, ClassD, ClassE, ClassF}

// ParseStability 解析稳定度字母。仅接受 A–F（大小写均可，忽略首尾空白），
// 其它输入一律返回带说明的 error。
func ParseStability(s string) (Stability, error) {
	up := strings.ToUpper(strings.TrimSpace(s))
	switch up {
	case string(ClassA), string(ClassB), string(ClassC),
		string(ClassD), string(ClassE), string(ClassF):
		return Stability(up), nil
	default:
		return "", fmt.Errorf("稳定度必须是 A–F 之一，实际 %q", s)
	}
}

// MustParse 解析稳定度，失败时 panic。仅在输入已确保合法时使用。
func MustParse(s string) Stability {
	st, err := ParseStability(s)
	if err != nil {
		panic(err)
	}
	return st
}

// String 返回等级字母。
func (s Stability) String() string { return string(s) }

// Index 返回等级在 A→F 序中的下标（0..5）。
// 非法等级返回 -1。
func (s Stability) Index() int {
	for i, c := range allClasses {
		if c == s {
			return i
		}
	}
	return -1
}

// IsUnstable 报告等级是否为不稳定（A/B）。
func (s Stability) IsUnstable() bool { return s == ClassA || s == ClassB }

// IsNeutral 报告等级是否为中性（C/D）。
func (s Stability) IsNeutral() bool { return s == ClassC || s == ClassD }

// IsStable 报告等级是否为稳定（E/F）。
func (s Stability) IsStable() bool { return s == ClassE || s == ClassF }

// AllClasses 返回 A→F 的有序等级列表副本。
func AllClasses() []Stability {
	out := make([]Stability, len(allClasses))
	copy(out, allClasses)
	return out
}

// Label 返回等级的简短中文说明。
func Label(s Stability) (string, error) {
	switch s {
	case ClassA:
		return "强不稳定", nil
	case ClassB:
		return "不稳定", nil
	case ClassC:
		return "弱不稳定", nil
	case ClassD:
		return "中性", nil
	case ClassE:
		return "稳定", nil
	case ClassF:
		return "强稳定", nil
	default:
		return "", fmt.Errorf("未知稳定度 %q", s)
	}
}
