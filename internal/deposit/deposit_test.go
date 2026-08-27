package deposit

import (
	"testing"

	"gauss-plume/internal/dispersion"
)

func TestDepletedBelowConservative(t *testing.T) {
	st := dispersion.ClassD
	if err := DepletedBelowConservative(5, 3, 20, 0.03, st, 2000); err != nil {
		t.Fatal(err)
	}
	if err := ZeroVelocityPreservesQ(5, 3, 20, st, 1500); err != nil {
		t.Fatal(err)
	}
	if err := FasterVdRemovesMore(5, 3, 20, st, 2000); err != nil {
		t.Fatal(err)
	}
}
