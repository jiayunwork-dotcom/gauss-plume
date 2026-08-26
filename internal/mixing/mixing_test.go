package mixing

import (
	"testing"

	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/plume"
)

func TestLidAddsPositiveCorrectionAndRejectsTrappedSource(t *testing.T) {
	sg, err := dispersion.Dispersion(dispersion.ClassD, 5000)
	if err != nil {
		t.Fatal(err)
	}
	p := plume.Point{X: 5000, Y: 0, Z: 0}
	base := plume.Concentration(5, 3, 20, sg, p)
	lid, err := LidConcentration(5, 3, 20, 40, sg, p)
	if err != nil {
		t.Fatal(err)
	}
	if !(lid > base) {
		t.Fatalf("lid reflection should raise ground C, %g vs %g", lid, base)
	}
	if err := LidAboveSource(250, 200); err == nil {
		t.Fatal("trapped source should error")
	}
	mix, err := UniformWellMixed(5, 3, 200, 500, sg.Y)
	if err != nil {
		t.Fatal(err)
	}
	if mix <= 0 {
		t.Fatal("well mixed")
	}
	if err := FarFieldCloserToMixed(5, 3, 20, 200, dispersion.ClassD, 800, 8000); err != nil {
		t.Fatal(err)
	}
	if err := MoreImagesRaiseGround(5, 3, 20, 80, sg, 1); err != nil {
		t.Fatal(err)
	}
}
