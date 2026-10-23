package httpapi

import (
	"encoding/json"
	"net/http"

	"gamification/internal/transport/http/apigen"
)

func writeError(w http.ResponseWriter, status int, code apigen.ErrorBodyErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(apigen.ErrorBody{
		Error: struct {
			Code    apigen.ErrorBodyErrorCode `json:"code"`
			Message string                    `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func notImplemented(w http.ResponseWriter) {
	writeError(w, http.StatusNotImplemented, apigen.ErrorBodyErrorCodeNotImplemented, "not implemented")
}
