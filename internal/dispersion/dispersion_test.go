package dispersion

import (
	"math"
	"testing"
)

func TestParseStabilityValidAndInvalid(t *testing.T) {
	valid := []string{"A", "B", "C", "D", "E", "F", "a", "d", " f "}
	for _, in := range valid {
		if _, err := ParseStability(in); err != nil {
			t.Errorf("ParseStability(%q) returned error %v, want nil", in, err)
		}
	}
	invalid := []string{"", "G", "g", "H", "AA", "1", "stable", "a1"}
	for _, in := range invalid {
		if _, err := ParseStability(in); err == nil {
			t.Errorf("ParseStability(%q) returned nil error, want error", in)
		}
	}
}

func TestSigmaIncreasesWithDistance(t *testing.T) {
	classes := AllClasses()
	for _, c := range classes {
		near, err := Dispersion(c, 500)
		if err != nil {
			t.Fatalf("Dispersion(%s, 500) error: %v", c, err)
		}
		far, err := Dispersion(c, 5000)
		if err != nil {
			t.Fatalf("Dispersion(%s, 5000) error: %v", c, err)
		}
		if far.Y <= near.Y {
			t.Errorf("%s: sigma_y(5000)=%g not greater than sigma_y(500)=%g", c, far.Y, near.Y)
		}
		if far.Z <= near.Z {
			t.Errorf("%s: sigma_z(5000)=%g not greater than sigma_z(500)=%g", c, far.Z, near.Z)
		}
	}
}

func TestSigmaStabilityOrdering(t *testing.T) {
	for _, x := range []float64{100, 500, 2000, 10000} {
		res, err := CheckStabilityOrdering(x)
		if err != nil {
			t.Fatalf("CheckStabilityOrdering(%g) error: %v", x, err)
		}
		if !res.Monotone() {
			t.Errorf("at x=%g sigma_z not strictly decreasing with stability: %v", x, res.SigmaZ)
		}
	}
}

func TestSigmaGrowthExponentInRange(t *testing.T) {
	// The Pasquill-Gifford convention pins sigma ~ x^b with b in [0.80, 1.00].
	// A regression that replaces the power law with a linear sigma (b=1.0)
	// for every class, or with a wrong exponent, must be caught here.
	for _, c := range AllClasses() {
		ex, err := GrowthExponents(c, 500)
		if err != nil {
			t.Fatalf("GrowthExponents(%s, 500) error: %v", c, err)
		}
		if ex.Y < 0.80 || ex.Y > 1.00 {
			t.Errorf("%s: sigma_y growth exponent %g outside [0.80, 1.00]", c, ex.Y)
		}
		if ex.Z < 0.80 || ex.Z > 1.00 {
			t.Errorf("%s: sigma_z growth exponent %g outside [0.80, 1.00]", c, ex.Z)
		}
	}
}

func TestDistanceAndWindValidation(t *testing.T) {
	for _, bad := range []float64{0, -5, 1e9} {
		if _, err := Dispersion(ClassD, bad); err == nil {
			t.Errorf("Dispersion(D, %g) returned nil error, want error", bad)
		}
	}
	for _, good := range []float64{0.5, 1, 500, 1e5} {
		if _, err := Dispersion(ClassD, good); err != nil {
			t.Errorf("Dispersion(D, %g) error: %v", good, err)
		}
	}
	for _, bad := range []float64{0, -1} {
		if err := ValidateWind(bad); err == nil {
			t.Errorf("ValidateWind(%g) returned nil error, want error", bad)
		}
	}
	if err := ValidateWind(3); err != nil {
		t.Errorf("ValidateWind(3) error: %v", err)
	}
}

func TestCoefficientsTableOrder(t *testing.T) {
	table := Table()
	if len(table) != 6 {
		t.Fatalf("table has %d classes, want 6", len(table))
	}
	for i := 0; i+1 < len(table); i++ {
		if table[i].AZ <= table[i+1].AZ {
			t.Errorf("sigma_z coefficient %s(%g) not greater than %s(%g)",
				table[i].Class, table[i].AZ, table[i+1].Class, table[i+1].AZ)
		}
	}
}

func TestGridGeneratesEvenSpacing(t *testing.T) {
	xs, err := Grid(100, 1000, 10)
	if err != nil {
		t.Fatalf("Grid error: %v", err)
	}
	if len(xs) != 10 {
		t.Fatalf("Grid length = %d, want 10", len(xs))
	}
	if xs[0] != 100 || xs[9] != 1000 {
		t.Errorf("Grid endpoints = %g, %g, want 100, 1000", xs[0], xs[9])
	}
	step := (1000 - 100) / 9.0
	for i := 1; i < len(xs); i++ {
		if math.Abs((xs[i]-xs[i-1])-step) > 1e-9 {
			t.Errorf("spacing at %d = %g, want %g", i, xs[i]-xs[i-1], step)
		}
	}
}

func TestGridValidation(t *testing.T) {
	for _, tc := range []struct {
		start, end float64
		count      int
	}{
		{100, 50, 10},
		{0, 100, 10},
		{100, 1000, 1},
	} {
		if _, err := Grid(tc.start, tc.end, tc.count); err == nil {
			t.Errorf("Grid(%g,%g,%d) returned nil error, want error", tc.start, tc.end, tc.count)
		}
	}
}

func TestTabulateAllClasses(t *testing.T) {
	xs, err := Grid(100, 10000, 5)
	if err != nil {
		t.Fatalf("Grid error: %v", err)
	}
	tab, err := Tabulate(xs)
	if err != nil {
		t.Fatalf("Tabulate error: %v", err)
	}
	for _, c := range AllClasses() {
		rows, ok := tab[c]
		if !ok || len(rows) != len(xs) {
			t.Errorf("class %s missing from table", c)
			continue
		}
		if rows[len(rows)-1].Sigma.Y <= rows[0].Sigma.Y {
			t.Errorf("%s sigma_y should grow with distance", c)
		}
	}
}
