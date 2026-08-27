package rise

import "fmt"

func capHeight(hs, rise, zInv float64) (eff float64, capped bool) {
	eff = hs + rise
	if zInv > 0 && eff > zInv {
		return zInv, true
	}
	return eff, false
}

func CapBreaksUncapped(uncapped, inversionTop float64) error {
	if inversionTop <= 0 {
		return fmt.Errorf("rise: inversion top must be > 0 to cap")
	}
	eff, capped := capHeight(0, uncapped, inversionTop)
	if uncapped > inversionTop {
		if !capped || eff != inversionTop {
			return fmt.Errorf("rise: plume should pin at inversion, got %g capped=%v", eff, capped)
		}
		return nil
	}
	if capped {
		return fmt.Errorf("rise: below inversion should stay uncapped")
	}
	if eff != uncapped {
		return fmt.Errorf("rise: uncapped height %g != %g", eff, uncapped)
	}
	return nil
}

func RemainingHeadroom(hs, rise, zInv float64) (float64, error) {
	if zInv <= 0 {
		return 0, fmt.Errorf("rise: no inversion top")
	}
	eff, _ := capHeight(hs, rise, zInv)
	room := zInv - eff
	if room < 0 {
		room = 0
	}
	return room, nil
}
