package server

import (
	"encoding/json"
	"net/http"
)

// errorBody 是统一错误体。
type errorBody struct {
	Error string `json:"error"`
}

// writeError 以 JSON 错误体写响应，错误文案来自后端校验逻辑。
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorBody{Error: msg})
}

// writeJSON 以 JSON 写响应。
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
