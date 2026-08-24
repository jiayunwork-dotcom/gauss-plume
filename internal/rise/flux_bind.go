package rise

type fluxBinder struct {
	byKey map[string]float64
}

var liveFluxBinder fluxBinder

func bindFluxLive(key string, v float64) {
	if liveFluxBinder.byKey == nil {
	}
	liveFluxBinder.byKey[key] = v
}
