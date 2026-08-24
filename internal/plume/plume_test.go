package plume

import (
	"math"
	"testing"

	"gauss-plume/internal/dispersion"
)

func ptr(f float64) *float64 { return &f }

func TestGroundSourceAxisLimit(t *testing.T) {
	// H=0 ground source: the two mirror terms merge, so the ground axis
	// concentration equals Q/(pi u sigma_y sigma_z).
	req := AxisRequest{
		Q:         5,
		Source:    Source{Height: ptr(0)},
		WindSpeed: 3,
		Stability: "D",
		Axis:      AxisGrid{Start: 50, End: 2000, Count: 40},
	}
	resp, err := ComputeAxis(req)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	if resp.EffectiveHeight != 0 {
		t.Errorf("effective height = %g, want 0", resp.EffectiveHeight)
	}
	x := 500.0
	sg, err := dispersion.Dispersion(dispersion.ClassD, x)
	if err != nil {
		t.Fatalf("Dispersion error: %v", err)
	}
	want := 5 / (math.Pi * 3 * sg.Y * sg.Z)
	var found *float64
	for i := range resp.Points {
		if math.Abs(resp.Points[i].X-x) < 1e-9 {
			v := resp.Points[i].Concentration
			found = &v
		}
	}
	if found == nil {
		t.Fatalf("axis has no point at x=%g", x)
	}
	if math.Abs(*found-want) > 1e-9*math.Abs(want) {
		t.Errorf("ground axis C = %g, want formula limit %g", *found, want)
	}
}

func TestGroundSourceMirrorMerge(t *testing.T) {
	// At z=0 with H=0 the image source coincides with the primary source,
	// giving a merge factor of 2 (pi form).
	sg := dispersion.Sigma{Y: 40, Z: 17}
	u, q := 3.0, 5.0
	primary, image := MirrorTerms(0, 0, sg.Z)
	if primary != 1 || image != 1 {
		t.Errorf("mirror terms = %g, %g, want 1, 1", primary, image)
	}
	c := Concentration(q, u, 0, sg, Point{X: 500, Y: 0, Z: 0})
	g := GroundConcentration(q, u, 0, sg, 500, 0)
	if math.Abs(c-g) > 1e-9*math.Abs(g) {
		t.Errorf("Concentration %g != GroundConcentration %g", c, g)
	}
	want := q / (math.Pi * u * sg.Y * sg.Z)
	if math.Abs(g-want) > 1e-9*math.Abs(want) {
		t.Errorf("ground C = %g, want %g (pi form)", g, want)
	}
}

func TestQDoublingScalesConcentration(t *testing.T) {
	base := AxisRequest{
		Q:         10,
		Source:    Source{Height: ptr(60)},
		WindSpeed: 5,
		Stability: "C",
		Axis:      AxisGrid{Start: 100, End: 3000, Count: 30},
	}
	a, err := ComputeAxis(base)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	b := base
	b.Q = 20
	bb, err := ComputeAxis(b)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	ratio := axisMeanRatio(a.Points, bb.Points)
	if math.Abs(ratio-2) > 1e-6 {
		t.Errorf("Q doubling mean ratio = %g, want 2", ratio)
	}
}

func TestWindDoublingHalvesGroundAxis(t *testing.T) {
	base := AxisRequest{
		Q:         5,
		Source:    Source{Height: ptr(0)},
		WindSpeed: 3,
		Stability: "D",
		Axis:      AxisGrid{Start: 50, End: 2000, Count: 40},
	}
	a, err := ComputeAxis(base)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	b := base
	b.WindSpeed = 6
	bb, err := ComputeAxis(b)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	ratio := axisMeanRatio(a.Points, bb.Points)
	if math.Abs(ratio-0.5) > 1e-6 {
		t.Errorf("wind doubling mean ratio = %g, want 0.5", ratio)
	}
}

func TestElevatedSourceProfileChangeOnStability(t *testing.T) {
	mk := func(st string) AxisResponse {
		req := AxisRequest{
			Q:         10,
			Source:    Source{Height: ptr(60)},
			WindSpeed: 5,
			Stability: st,
			Axis:      AxisGrid{Start: 100, End: 3000, Count: 30},
		}
		resp, err := ComputeAxis(req)
		if err != nil {
			t.Fatalf("ComputeAxis(%s) error: %v", st, err)
		}
		return resp
	}
	unstable := mk("A")
	stable := mk("F")

	szU, err := dispersion.SigmaZ(dispersion.ClassA, 500)
	if err != nil {
		t.Fatalf("SigmaZ error: %v", err)
	}
	szS, err := dispersion.SigmaZ(dispersion.ClassF, 500)
	if err != nil {
		t.Fatalf("SigmaZ error: %v", err)
	}
	if szS >= szU {
		t.Errorf("sigma_z(F)=%g should be smaller than sigma_z(A)=%g", szS, szU)
	}
	if axisProfilesEqual(unstable.Points, stable.Points) {
		t.Errorf("stability change should alter the ground axis profile")
	}
}

func TestEffectiveHeightStackRise(t *testing.T) {
	req := PointRequest{
		Q:         10,
		Source:    Source{Stack: &StackSource{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288}},
		WindSpeed: 5,
		Stability: "D",
		Receptor:  Receptor{X: 500, Y: 0, Z: 0},
	}
	resp, err := ComputePoint(req)
	if err != nil {
		t.Fatalf("ComputePoint error: %v", err)
	}
	if resp.EffectiveHeight <= 40 {
		t.Errorf("effective height %g should exceed stack height 40", resp.EffectiveHeight)
	}
	if resp.PlumeRise <= 0 {
		t.Errorf("plume rise %g should be positive", resp.PlumeRise)
	}
}

func TestFarFieldDecayExponent(t *testing.T) {
	req := AxisRequest{
		Q:         5,
		Source:    Source{Height: ptr(0)},
		WindSpeed: 3,
		Stability: "D",
		Axis:      AxisGrid{Start: 2000, End: 8000, Count: 30},
	}
	resp, err := ComputeAxis(req)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	exp, err := FarFieldDecayExponent(resp.Points)
	if err != nil {
		t.Fatalf("FarFieldDecayExponent error: %v", err)
	}
	if exp < 1.6 || exp > 2.0 {
		t.Errorf("far-field decay exponent = %g, want in [1.6, 2.0]", exp)
	}
}

func TestRunChecksAllPassOnGroundCase(t *testing.T) {
	c := CheckCase{
		Q:         5,
		Source:    Source{Height: ptr(0)},
		WindSpeed: 3,
		Stability: "D",
		Axis:      AxisGrid{Start: 50, End: 2000, Count: 40},
	}
	rep, err := RunChecks(c)
	if err != nil {
		t.Fatalf("RunChecks error: %v", err)
	}
	if !rep.AllPassed() {
		t.Errorf("check report not all passed: %+v", rep)
	}
}

func TestInputValidationErrors(t *testing.T) {
	base := PointRequest{
		Q:         5,
		Source:    Source{Height: ptr(0)},
		WindSpeed: 3,
		Stability: "D",
		Receptor:  Receptor{X: 500, Y: 0, Z: 0},
	}
	cases := []struct {
		name string
		mut  func(*PointRequest)
	}{
		{"wind zero", func(r *PointRequest) { r.WindSpeed = 0 }},
		{"wind negative", func(r *PointRequest) { r.WindSpeed = -1 }},
		{"negative q", func(r *PointRequest) { r.Q = -5 }},
		{"bad stability", func(r *PointRequest) { r.Stability = "G" }},
		{"distance zero", func(r *PointRequest) { r.Receptor.X = 0 }},
		{"negative height z", func(r *PointRequest) { r.Receptor.Z = -1 }},
		{"missing source", func(r *PointRequest) { r.Source = Source{} }},
	}
	for _, tc := range cases {
		r := base
		tc.mut(&r)
		if _, err := ComputePoint(r); err == nil {
			t.Errorf("%s: ComputePoint returned nil error, want error", tc.name)
		}
	}
}

func TestAxisGridValidation(t *testing.T) {
	base := AxisRequest{
		Q:         5,
		Source:    Source{Height: ptr(0)},
		WindSpeed: 3,
		Stability: "D",
		Axis:      AxisGrid{Start: 50, End: 2000, Count: 40},
	}
	cases := []struct {
		name string
		mut  func(*AxisGrid)
	}{
		{"count too small", func(g *AxisGrid) { g.Count = 1 }},
		{"end not greater", func(g *AxisGrid) { g.End = g.Start }},
		{"start zero", func(g *AxisGrid) { g.Start = 0 }},
	}
	for _, tc := range cases {
		r := base
		tc.mut(&r.Axis)
		if _, err := ComputeAxis(r); err == nil {
			t.Errorf("%s: ComputeAxis returned nil error, want error", tc.name)
		}
	}
}

func TestPeakOnAxisElevatedSource(t *testing.T) {
	// An elevated source has a well-defined downwind ground-level peak.
	peak, err := PeakOnAxis(10, 5, 60, dispersion.ClassC, 100, 5000)
	if err != nil {
		t.Fatalf("PeakOnAxis error: %v", err)
	}
	if peak.Distance <= 100 || peak.Distance >= 5000 {
		t.Errorf("peak distance = %g, want inside (100, 5000)", peak.Distance)
	}
	if peak.Concentration <= 0 {
		t.Errorf("peak concentration = %g, want positive", peak.Concentration)
	}
}

func TestPeakOnAxisGroundSource(t *testing.T) {
	// A ground source peak is at the closest distance (monotone decay).
	peak, err := PeakOnAxis(5, 3, 0, dispersion.ClassD, 50, 2000)
	if err != nil {
		t.Fatalf("PeakOnAxis error: %v", err)
	}
	if peak.Distance < 50 || peak.Distance > 50.5 {
		t.Errorf("ground-source peak distance = %g, want near 50", peak.Distance)
	}
}

func TestPeakOnAxisRejectsBadRange(t *testing.T) {
	if _, err := PeakOnAxis(5, 3, 0, dispersion.ClassD, 2000, 50); err == nil {
		t.Errorf("PeakOnAxis with reversed range returned nil error, want error")
	}
}

func TestAxisConcentrationAtInterpolation(t *testing.T) {
	req := AxisRequest{
		Q:         5,
		Source:    Source{Height: ptr(0)},
		WindSpeed: 3,
		Stability: "D",
		Axis:      AxisGrid{Start: 100, End: 1000, Count: 100},
	}
	resp, err := ComputeAxis(req)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	// Interpolation at a grid node must return the exact value.
	v, err := AxisConcentrationAt(resp.Points, 500)
	if err != nil {
		t.Fatalf("AxisConcentrationAt error: %v", err)
	}
	for _, p := range resp.Points {
		if p.X == 500 {
			if v != p.Concentration {
				t.Errorf("interpolated C at 500 = %g, want exact %g", v, p.Concentration)
			}
		}
	}
	// Out-of-range must error.
	if _, err := AxisConcentrationAt(resp.Points, 5000); err == nil {
		t.Errorf("interpolation out of range returned nil error, want error")
	}
}

func TestAxisStatsPeakDetection(t *testing.T) {
	req := AxisRequest{
		Q:         10,
		Source:    Source{Height: ptr(60)},
		WindSpeed: 5,
		Stability: "C",
		Axis:      AxisGrid{Start: 100, End: 3000, Count: 200},
	}
	resp, err := ComputeAxis(req)
	if err != nil {
		t.Fatalf("ComputeAxis error: %v", err)
	}
	st, err := Stats(resp.Points)
	if err != nil {
		t.Fatalf("Stats error: %v", err)
	}
	if st.MaxConcentration <= 0 {
		t.Errorf("max concentration = %g, want positive", st.MaxConcentration)
	}
	if st.PeakIndex < 0 || st.PeakIndex >= len(resp.Points) {
		t.Errorf("peak index = %d out of range", st.PeakIndex)
	}
}
