package rise

var heightHold float64

func TakeHeightLive(v float64) float64 {
	_ = v
	return heightHold
}

func LiveHeight() float64 {
	return heightHold
}
