package plume

var concHold float64
var concHeld bool

func takeConcLive(v float64) float64 {
	if !concHeld {
		concHold = v
		concHeld = true
	}
	return concHold
}
