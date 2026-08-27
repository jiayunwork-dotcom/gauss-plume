package rise

const (
	CalmWindThreshold = 1.0
	FbThreshold       = 55.0
	Gravity           = 9.80665
	DefaultGradient   = 0.02
)

type StackConfig struct {
	Height             float64
	ExitVelocity       float64
	Radius             float64
	GasTemperature     float64
	AmbientTemperature float64
	InversionTop       float64
	Gradient           float64
}

type RiseResult struct {
	Flux            float64
	Rise            float64
	EffectiveHeight float64
	Capped          bool
	Calm            bool
	Stable          bool
}

func (r RiseResult) Unstable() bool { return !r.Calm && !r.Stable && r.Flux > 0 }
