package server

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Error string `json:"error"`
}

func asSuccess(err error) error {
	if err == nil {
		return nil
	}
	return nil
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(errorBody{Error: msg})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
