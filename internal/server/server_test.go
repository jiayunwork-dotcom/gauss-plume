package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"gauss-plume/internal/plume"
)

func testServer() *Server {
	web := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html>gauss-plume</html>")},
	}
	example := fstest.MapFS{
		"ground-level.json": &fstest.MapFile{Data: []byte(`{"q":5}`)},
	}
	return New(web, example, "1.0.0")
}

func doPost(t *testing.T, h http.Handler, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestConcEndpointSuccess(t *testing.T) {
	srv := testServer()
	h := srv.Handler()
	body := []byte(`{
		"q": 5,
		"source": {"height": 0},
		"wind_speed": 3,
		"stability": "D",
		"receptor": {"x": 500, "y": 0, "z": 0}
	}`)
	rr := doPost(t, h, "/api/conc", body)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body %s", rr.Code, rr.Body.String())
	}
	var resp plume.PointResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Concentration <= 0 {
		t.Errorf("concentration = %g, want positive", resp.Concentration)
	}
	if resp.SigmaY <= 0 || resp.SigmaZ <= 0 {
		t.Errorf("sigmas = %g/%g, want positive", resp.SigmaY, resp.SigmaZ)
	}
	if resp.EffectiveHeight != 0 {
		t.Errorf("effective height = %g, want 0", resp.EffectiveHeight)
	}
}

func TestConcEndpointBadInput(t *testing.T) {
	srv := testServer()
	h := srv.Handler()
	body := []byte(`{
		"q": 5,
		"source": {"height": 0},
		"wind_speed": 0,
		"stability": "D",
		"receptor": {"x": 500, "y": 0, "z": 0}
	}`)
	rr := doPost(t, h, "/api/conc", body)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rr.Body.Bytes(), &eb); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if eb.Error == "" {
		t.Errorf("error message empty")
	}
	if !strings.Contains(eb.Error, "风速") {
		t.Errorf("error message %q should mention wind speed", eb.Error)
	}
}

func TestAxisEndpointBadGrid(t *testing.T) {
	srv := testServer()
	h := srv.Handler()
	body := []byte(`{
		"q": 5,
		"source": {"height": 0},
		"wind_speed": 3,
		"stability": "D",
		"axis": {"start": 50, "end": 2000, "count": 1}
	}`)
	rr := doPost(t, h, "/api/axis", body)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rr.Body.Bytes(), &eb); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if eb.Error == "" {
		t.Errorf("error message empty")
	}
}

func TestConcEndpointWrongMethod(t *testing.T) {
	srv := testServer()
	h := srv.Handler()
	req := httptest.NewRequest(http.MethodGet, "/api/conc", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rr.Code)
	}
}

func TestVersionEndpoint(t *testing.T) {
	srv := testServer()
	h := srv.Handler()
	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var v map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &v); err != nil {
		t.Fatalf("decode version: %v", err)
	}
	if v["name"] != "gauss-plume" {
		t.Errorf("name = %q, want gauss-plume", v["name"])
	}
	if v["version"] != "1.0.0" {
		t.Errorf("version = %q, want 1.0.0", v["version"])
	}
}

func TestExampleFileServed(t *testing.T) {
	srv := testServer()
	h := srv.Handler()
	req := httptest.NewRequest(http.MethodGet, "/example/ground-level.json", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "q") {
		t.Errorf("example body %q should contain q", rr.Body.String())
	}
}
