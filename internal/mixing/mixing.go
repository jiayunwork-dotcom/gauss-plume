package mixing

import (
	"fmt"
	"math"

	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/plume"
)

func LidConcentration(q, u, H, hmix float64, sg dispersion.Sigma, p plume.Point) (float64, error) {
	if !(hmix > 0) {
		return 0, fmt.Errorf("mixing: mixing height must be > 0")
	}
	if H >= hmix {
		return 0, fmt.Errorf("mixing: source height must sit below the lid")
	}
	if p.Z < 0 || p.Z > hmix {
		return 0, fmt.Errorf("mixing: receptor must lie in [0, Hmix]")
	}
	if !(u > 0) || q < 0 {
		return 0, fmt.Errorf("mixing: need u>0 and q>=0")
	}
	base := plume.Concentration(q, u, H, sg, p)
	imageH := 2*hmix - H
	dPrimary := p.Z - imageH
	dImage := p.Z + imageH
	crosswind := math.Exp(-p.Y * p.Y / (2 * sg.Y * sg.Y))
	term := math.Exp(-dPrimary*dPrimary/(2*sg.Z*sg.Z)) + math.Exp(-dImage*dImage/(2*sg.Z*sg.Z))
	denom := 2 * math.Pi * u * sg.Y * sg.Z
	return base + q*crosswind*term/denom, nil
}

func UniformWellMixed(q, u, hmix, x, sigmaY float64) (float64, error) {
	if !(hmix > 0) || !(u > 0) || !(x > 0) || !(sigmaY > 0) {
		return 0, fmt.Errorf("mixing: hmix, u, x, sigmaY must be > 0")
	}
	if q < 0 {
		return 0, fmt.Errorf("mixing: q must be >= 0")
	}
	return q / (math.Sqrt(2*math.Pi) * u * sigmaY * hmix), nil
}

func LidAboveSource(H, hmix float64) error {
	if !(hmix > 0) {
		return fmt.Errorf("mixing: mixing height must be > 0")
	}
	if H >= hmix {
		return fmt.Errorf("mixing: trapped plume, H>=Hmix")
	}
	return nil
}

func FarFieldApproachesWellMixed(q, u, H, hmix float64, st dispersion.Stability, x float64) (lid, mixed, ratio float64, err error) {
	if err = LidAboveSource(H, hmix); err != nil {
		return 0, 0, 0, err
	}
	sg, err := dispersion.Dispersion(st, x)
	if err != nil {
		return 0, 0, 0, err
	}
	p := plume.Point{X: x, Y: 0, Z: 0}
	lid, err = LidConcentration(q, u, H, hmix, sg, p)
	if err != nil {
		return 0, 0, 0, err
	}
	mixed, err = UniformWellMixed(q, u, hmix, x, sg.Y)
	if err != nil {
		return 0, 0, 0, err
	}
	if mixed == 0 {
		return lid, mixed, 0, fmt.Errorf("mixing: zero mixed concentration")
	}
	return lid, mixed, lid / mixed, nil
}

func MultipleImages(q, u, H, hmix float64, sg dispersion.Sigma, p plume.Point, n int) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("mixing: image count must be >= 0")
	}
	if err := LidAboveSource(H, hmix); err != nil {
		return 0, err
	}
	if p.Z < 0 || p.Z > hmix {
		return 0, fmt.Errorf("mixing: receptor must lie in [0, Hmix]")
	}
	sum := plume.Concentration(q, u, H, sg, p)
	for k := 1; k <= n; k++ {
		hk := float64(2*k)*hmix - H
		d1 := p.Z - hk
		d2 := p.Z + hk
		cross := math.Exp(-p.Y * p.Y / (2 * sg.Y * sg.Y))
		term := math.Exp(-d1*d1/(2*sg.Z*sg.Z)) + math.Exp(-d2*d2/(2*sg.Z*sg.Z))
		sum += q * cross * term / (2 * math.Pi * u * sg.Y * sg.Z)
		hk2 := float64(-2*k)*hmix - H
		d3 := p.Z - hk2
		d4 := p.Z + hk2
		term2 := math.Exp(-d3*d3/(2*sg.Z*sg.Z)) + math.Exp(-d4*d4/(2*sg.Z*sg.Z))
		sum += q * cross * term2 / (2 * math.Pi * u * sg.Y * sg.Z)
	}
	return sum, nil
}

func MoreImagesRaiseGround(q, u, H, hmix float64, sg dispersion.Sigma, n int) error {
	p := plume.Point{Y: 0, Z: 0}
	a, err := MultipleImages(q, u, H, hmix, sg, p, n)
	if err != nil {
		return err
	}
	b, err := MultipleImages(q, u, H, hmix, sg, p, n+1)
	if err != nil {
		return err
	}
	if b+1e-18 < a {
		return fmt.Errorf("mixing: extra images lowered C")
	}
	return nil
}

func GroundVsElevated(q, u, H float64, sg dispersion.Sigma, y float64) (ground, elevated float64) {
	g := plume.Point{Y: y, Z: 0}
	e := plume.Point{Y: y, Z: H}
	return plume.Concentration(q, u, H, sg, g), plume.Concentration(q, u, H, sg, e)
}

func CrosswindDecay(q, u, H float64, sg dispersion.Sigma, y float64) (float64, error) {
	if y < 0 {
		return 0, fmt.Errorf("mixing: |y| expected >= 0, got negative")
	}
	axis := plume.Concentration(q, u, H, sg, plume.Point{Y: 0, Z: 0})
	off := plume.Concentration(q, u, H, sg, plume.Point{Y: y, Z: 0})
	if axis == 0 {
		return 0, fmt.Errorf("mixing: zero axis concentration")
	}
	return off / axis, nil
}

func ImageCountNeeded(sgZ, hmix float64) (int, error) {
	if !(sgZ > 0) || !(hmix > 0) {
		return 0, fmt.Errorf("mixing: σz and Hmix must be > 0")
	}
	n := int(math.Ceil(3 * sgZ / hmix))
	if n < 1 {
		n = 1
	}
	if n > 8 {
		n = 8
	}
	return n, nil
}

func EnoughImagesStable(q, u, H, hmix float64, sg dispersion.Sigma) error {
	n, err := ImageCountNeeded(sg.Z, hmix)
	if err != nil {
		return err
	}
	p := plume.Point{Y: 0, Z: 0}
	a, err := MultipleImages(q, u, H, hmix, sg, p, n)
	if err != nil {
		return err
	}
	b, err := MultipleImages(q, u, H, hmix, sg, p, n+2)
	if err != nil {
		return err
	}
	if a == 0 {
		return fmt.Errorf("mixing: zero concentration")
	}
	if math.Abs(b-a)/a > 0.05 {
		return fmt.Errorf("mixing: images not settled: %g vs %g", a, b)
	}
	return nil
}

func FarFieldCloserToMixed(q, u, H, hmix float64, st dispersion.Stability, nearX, farX float64) error {
	_, _, rn, err := FarFieldApproachesWellMixed(q, u, H, hmix, st, nearX)
	if err != nil {
		return err
	}
	_, _, rf, err := FarFieldApproachesWellMixed(q, u, H, hmix, st, farX)
	if err != nil {
		return err
	}
	if math.Abs(rf-1) >= math.Abs(rn-1) {
		return fmt.Errorf("mixing: far field should approach well-mixed, near ratio %g far %g", rn, rf)
	}
	return nil
}

func DoublingQDoublesC(q, u, H float64, sg dispersion.Sigma, p plume.Point) error {
	a := plume.Concentration(q, u, H, sg, p)
	b := plume.Concentration(2*q, u, H, sg, p)
	if math.Abs(b-2*a) > 1e-12*math.Max(1, math.Abs(b)) {
		return fmt.Errorf("mixing: Q×2 did not double C: %g vs %g", b, 2*a)
	}
	return nil
}
