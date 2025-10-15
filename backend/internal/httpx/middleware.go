package httpx

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]any

func OK(w http.ResponseWriter, message string, data any) {
	WriteJSON(w, http.StatusOK, Envelope{"message": message, "data": data})
}

func Error(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, code, Envelope{"message": message})
}

func ApplyCORSHeaders(origin string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

func WithCORS(origin string, h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ApplyCORSHeaders(origin, w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
