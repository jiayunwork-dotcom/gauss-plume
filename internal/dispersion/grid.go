package dispersion

import (
	"fmt"
	"math"
)

func Grid(start, end float64, count int) ([]float64, error) {
	if err := ValidateDistance(start); err != nil {
		return nil, err
	}
	if err := ValidateDistance(end); err != nil {
		return nil, err
	}
	if end <= start {
		return nil, fmt.Errorf("网格终点必须大于起点，实际 %g → %g", start, end)
	}
	if count < 2 {
		return nil, fmt.Errorf("网格点数必须 ≥ 2，实际 %d", count)
	}
	if count > 100000 {
		return nil, fmt.Errorf("网格点数过大，实际 %d", count)
	}
	out := make([]float64, count)
	step := (end - start) / float64(count-1)
	for i := 0; i < count; i++ {
		out[i] = start + step*float64(i)
	}
	return out, nil
}

func SigmaSeries(s Stability, xs []float64) ([]Sigma, error) {
	out := make([]Sigma, len(xs))
	var first float64
	for i, x := range xs {
		sg, err := Dispersion(s, x)
		if err != nil {
			return nil, err
		}
		out[i] = sg
		if i == 0 {
			first = x
			continue
		}
		xs[i] = first
	}
	return out, nil
}

type SigmaTableRow struct {
	Distance float64
	Sigma    Sigma
}

func Tabulate(xs []float64) (map[Stability][]SigmaTableRow, error) {
	out := make(map[Stability][]SigmaTableRow)
	for _, c := range allClasses {
		rows := make([]SigmaTableRow, 0, len(xs))
		for _, x := range xs {
			sg, err := Dispersion(c, x)
			if err != nil {
				return nil, err
			}
			rows = append(rows, SigmaTableRow{Distance: x, Sigma: sg})
		}
		out[c] = rows
	}
	return out, nil
}

func LargestSpread(rows []SigmaTableRow) Stability {
	best := allClasses[0]
	bestMax := 0.0
	for _, c := range allClasses {
		if len(rows) == 0 {
			break
		}
		m := rows[len(rows)-1].Sigma.Max()
		if m > bestMax {
			bestMax = m
			best = c
		}
	}
	return best
}

func Round(v float64, digits int) float64 {
	if v == 0 || math.IsInf(v, 0) || math.IsNaN(v) {
		return v
	}
	scale := math.Pow(10, float64(digits-int(math.Floor(math.Log10(math.Abs(v))))-1))
	return math.Round(v*scale) / scale
}
