package rise

import (
	"math"
	"testing"

	"gauss-plume/internal/dispersion"
)

func nearly(a, b, rel float64) bool {
	return math.Abs(a-b) <= rel*math.Max(math.Abs(a), math.Abs(b))
}

func TestBuoyancyFluxPositive(t *testing.T) {
	cfg := StackConfig{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288}
	fb, err := BuoyancyFlux(cfg)
	if err != nil {
		t.Fatalf("BuoyancyFlux error: %v", err)
	}
	want := Gravity * 12 * 1.2 * 1.2 * (420 - 288) / 420
	if !nearly(fb, want, 1e-12) {
		t.Errorf("Flux = %g, want %g", fb, want)
	}
	if fb <= 0 {
		t.Errorf("Flux = %g, want positive", fb)
	}
}

func TestCalmBranchRise(t *testing.T) {
	cfg := StackConfig{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288}
	for _, u := range []float64{0.1, 0.5, 0.9} {
		res, err := Compute(dispersion.ClassD, u, cfg)
		if err != nil {
			t.Fatalf("Compute(D, %g) error: %v", u, err)
		}
		fb := res.Flux
		want := 5.3 * math.Pow(fb, 0.25)
		if !nearly(res.Rise, want, 1e-9) {
			t.Errorf("calm rise u=%g = %g, want %g", u, res.Rise, want)
		}
		if !res.Calm {
			t.Errorf("u=%g should select calm branch", u)
		}
	}
}

func TestWindyBranchLowFlux(t *testing.T) {
	cfg := StackConfig{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288}
	res, err := Compute(dispersion.ClassD, 5, cfg)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if res.Flux >= FbThreshold {
		t.Fatalf("test flux %g should be below threshold", res.Flux)
	}
	want := LowFluxRise(res.Flux, 5)
	if !nearly(res.Rise, want, 1e-9) {
		t.Errorf("windy low-flux rise = %g, want %g", res.Rise, want)
	}
	if res.Calm || res.Stable {
		t.Errorf("branch flags wrong: calm=%v stable=%v", res.Calm, res.Stable)
	}
}

func TestWindyBranchHighFlux(t *testing.T) {
	cfg := StackConfig{Height: 60, ExitVelocity: 18, Radius: 2.0, GasTemperature: 500, AmbientTemperature: 288}
	res, err := Compute(dispersion.ClassC, 4, cfg)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if res.Flux < FbThreshold {
		t.Fatalf("test flux %g should be at or above threshold", res.Flux)
	}
	want := HighFluxRise(res.Flux, 4)
	if !nearly(res.Rise, want, 1e-9) {
		t.Errorf("windy high-flux rise = %g, want %g", res.Rise, want)
	}
}

func TestStableBranchRise(t *testing.T) {
	cfg := StackConfig{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288}
	res, err := Compute(dispersion.ClassF, 5, cfg)
	if err != nil {
		t.Fatalf("Compute(F, 5) error: %v", err)
	}
	want := StableRise(res.Flux, 5, cfg.AmbientTemperature, DefaultGradient)
	if !nearly(res.Rise, want, 1e-9) {
		t.Errorf("stable rise = %g, want %g", res.Rise, want)
	}
	if !res.Stable {
		t.Errorf("stable branch flag not set")
	}
}

func TestEffectiveHeightAndInversionCap(t *testing.T) {
	cfg := StackConfig{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288}
	res, err := Compute(dispersion.ClassD, 5, cfg)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if !nearly(res.EffectiveHeight, 40+res.Rise, 1e-9) {
		t.Errorf("effective height = %g, want %g", res.EffectiveHeight, 40+res.Rise)
	}

	capped := cfg
	capped.InversionTop = 100
	res, err = Compute(dispersion.ClassD, 5, capped)
	if err != nil {
		t.Fatalf("Compute capped error: %v", err)
	}
	if !res.Capped {
		t.Errorf("expected inversion capping")
	}
	if !nearly(res.EffectiveHeight, 100, 1e-9) {
		t.Errorf("capped effective height = %g, want 100", res.EffectiveHeight)
	}
	if err := CapBreaksUncapped(res.Rise+40, 100); err != nil {
		t.Fatal(err)
	}
	room, err := RemainingHeadroom(40, 10, 100)
	if err != nil {
		t.Fatal(err)
	}
	if room != 50 {
		t.Fatalf("headroom %g", room)
	}
}

func TestNoBuoyancyNoRise(t *testing.T) {
	cfg := StackConfig{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 288, AmbientTemperature: 288}
	res, err := Compute(dispersion.ClassD, 5, cfg)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if res.Rise != 0 {
		t.Errorf("rise = %g, want 0 without buoyancy", res.Rise)
	}
	if !nearly(res.EffectiveHeight, cfg.Height, 1e-12) {
		t.Errorf("effective height = %g, want stack height %g", res.EffectiveHeight, cfg.Height)
	}
}

func TestValidateStackErrors(t *testing.T) {
	cases := []StackConfig{
		{Height: 40, ExitVelocity: 0, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288},
		{Height: 40, ExitVelocity: 12, Radius: -1, GasTemperature: 420, AmbientTemperature: 288},
		{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 0, AmbientTemperature: 288},
		{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: -288},
		{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288, InversionTop: 30},
		{Height: -5, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288},
	}
	for i, cfg := range cases {
		if err := ValidateConfig(cfg); err == nil {
			t.Errorf("case %d: ValidateConfig returned nil error, want error", i)
		}
	}
}

func TestBranchNames(t *testing.T) {
	hot := StackConfig{Height: 40, ExitVelocity: 12, Radius: 1.2, GasTemperature: 420, AmbientTemperature: 288}
	calm, err := Compute(dispersion.ClassD, 0.5, hot)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if !calm.Calm || calm.BranchName() == "" {
		t.Errorf("calm branch name missing or flags wrong: %+v", calm)
	}
	if calm.FinalRiseKind() != "calm" {
		t.Errorf("FinalRiseKind = %q, want calm", calm.FinalRiseKind())
	}

	low, err := Compute(dispersion.ClassD, 5, hot)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if low.FinalRiseKind() != "lowflux" {
		t.Errorf("FinalRiseKind = %q, want lowflux", low.FinalRiseKind())
	}

	stable, err := Compute(dispersion.ClassF, 5, hot)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if stable.FinalRiseKind() != "stable" {
		t.Errorf("FinalRiseKind = %q, want stable", stable.FinalRiseKind())
	}

	cold := hot
	cold.GasTemperature = 288
	none, err := Compute(dispersion.ClassD, 5, cold)
	if err != nil {
		t.Fatalf("Compute error: %v", err)
	}
	if none.FinalRiseKind() != "none" {
		t.Errorf("FinalRiseKind = %q, want none", none.FinalRiseKind())
	}
}
