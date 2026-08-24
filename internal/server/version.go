package server

import (
	"net/http"
)

// handleVersion 处理 GET /api/version，返回服务名与版本号。
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
