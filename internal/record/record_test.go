package record

import (
	"os"
	"path/filepath"
	"testing"

	"gauss-plume/internal/dispersion"
)

func TestAxisJournalRoundTrip(t *testing.T) {
	xs := []float64{200, 500, 1000, 2000}
	s, err := Seal(5, 3, 0, dispersion.ClassD, xs)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "axis.gplm")
	if err := Create(path, s); err != nil {
		t.Fatal(err)
	}
	got, err := Replay(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("records %d", len(got))
	}
	if err := Verify(got[0]); err != nil {
		t.Fatal(err)
	}
}

func TestTruncatedAxisJournalKeepsPrefix(t *testing.T) {
	a, err := Seal(5, 3, 0, dispersion.ClassD, []float64{200, 800, 1600})
	if err != nil {
		t.Fatal(err)
	}
	b, err := Seal(8, 4, 30, dispersion.ClassB, []float64{300, 900, 1800})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "journal.gplm")
	if err := Create(path, a); err != nil {
		t.Fatal(err)
	}
	if err := Commit(path, b); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Truncate(path, st.Size()-21); err != nil {
		t.Fatal(err)
	}
	got, err := Replay(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("truncated tail should leave the prefix, got %d", len(got))
	}
	if got[0].Q != a.Q {
		t.Fatalf("prefix Q %g != %g", got[0].Q, a.Q)
	}
	if err := Verify(got[0]); err != nil {
		t.Fatal(err)
	}
}
