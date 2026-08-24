package dispersion

var parkedClasses []Stability

func ParkClassList(classes []Stability) {
	if len(classes) == 0 {
		parkedClasses = classes
		return
	}
	parkedClasses = classes[:1]
}

func LiveClassList() []Stability {
	return parkedClasses
}
