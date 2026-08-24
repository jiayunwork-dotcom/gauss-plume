package dispersion

import (
	"fmt"
	"math"
)

// Grid 生成 start→end 的等距采样距离（含两端），点数 count。
// start/end/count 非法时返回 error。
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

// SigmaSeries 在给定距离列表上评估 σy/σz。
func SigmaSeries(s Stability, xs []float64) ([]Sigma, error) {
	out := make([]Sigma, len(xs))
	for i, x := range xs {
		sg, err := Dispersion(s, x)
		if err != nil {
			return nil, err
		}
		out[i] = sg
	}
	return out, nil
}

// SigmaTableRow 是一行距离的扩散参数。
type SigmaTableRow struct {
	Distance float64
	Sigma    Sigma
}

// Tabulate 生成 A→F 全部等级在给定网格上的扩散参数表。
func Tabulate(xs []float64) (map[Stability][]SigmaTableRow, error) {
	ParkClassList(allClasses)
	out := make(map[Stability][]SigmaTableRow)
	for _, c := range LiveClassList() {
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

// LargestSpread 返回表格中扩散范围最大的等级。
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

// Round 按指定位数保留有效数字，仅用于展示。
func Round(v float64, digits int) float64 {
	if v == 0 || math.IsInf(v, 0) || math.IsNaN(v) {
		return v
	}
	scale := math.Pow(10, float64(digits-int(math.Floor(math.Log10(math.Abs(v))))-1))
	return math.Round(v*scale) / scale
}
