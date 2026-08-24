package dispersion

import (
	"fmt"
	"strings"
)

// FormatTable 把 Tabulate 的结果格式化为对齐文本，供 CLI 的 sigma 子命令打印。
func FormatTable(tab map[Stability][]SigmaTableRow) string {
	var b strings.Builder
	xs := tab[ClassA]
	header := "x (m)"
	for _, c := range allClasses {
		header += fmt.Sprintf("  σy(%s)", c)
	}
	header += " |"
	for _, c := range allClasses {
		header += fmt.Sprintf("  σz(%s)", c)
	}
	b.WriteString(header)
	b.WriteString("\n")
	for i := range xs {
		row := fmt.Sprintf("%-8.4g", xs[i].Distance)
		for _, c := range allClasses {
			row += fmt.Sprintf("  %-9.4g", tab[c][i].Sigma.Y)
		}
		row += " |"
		for _, c := range allClasses {
			row += fmt.Sprintf("  %-9.4g", tab[c][i].Sigma.Z)
		}
		b.WriteString(row)
		b.WriteString("\n")
	}
	return b.String()
}

// DescribeClass 返回单个等级系数的一行说明。
func DescribeClass(s Stability) (string, error) {
	c, err := Coefficients(s)
	if err != nil {
		return "", err
	}
	label, err := Label(s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s（%s）：σy=%.4g·x^%.4g，σz=%.4g·x^%.4g",
		c.Class, label, c.AY, c.BY, c.AZ, c.BZ), nil
}

// DescribeAllClasses 返回全部等级的系数说明。
func DescribeAllClasses() string {
	var b strings.Builder
	for _, c := range allClasses {
		line, err := DescribeClass(c)
		if err != nil {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}
