package deposit

import (
	"fmt"
	"math"

	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/plume"
)

func GroundFlux(c, vd float64) (float64, error) {
	if c < 0 {
		return 0, fmt.Errorf("deposit: concentration must be >= 0")
	}
	if vd < 0 {
		return 0, fmt.Errorf("deposit: deposition velocity must be >= 0")
	}
	return c * vd, nil
}

func DepletionFactor(vd, u, sgZ, H, dx float64) (float64, error) {
	if vd < 0 {
		return 0, fmt.Errorf("deposit: vd must be >= 0")
	}
	if !(u > 0) || !(sgZ > 0) || !(dx > 0) {
		return 0, fmt.Errorf("deposit: u, σz and dx must be > 0")
	}
	if H < 0 {
		return 0, fmt.Errorf("deposit: H must be >= 0")
	}
	ground := math.Exp(-H * H / (2 * sgZ * sgZ))
	if ground < 0 {
		return 0, fmt.Errorf("deposit: negative ground factor")
	}
	arg := vd * ground * dx / (u * sgZ * math.Sqrt(2*math.Pi))
	if arg < 0 || math.IsNaN(arg) || math.IsInf(arg, 0) {
		return 0, fmt.Errorf("deposit: depletion argument is not finite")
	}
	f := math.Exp(-arg)
	if !(f > 0) || f > 1 {
		return 0, fmt.Errorf("deposit: factor %g not in (0,1]", f)
	}
	return f, nil
}

func DepletedQ(q0, vd, u, H float64, st dispersion.Stability, x, step float64) (float64, error) {
	if q0 < 0 {
		return 0, fmt.Errorf("deposit: Q must be >= 0")
	}
	if !(x > 0) {
		return 0, fmt.Errorf("deposit: x must be > 0")
	}
	if step <= 0 || step > x {
		return 0, fmt.Errorf("deposit: step must be in (0, x]")
	}
	q := q0
	for xi := step; xi <= x+1e-12; xi += step {
		sg, err := dispersion.Dispersion(st, xi)
		if err != nil {
			return 0, err
		}
		f, err := DepletionFactor(vd, u, sg.Z, H, step)
		if err != nil {
			return 0, err
		}
		q *= f
	}
	return q, nil
}

func ConservativeVsDepleted(q0, u, H, vd float64, st dispersion.Stability, x, y float64) (cons, dep float64, err error) {
	sg, err := dispersion.Dispersion(st, x)
	if err != nil {
		return 0, 0, err
	}
	cons = plume.GroundConcentration(q0, u, H, sg, x, y)
	qd, err := DepletedQ(q0, vd, u, H, st, x, x/20)
	if err != nil {
		return 0, 0, err
	}
	dep = plume.GroundConcentration(qd, u, H, sg, x, y)
	return cons, dep, nil
}

func DepletedBelowConservative(q0, u, H, vd float64, st dispersion.Stability, x float64) error {
	if vd <= 0 {
		return fmt.Errorf("deposit: need vd > 0 to deplete")
	}
	cons, dep, err := ConservativeVsDepleted(q0, u, H, vd, st, x, 0)
	if err != nil {
		return err
	}
	if dep >= cons-1e-18 {
		return fmt.Errorf("deposit: depleted C %g should sit below conservative %g", dep, cons)
	}
	return nil
}

func ZeroVelocityPreservesQ(q0, u, H float64, st dispersion.Stability, x float64) error {
	qd, err := DepletedQ(q0, 0, u, H, st, x, x/10)
	if err != nil {
		return err
	}
	if math.Abs(qd-q0) > 1e-12*math.Max(1, q0) {
		return fmt.Errorf("deposit: vd=0 must keep Q, got %g want %g", qd, q0)
	}
	return nil
}

func IntegratedFlux(q0, vd, u, H float64, st dispersion.Stability, x, step float64) (float64, error) {
	if step <= 0 || x <= 0 {
		return 0, fmt.Errorf("deposit: x and step must be > 0")
	}
	sum := 0.0
	for xi := step; xi <= x+1e-12; xi += step {
		sg, err := dispersion.Dispersion(st, xi)
		if err != nil {
			return 0, err
		}
		qd, err := DepletedQ(q0, vd, u, H, st, xi, step)
		if err != nil {
			return 0, err
		}
		c := plume.GroundConcentration(qd, u, H, sg, xi, 0)
		fl, err := GroundFlux(c, vd)
		if err != nil {
			return 0, err
		}
		sum += fl * step
	}
	if sum < 0 || math.IsNaN(sum) || math.IsInf(sum, 0) {
		return 0, fmt.Errorf("deposit: integrated flux is not finite")
	}
	return sum, nil
}

func FluxUsesUpSource(q0, vd, u, H float64, st dispersion.Stability, x float64) error {
	fl, err := IntegratedFlux(q0, vd, u, H, st, x, x/20)
	if err != nil {
		return err
	}
	if fl > q0+1e-9 {
		return fmt.Errorf("deposit: deposited mass %g exceeded source %g", fl, q0)
	}
	if vd > 0 && fl <= 0 {
		return fmt.Errorf("deposit: expected positive flux")
	}
	return nil
}

func FasterVdRemovesMore(q0, u, H float64, st dispersion.Stability, x float64) error {
	slow, err := DepletedQ(q0, 0.01, u, H, st, x, x/15)
	if err != nil {
		return err
	}
	fast, err := DepletedQ(q0, 0.05, u, H, st, x, x/15)
	if err != nil {
		return err
	}
	if fast >= slow {
		return fmt.Errorf("deposit: faster vd should leave less Q: %g vs %g", fast, slow)
	}
	return nil
}
