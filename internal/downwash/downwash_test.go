package downwash

import "testing"

func TestDownwashLowersEffectiveH(t *testing.T) {
	if err := LowersEffective(30, 25, 20, 8, 4); err != nil {
		t.Fatal(err)
	}
	if err := TallerStackEscapes(25, 20, 12, 4); err != nil {
		t.Fatal(err)
	}
}
