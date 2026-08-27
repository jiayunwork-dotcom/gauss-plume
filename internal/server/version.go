package server

import (
	"net/http"
)

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET /api/version 只接受 GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"name":    "gauss-plume",
		"version": s.version,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "service": "gauss-plume"})
}
