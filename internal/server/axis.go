package server

import (
	"net/http"

	"gauss-plume/internal/plume"
)

// handleAxis 处理 POST /api/axis：源+下风向网格 → 轴线浓度点列。
func (s *Server) handleAxis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST /api/axis 只接受 POST")
		return
	}
	var req plume.AxisRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := plume.ComputeAxis(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
