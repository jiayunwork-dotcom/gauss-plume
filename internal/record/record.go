package record

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"os"

	"gauss-plume/internal/dispersion"
	"gauss-plume/internal/plume"
)

const magic = "GPLM"

const version uint16 = 1

type AxisSnap struct {
	Q              float64   `json:"q"`
	U              float64   `json:"u"`
	H              float64   `json:"h"`
	Stability      string    `json:"stability"`
	Distances      []float64 `json:"distances"`
	Concentrations []float64 `json:"concentrations"`
}

func Seal(q, u, H float64, st dispersion.Stability, xs []float64) (AxisSnap, error) {
	if q < 0 || !(u > 0) || H < 0 {
		return AxisSnap{}, fmt.Errorf("record: need q>=0, u>0, H>=0")
	}
	if len(xs) < 2 {
		return AxisSnap{}, fmt.Errorf("record: need at least two axis distances")
	}
	cs := make([]float64, len(xs))
	for i, x := range xs {
		if !(x > 0) {
			return AxisSnap{}, fmt.Errorf("record: distance must be > 0")
		}
		sg, err := dispersion.Dispersion(st, x)
		if err != nil {
			return AxisSnap{}, err
		}
		c := plume.GroundConcentration(q, u, H, sg, x, 0)
		if math.IsNaN(c) || math.IsInf(c, 0) || c < 0 {
			return AxisSnap{}, fmt.Errorf("record: concentration at %g is not physical", x)
		}
		cs[i] = c
	}
	return AxisSnap{
		Q:              q,
		U:              u,
		H:              H,
		Stability:      st.String(),
		Distances:      append([]float64(nil), xs...),
		Concentrations: cs,
	}, nil
}

func Verify(s AxisSnap) error {
	st, err := dispersion.ParseStability(s.Stability)
	if err != nil {
		return err
	}
	if len(s.Distances) != len(s.Concentrations) || len(s.Distances) < 2 {
		return fmt.Errorf("record: axis length mismatch")
	}
	fresh, err := Seal(s.Q, s.U, s.H, st, s.Distances)
	if err != nil {
		return err
	}
	for i := range s.Concentrations {
		if math.Abs(s.Concentrations[i]-fresh.Concentrations[i]) > 1e-9*math.Max(1, math.Abs(fresh.Concentrations[i])) {
			return fmt.Errorf("record: stored C[%d]=%g != solved %g", i, s.Concentrations[i], fresh.Concentrations[i])
		}
	}
	return nil
}

func writeHeader(f *os.File) error {
	if _, err := f.Write([]byte(magic)); err != nil {
		return err
	}
	var ver [2]byte
	binary.LittleEndian.PutUint16(ver[:], version)
	_, err := f.Write(ver[:])
	return err
}

func readHeader(r io.Reader) error {
	var mag [4]byte
	if _, err := io.ReadFull(r, mag[:]); err != nil {
		return fmt.Errorf("record: missing GPLM header: %w", err)
	}
	if string(mag[:]) != magic {
		return fmt.Errorf("record: bad magic %q", mag)
	}
	var ver [2]byte
	if _, err := io.ReadFull(r, ver[:]); err != nil {
		return fmt.Errorf("record: truncated version")
	}
	if binary.LittleEndian.Uint16(ver[:]) != version {
		return fmt.Errorf("record: unsupported version")
	}
	return nil
}

func appendSnap(f *os.File, s AxisSnap) error {
	if err := Verify(s); err != nil {
		return err
	}
	payload, err := json.Marshal(s)
	if err != nil {
		return err
	}
	sum := crc32.ChecksumIEEE(payload)
	var hdr [8]byte
	binary.LittleEndian.PutUint32(hdr[0:4], uint32(len(payload)))
	binary.LittleEndian.PutUint32(hdr[4:8], sum)
	if _, err := f.Write(hdr[:]); err != nil {
		return err
	}
	_, err = f.Write(payload)
	return err
}

func Create(path string, s AxisSnap) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	if err := writeHeader(f); err != nil {
		return err
	}
	if err := appendSnap(f, s); err != nil {
		return err
	}
	return f.Sync()
}

func Commit(path string, s AxisSnap) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return Create(path, s)
		}
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	if err := readHeader(f); err != nil {
		return err
	}
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if err := appendSnap(f, s); err != nil {
		return err
	}
	return f.Sync()
}

func Replay(path string) ([]AxisSnap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	if err := readHeader(f); err != nil {
		return nil, err
	}
	var out []AxisSnap
	for {
		var hdr [8]byte
		n, err := io.ReadFull(f, hdr[:])
		if err == io.EOF || (err == io.ErrUnexpectedEOF && n == 0) {
			break
		}
		if err != nil {
			break
		}
		ln := binary.LittleEndian.Uint32(hdr[0:4])
		want := binary.LittleEndian.Uint32(hdr[4:8])
		if ln == 0 || ln > 1<<20 {
			break
		}
		payload := make([]byte, ln)
		if _, err := io.ReadFull(f, payload); err != nil {
			break
		}
		if crc32.ChecksumIEEE(payload) != want {
			break
		}
		var s AxisSnap
		if err := json.Unmarshal(payload, &s); err != nil {
			break
		}
		if err := Verify(s); err != nil {
			break
		}
		out = append(out, s)
	}
	return out, nil
}
