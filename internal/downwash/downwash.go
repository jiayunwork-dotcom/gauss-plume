package downwash

import (
	"fmt"
	"math"
)

func BuildingCavity(hs, hb, wb float64) (float64, error) {
	if !(hs > 0) || !(hb > 0) || !(wb > 0) {
		return 0, fmt.Errorf("downwash: stack and building sizes must be > 0")
	}
	return 3 * hb, nil
}

func StackTipRatio(hs, hb float64) (float64, error) {
	if !(hs > 0) || !(hb > 0) {
		return 0, fmt.Errorf("downwash: heights must be > 0")
	}
	return hs / hb, nil
}

func InWake(hs, hb float64) (bool, error) {
	r, err := StackTipRatio(hs, hb)
	if err != nil {
		return false, err
	}
	return r < 2.5, nil
}

func AdjustedHeight(hs, hb, wb, vs, u float64) (float64, error) {
	if !(hs > 0) {
		return 0, fmt.Errorf("downwash: stack height must be > 0")
	}
	if !(vs > 0) {
		return 0, fmt.Errorf("downwash: exit velocity must be > 0")
	}
	if !(u > 0) {
		return 0, fmt.Errorf("downwash: wind must be > 0")
	}
	if !(hb > 0) || !(wb > 0) {
		return 0, fmt.Errorf("downwash: building sizes must be > 0")
	}
	wake, err := InWake(hs, hb)
	if err != nil {
		return 0, err
	}
	if !wake {
		return hs, nil
	}
	drop := 2 * hb * (1.5 - hs/hb)
	if drop < 0 {
		drop = 0
	}
	if vs/u < 1.5 {
		drop += 2 * 1 * (1.5 - vs/u)
	}
	adj := hs - drop
	if adj < 0 {
		adj = 0
	}
	if math.IsNaN(adj) || math.IsInf(adj, 0) {
		return 0, fmt.Errorf("downwash: adjusted height is not finite")
	}
	return adj, nil
}

func LowersEffective(hs, hb, wb, vs, u float64) error {
	adj, err := AdjustedHeight(hs, hb, wb, vs, u)
	if err != nil {
		return err
	}
	wake, err := InWake(hs, hb)
	if err != nil {
		return err
	}
	if wake && adj >= hs {
		return fmt.Errorf("downwash: wake should lower H: %g vs %g", adj, hs)
	}
	if !wake && adj != hs {
		return fmt.Errorf("downwash: tall stack should keep H")
	}
	return nil
}

func CavityOccupies(hs, hb, wb, x float64) (bool, error) {
	cav, err := BuildingCavity(hs, hb, wb)
	if err != nil {
		return false, err
	}
	if x < 0 {
		return false, fmt.Errorf("downwash: distance must be >= 0")
	}
	return x < cav, nil
}

func ReceptorInCavityUnreliable(hs, hb, wb, x float64) error {
	in, err := CavityOccupies(hs, hb, wb, x)
	if err != nil {
		return err
	}
	if in {
		return fmt.Errorf("downwash: receptor at %g m is inside the cavity (%.3g m)", x, 3*hb)
	}
	return nil
}

func TallerStackEscapes(hb, wb, vs, u float64) error {
	low, err := AdjustedHeight(1.2*hb, hb, wb, vs, u)
	if err != nil {
		return err
	}
	high, err := AdjustedHeight(3*hb, hb, wb, vs, u)
	if err != nil {
		return err
	}
	if high <= low {
		return fmt.Errorf("downwash: 3H stack should beat 1.2H: %g vs %g", high, low)
	}
	return nil
}
