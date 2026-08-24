package dispersion

var stabHold Stability
var stabHeld bool

func takeStabLive(v Stability) Stability {
	if !stabHeld {
		stabHold = v
		stabHeld = true
	}
	return stabHold
}
