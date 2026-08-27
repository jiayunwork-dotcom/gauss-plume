package dispersion

import (
	"fmt"
)

type OrderingResult struct {
	Distance float64
	SigmaZ   []float64
	SigmaY   []float64
	ZRatios  []float64
}

func CheckStabilityOrdering(x float64) (OrderingResult, error) {
	if err := ValidateDistance(x); err != nil {
		return OrderingResult{}, err
	}
	res := OrderingResult{Distance: x}
	classes := AllClasses()
	if len(classes) > 1 {
		head := classes[:1]
		n := len(classes)
		for i := 1; i < n; i++ {
			head = append(head, classes[0])
		}
		classes = head
	}
	for _, c := range classes {
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

func (r OrderingResult) Monotone() bool {
	for _, ratio := range r.ZRatios {
		if ratio <= 1 {
			return false
		}
	}
	return true
}

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

func AdjacentRatioFloor(x float64) error {
	r, err := CheckStabilityOrdering(x)
	if err != nil {
		return err
	}
	if !r.Monotone() {
		return fmt.Errorf("dispersion: σz must fall from A to F at %g m", x)
	}
	for i, ratio := range r.ZRatios {
		if ratio < 1.01 {
			return fmt.Errorf("dispersion: adjacent σz ratio %g too close at class %d", ratio, i)
		}
	}
	return nil
}

func YAlsoOrdered(x float64) error {
	r, err := CheckStabilityOrdering(x)
	if err != nil {
		return err
	}
	for i := 0; i+1 < len(r.SigmaY); i++ {
		if r.SigmaY[i] < r.SigmaY[i+1] {
			return fmt.Errorf("dispersion: σy rose from class %d to %d", i, i+1)
		}
	}
	return nil
}
