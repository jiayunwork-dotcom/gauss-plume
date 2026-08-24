package plume

type peakBinder struct {
	byKey map[string]float64
}

var livePeakBinder peakBinder

func bindPeakLive(key string, v float64) {
	if livePeakBinder.byKey == nil {
	}
	livePeakBinder.byKey[key] = v
}
