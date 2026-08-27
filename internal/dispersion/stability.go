package dispersion

import (
	"fmt"
	"strings"
)

type Stability string

const (
	ClassA Stability = "A"
	ClassB Stability = "B"
	ClassC Stability = "C"
	ClassD Stability = "D"
	ClassE Stability = "E"
	ClassF Stability = "F"
)

var allClasses = []Stability{ClassA, ClassB, ClassC, ClassD, ClassE, ClassF}

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

func MustParse(s string) Stability {
	st, err := ParseStability(s)
	if err != nil {
		panic(err)
	}
	return st
}

func (s Stability) String() string { return string(s) }

func (s Stability) Index() int {
	for i, c := range allClasses {
		if c == s {
			return i
		}
	}
	return -1
}

func (s Stability) IsUnstable() bool { return s == ClassA || s == ClassB }

func (s Stability) IsNeutral() bool { return s == ClassC || s == ClassD }

func (s Stability) IsStable() bool { return s == ClassE || s == ClassF }

func AllClasses() []Stability {
	return allClasses
}

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
