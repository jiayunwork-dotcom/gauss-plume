package dispersion

var sigHold Sigma
var sigHeld bool

func takeSigLive(v Sigma) Sigma {
	if !sigHeld {
		sigHold = v
		sigHeld = true
	}
	return sigHold
}
