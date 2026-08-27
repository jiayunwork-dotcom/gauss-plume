package dispersion

import (
	"fmt"
	"strings"
)

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
